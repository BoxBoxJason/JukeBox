package db_controller

import (
	"mime/multipart"
	"regexp"
	"strconv"
	"strings"

	"github.com/boxboxjason/jukebox/internal/constants"
	db_model "github.com/boxboxjason/jukebox/internal/model"
	"github.com/boxboxjason/jukebox/pkg/logger"
	"github.com/boxboxjason/jukebox/pkg/utils/cryptutils"
	"github.com/boxboxjason/jukebox/pkg/utils/fileutils"
	"github.com/boxboxjason/jukebox/pkg/utils/httputils"
	"gorm.io/gorm"
)

var (
	VALID_USERNAME = regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
	VALID_EMAIL    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	VALID_PASSWORD = regexp.MustCompile(`^.{6,99}$`)
)

// ================= CRUD Operations =================

// ================= Create =================

// CreateUser creates a new user in the database after checking the validity of the input fields
// And the uniqueness of the username and email
func CreateUser(db *gorm.DB, query_params *db_model.UsersPostRequestParams) (*db_model.User, error) {
	// Validate user input
	valid_username := VALID_USERNAME.MatchString(query_params.Username)
	valid_email := VALID_EMAIL.MatchString(query_params.Email)
	valid_password := VALID_PASSWORD.MatchString(query_params.Password)

	invalid_fields := make([]string, 0)
	if !valid_username {
		invalid_fields = append(invalid_fields, "username")
	}
	if !valid_email {
		invalid_fields = append(invalid_fields, "email")
	}
	if !valid_password {
		invalid_fields = append(invalid_fields, "password is too short")
	}
	if len(invalid_fields) > 0 {
		return &db_model.User{}, httputils.NewBadRequestError("Invalid fields: " + strings.Join(invalid_fields, ", "))
	}

	hashed_password, err := cryptutils.HashString(query_params.Password)
	if err != nil {
		logger.Error("Unable to hash the password during user creation", err)
		return &db_model.User{}, httputils.NewInternalServerError("Unable to hash the password")
	}

	user := db_model.User{
		Username:        query_params.Username,
		Email:           query_params.Email,
		Hashed_Password: hashed_password,
	}

	// Open db connection
	if db == nil {
		db, err = db_model.OpenConnection()
		if err != nil {
			return &db_model.User{}, err
		}
		defer db_model.CloseConnection(db)
	}

	// Check if user already exists
	_, err = db_model.GetUserByUsername(db, query_params.Username)
	if err == nil {
		return &db_model.User{}, httputils.NewConflictError("Username already exists")
	}
	_, err = db_model.GetUserByEmail(db, query_params.Email)
	if err == nil {
		return &db_model.User{}, httputils.NewConflictError("Email already exists")
	}

	// Create user
	err = user.CreateUser(db)
	if err != nil {
		logger.Error("Unable to create the user in the database")
	} else {
		logger.Info("User", query_params.Username, "created successfully")
	}
	return &user, err
}

// ================= Read =================

// GetUser retrieves a user from the database by ID
func GetUser(db *gorm.DB, id int) (*db_model.User, error) {
	// Open db connection
	if db == nil {
		var err error
		db, err = db_model.OpenConnection()
		if err != nil {
			return nil, err
		}
		defer db_model.CloseConnection(db)
	}

	user, err := db_model.GetUserByID(db, id)
	if err != nil {
		return nil, httputils.NewNotFoundError("User not found")
	}

	return user, nil
}

// GetUserByPartialUsername retrieves a user from the database by partial username
func GetUsersByPartialUsername(db *gorm.DB, partial_username string) ([]*db_model.User, error) {
	// Open db connection
	if db == nil {
		var err error
		db, err = db_model.OpenConnection()
		if err != nil {
			return nil, err
		}
		defer db_model.CloseConnection(db)
	}

	return db_model.GetUsersByUsername(db, partial_username)
}

// GetUsers retrieves all users from the database, applies filters if provided
func GetUsers(db *gorm.DB, query_params *db_model.UsersGetRequestParams) ([]*db_model.User, error) {
	// Sanity checks
	if len(query_params.Admin) > 1 {
		return nil, httputils.NewBadRequestError("Only one value is allowed for the admin parameter")
	} else if query_params.SubscriberTier < 0 || query_params.SubscriberTier > 3 {
		return nil, httputils.NewBadRequestError("Invalid subscriber tier, must be between 0 and 3")
	}
	for _, id := range query_params.ID {
		if id < 0 {
			return nil, httputils.NewBadRequestError("Invalid user ID, must be a positive integer")
		}
	}
	// Open db connection
	if db == nil {
		var err error
		db, err = db_model.OpenConnection()
		if err != nil {
			return nil, err
		}
		defer db_model.CloseConnection(db)
	}

	return db_model.GetUsersByFilters(db, query_params)
}

// ================= Update =================

