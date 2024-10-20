package user_controller

import (
	"regexp"
	"strings"

	"github.com/boxboxjason/jukebox/internal/controller/token_controller"
	db_model "github.com/boxboxjason/jukebox/internal/model"
	"github.com/boxboxjason/jukebox/pkg/customerrors"
	"github.com/boxboxjason/jukebox/pkg/logger"
	"github.com/boxboxjason/jukebox/pkg/utils/cryptutils"
)

var (
	VALID_USERNAME = regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
	VALID_EMAIL    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	VALID_PASSWORD = regexp.MustCompile(`^.{6,}$`)
)

// ================= CRUD Operations =================

// ================= Create =================

// CreateUser creates a new user in the database after checking the validity of the input fields
// And the uniqueness of the username and email
func CreateUser(username string, email string, password string) error {
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
		return customerrors.NewBadRequestError("Invalid fields: " + strings.Join(invalid_fields, ", "))
	}

	hashed_password, err := cryptutils.HashString(password)
	if err != nil {
		logger.Error("Unable to hash the password during user creation", err)
		return customerrors.NewInternalServerError("Unable to hash the password")
	}

	user := db_model.User{
		Username:        username,
		Email:           email,
		Hashed_Password: hashed_password,
	}

	// Open db connection
	db, err := db_model.OpenConnection()
	if err != nil {
		return err
	}
	defer db_model.CloseConnection(db)

	// Check if user already exists
	_, err = db_model.GetUserByUsername(db, username)
	if err == nil {
		return customerrors.NewConflictError("Username already exists")
	}
	_, err = db_model.GetUserByEmail(db, email)
	if err == nil {
		return customerrors.NewConflictError("Email already exists")
	}

	// Create user
	err = user.CreateUser(db)
	if err != nil {
		logger.Error("Unable to create the user in the database")
	} else {
		logger.Info("User", username, "created successfully")
	}
	return err
}

// ================= Read =================

// GetUser retrieves a user from the database by ID
func GetUser(id int) (map[string]interface{}, error) {
	// Open db connection
	db, err := db_model.OpenConnection()
	if err != nil {
		return nil, err
	}
	defer db_model.CloseConnection(db)

	user, err := db_model.GetUserByID(db, id)
	if err != nil {
		return nil, customerrors.NewNotFoundError("User not found")
	}

	return user.ToJSON(), nil
}

// GetUserByPartialUsername retrieves a user from the database by partial username
func GetUsersByPartialUsername(partial_username string) ([]map[string]interface{}, error) {
	// Open db connection
	db, err := db_model.OpenConnection()
	if err != nil {
		return nil, err
	}
	defer db_model.CloseConnection(db)

	users, err := db_model.GetUsersByUsername(db, partial_username)
	if err != nil {
		return nil, customerrors.NewNotFoundError("User not found")
	}

	users_json := make([]map[string]interface{}, len(users))
	for i, user := range users {
		users_json[i] = user.ToJSON()
	}

	return users_json, nil
}

// GetAllUsers retrieves all users from the database
func GetAllUsers() ([]map[string]interface{}, error) {
	// Open db connection
	db, err := db_model.OpenConnection()
	if err != nil {
		return nil, err
	}
	defer db_model.CloseConnection(db)

	users, err := db_model.GetAllUsers(db)
	if err != nil {
		return nil, customerrors.NewNotFoundError("User not found")
	}

	users_json := make([]map[string]interface{}, len(users))
	for i, user := range users {
		users_json[i] = user.ToJSON()
	}

	return users_json, nil
}

// LoginUserFromPassword logs in a user by checking the validity of the input fields
// And the correctness of the username and password
func LoginUserFromPassword(username_or_email string, password string) (string, string, error) {
	// Open db connection
	db, err := db_model.OpenConnection()
	if err != nil {
		return "", "", err
	}
	defer db_model.CloseConnection(db)

	// Retrieve the user (if it exists)
	user, err := db_model.GetUserByUsernameOREmail(db, username_or_email)
	if err != nil {
		return "", "", customerrors.NewUnauthorizedError("Invalid credentials combination")
	}

	// Check if the password matches
	if !user.CheckPasswordMatches(password) {
		return "", "", customerrors.NewUnauthorizedError("Invalid credentials combination")
	}

	// Generate the user's auth token
	access_token, refresh_token, err := token_controller.GenerateUserAuthTokens(db, user)
	if err != nil {
		return "", "", err
	}

	return access_token, refresh_token, nil
}

func LoginFromToken(user_id int, token_string string) error {
	// Open db connection
	db, err := db_model.OpenConnection()
	if err != nil {
		return err
	}
	defer db_model.CloseConnection(db)

	// Retrieve the user
	user, err := db_model.GetUserByID(db, user_id)
	if err != nil {
		return customerrors.NewUnauthorizedError("Invalid token")
	}

	// Check if the token matches
	_, err = user.CheckAuthTokenMatchesByType(db, token_string, db_model.ACCESS_TOKEN)
	if err != nil {
		return customerrors.NewUnauthorizedError("Invalid token")
	}

	return nil
}

// ================= Update =================

// ================= Delete =================
