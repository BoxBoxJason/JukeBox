package db_controller

import (
	"regexp"
	"strings"

	db_model "github.com/boxboxjason/jukebox/internal/model"
	"github.com/boxboxjason/jukebox/pkg/logger"
	"github.com/boxboxjason/jukebox/pkg/utils/cryptutils"
	"github.com/boxboxjason/jukebox/pkg/utils/httputils"
	"gorm.io/gorm"
)

var (
	VALID_USERNAME = regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
	VALID_EMAIL    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	VALID_PASSWORD = regexp.MustCompile(`^.{6,32}$`)
)

// ================= CRUD Operations =================

// ================= Create =================

// CreateUser creates a new user in the database after checking the validity of the input fields
// And the uniqueness of the username and email
func CreateUser(username string, email string, password string) (*db_model.User, error) {
	// Validate user input
	valid_username := VALID_USERNAME.MatchString(username)
	valid_email := VALID_EMAIL.MatchString(email)
	valid_password := VALID_PASSWORD.MatchString(password)

	invalid_fields := make([]string, 0)
	if !valid_username {
		invalid_fields = append(invalid_fields, "username")
	}
	if !valid_email {
		invalid_fields = append(invalid_fields, "email")
	}
	if !valid_password {
		invalid_fields = append(invalid_fields, "password")
	}
	if len(invalid_fields) > 0 {
		return &db_model.User{}, httputils.NewBadRequestError("Invalid fields: " + strings.Join(invalid_fields, ", "))
	}

	hashed_password, err := cryptutils.HashString(password)
	if err != nil {
		logger.Error("Unable to hash the password during user creation", err)
		return &db_model.User{}, httputils.NewInternalServerError("Unable to hash the password")
	}

	user := db_model.User{
		Username:        username,
		Email:           email,
		Hashed_Password: hashed_password,
	}

	// Open db connection
	db, err := db_model.OpenConnection()
	if err != nil {
		return &db_model.User{}, err
	}
	defer db_model.CloseConnection(db)

	// Check if user already exists
	_, err = db_model.GetUserByUsername(db, username)
	if err == nil {
		return &db_model.User{}, httputils.NewConflictError("Username already exists")
	}
	_, err = db_model.GetUserByEmail(db, email)
	if err == nil {
		return &db_model.User{}, httputils.NewConflictError("Email already exists")
	}

	// Create user
	err = user.CreateUser(db)
	if err != nil {
		logger.Error("Unable to create the user in the database")
	} else {
		logger.Info("User", username, "created successfully")
	}
	return &user, err
}

// ================= Read =================

// GetUser retrieves a user from the database by ID
func GetUser(id int) (*db_model.User, error) {
	// Open db connection
	db, err := db_model.OpenConnection()
	if err != nil {
		return nil, err
	}
	defer db_model.CloseConnection(db)

	user, err := db_model.GetUserByID(db, id)
	if err != nil {
		return nil, httputils.NewNotFoundError("User not found")
	}

	return user, nil
}

// GetUserByPartialUsername retrieves a user from the database by partial username
func GetUsersByPartialUsername(partial_username string) ([]*db_model.User, error) {
	// Open db connection
	db, err := db_model.OpenConnection()
	if err != nil {
		return nil, err
	}
	defer db_model.CloseConnection(db)

	return db_model.GetUsersByUsername(db, partial_username)
}

// GetUsers retrieves all users from the database, applies filters if provided
func GetUsers(ids []int, usernames []string, partial_usernames []string, emails []string, banned []bool, admin []bool, minimum_subscriber_tier int) ([]*db_model.User, error) {
	// Sanity checks
	if len(banned) > 1 {
		return nil, httputils.NewBadRequestError("Only one value is allowed for the banned parameter")
	} else if len(admin) > 1 {
		return nil, httputils.NewBadRequestError("Only one value is allowed for the admin parameter")
	} else if minimum_subscriber_tier < 0 || minimum_subscriber_tier > 3 {
		return nil, httputils.NewBadRequestError("Invalid subscriber tier, must be between 0 and 3")
	}
	for _, id := range ids {
		if id < 0 {
			return nil, httputils.NewBadRequestError("Invalid user ID, must be a positive integer")
		}
	}
	// Open db connection
	db, err := db_model.OpenConnection()
	if err != nil {
		return nil, err
	}
	defer db_model.CloseConnection(db)

	return db_model.GetUsersByFilters(db, ids, usernames, emails, partial_usernames, banned, admin, minimum_subscriber_tier)
}

// ================= Update =================

// ================= Delete =================
func UserHasPermissionToDeleteUser(user *db_model.User, user_to_delete *db_model.User) bool {
	return user.Admin || user.ID == user_to_delete.ID && !user_to_delete.Banned
}

func DeleteUser(db *gorm.DB, user *db_model.User) error {
	if db == nil {
		db, err := db_model.OpenConnection()
		if err != nil {
			return err
		}
		defer db_model.CloseConnection(db)
	}
	return user.DeleteUser(db)
}

func DeleteUsers(db *gorm.DB, requester *db_model.User, ids []int, usernames []string, emails []string, reason string) error {
	// Open db connection
	if db == nil {
		db, err := db_model.OpenConnection()
		if err != nil {
			return err
		}
		defer db_model.CloseConnection(db)
	}

	// Retrieve users
	users, err := db_model.GetUsersByFilters(db, ids, usernames, emails, nil, nil, nil, 0)
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
	logger.Info("User", requester.Username, "is deleting users", usernames_to_delete, "for reason:", reason)

	// Delete users
	return db_model.DeleteUsers(db, users)
}