// UpdateUser updates a user in the database after checking the validity of the input fields
func UpdateUser(db *gorm.DB, user *db_model.User, query_params *db_model.UsersRawPatchRequestParams) (*db_model.User, error) {
	// Validate user input
	valid_username := VALID_USERNAME.MatchString(query_params.Username)
	valid_email := VALID_EMAIL.MatchString(query_params.Email)
	valid_password := VALID_PASSWORD.MatchString(query_params.Password)

	invalid_fields := make([]string, 0)
	if len(query_params.Username) > 0 && !valid_username {
		invalid_fields = append(invalid_fields, "username")
	}
	if len(query_params.Email) > 0 && !valid_email {
		invalid_fields = append(invalid_fields, "email")
	}
	if len(query_params.Password) > 0 && !valid_password {
		invalid_fields = append(invalid_fields, "password")
	}
	if len(invalid_fields) > 0 {
		return &db_model.User{}, httputils.NewBadRequestError("Invalid fields: " + strings.Join(invalid_fields, ", "))
	}
	if query_params.Avatar != nil {
		err := UploadUserAvatar(query_params.Avatar, user.ID)
		if err != nil {
			return &db_model.User{}, err
		}
	}
	hashed_password := ""
	if len(query_params.Password) > 0 {
		new_hashed_password, err := cryptutils.HashString(query_params.Password)
		if err != nil {
			logger.Error("Unable to hash the password during user update", err)
			return &db_model.User{}, httputils.NewInternalServerError("Unable to hash the password")
		}
		hashed_password = new_hashed_password
	}

	// Open db connection
	if db == nil {
		var err error
		db, err = db_model.OpenConnection()
		if err != nil {
			return &db_model.User{}, err
		}
		defer db_model.CloseConnection(db)
	}

	// Check if user already exists
	if query_params.Username != user.Username {
		_, err := db_model.GetUserByUsername(db, query_params.Username)
		if err == nil {
			return &db_model.User{}, httputils.NewConflictError("Username already exists")
		}
	}
	if query_params.Email != user.Email {
		_, err := db_model.GetUserByEmail(db, query_params.Email)
		if err == nil {
			return &db_model.User{}, httputils.NewConflictError("Email already exists")
		}
	}

	// Update user
	if len(query_params.Username) > 0 {
		user.Username = query_params.Username
	}
	if len(query_params.Email) > 0 {
		user.Email = query_params.Email
	}
	if query_params.Avatar != nil {
		user.Avatar = strconv.Itoa(user.ID)
	}
	if len(hashed_password) > 0 {
		user.Hashed_Password = hashed_password
	}

	err := user.UpdateUser(db)
	if err != nil {
		logger.Error("Unable to update the user in the database")
	} else {
		logger.Info("User", user.Username, "updated successfully")
	}
	return user, err
}

// ================= Delete =================

// UserHasPermissionToDeleteUser checks if a user has permission to delete another user
func UserHasPermissionToDeleteUser(user *db_model.User, user_to_delete *db_model.User) bool {
	return user.Admin || user.ID == user_to_delete.ID && !user_to_delete.Banned
}

func DeleteUser(db *gorm.DB, user *db_model.User) error {
	if db == nil {
		var err error
		db, err = db_model.OpenConnection()
		if err != nil {
			return err
		}
		defer db_model.CloseConnection(db)
	}
	return user.DeleteUser(db)
}

func DeleteUsers(db *gorm.DB, requester *db_model.User, query_params *db_model.UsersDeleteRequestParams) error {
	// Open db connection
	if db == nil {
		var err error
		db, err = db_model.OpenConnection()
		if err != nil {
			return err
		}
		defer db_model.CloseConnection(db)
	}

	// Retrieve users
	users, err := db_model.GetUsersByFilters(db, &db_model.UsersGetRequestParams{
		ID:       query_params.ID,
		Username: query_params.Username,
		Email:    query_params.Email,
		Order:    query_params.Order,
		Limit:    query_params.Limit,
		Page:     query_params.Page,
		Offset:   query_params.Offset,
	})
	if err != nil {
		return err
	}

	// Check if requester has permission to delete users
	usernames_to_delete := make([]string, len(users))
	for i, user := range users {
		if !UserHasPermissionToDeleteUser(requester, user) {
			return httputils.NewForbiddenError("User does not have permission to delete user: " + user.Username)
		} else {
			usernames_to_delete[i] = user.Username
		}
	}

	// Log the deletion
	logger.Info("User", requester.Username, "is deleting users", usernames_to_delete, "for reason:", query_params.Reason)

	// Delete users
	return db_model.DeleteUsers(db, users)
}

// UploadUserAvatar uploads a user avatar to the server
func UploadUserAvatar(avatar multipart.File, user_id int) error {
	return fileutils.SaveImageFile(avatar, constants.AVATARS_DIR, strconv.Itoa(user_id))
}
