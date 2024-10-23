package db_controller

import (
	"regexp"
	"strings"

	db_model "github.com/boxboxjason/jukebox/internal/model"
	"github.com/boxboxjason/jukebox/pkg/logger"
	"github.com/boxboxjason/jukebox/pkg/utils/cryptutils"
	"github.com/boxboxjason/jukebox/pkg/utils/httputils"
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

// GetAllUsers retrieves all users from the database
func GetAllUsers() ([]*db_model.User, error) {
	// Open db connection
	db, err := db_model.OpenConnection()
	if err != nil {
		return nil, err
	}
	defer db_model.CloseConnection(db)

	return db_model.GetAllUsers(db)
}

// ================= Update =================

// ================= Delete =================
