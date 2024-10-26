package db_model

import (
	"testing"

	"github.com/boxboxjason/jukebox/pkg/utils/cryptutils"
)

func TestUserCreate(t *testing.T) {
	// Create a new user
	user := &User{
		Username:        "test_user_1",
		Email:           "test_email_1@test.com",
		Hashed_Password: "hashed_password",
	}
	db, err := OpenConnection()
	if err != nil {
		t.Errorf("Error opening connection: %v", err)
	}
	defer CloseConnection(db)

	err = user.CreateUser(db)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}
}

func TestUsersCreate(t *testing.T) {
	// Create a new user
	user1 := &User{
		Username:        "test_user_2",
		Email:           "test_email_2@test.com",
		Hashed_Password: "hashed_password",
	}
	user2 := &User{
		Username:        "test_user_3",
		Email:           "test_email_3@test.com",
		Hashed_Password: "hashed_password",
	}

	db, err := OpenConnection()
	if err != nil {
		t.Errorf("Error opening connection: %v", err)
	}
	defer CloseConnection(db)

	err = CreateUsers(db, []*User{user1, user2})
	if err != nil {
		t.Errorf("Error creating users: %v", err)
	}
}

func TestGetUserById(t *testing.T) {
	// Create a new user
	user := &User{
		Username:        "test_user_4",
		Email:           "test_email_4@test.com",
		Hashed_Password: "hashed_password",
	}
	db, err := OpenConnection()
	if err != nil {
		t.Errorf("Error opening connection: %v", err)
	}
	defer CloseConnection(db)

	err = user.CreateUser(db)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	// Retrieve the user by ID
	retrieved_user, err := GetUserByID(db, user.ID)
	if err != nil {
		t.Errorf("Error retrieving user: %v", err)
	}

	if retrieved_user.Username != user.Username {
		t.Errorf("Expected username %s, got %s", user.Username, retrieved_user.Username)
	}

	if retrieved_user.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, retrieved_user.Email)
	}

	if retrieved_user.Hashed_Password != user.Hashed_Password {
		t.Errorf("Expected hashed password %s, got %s", user.Hashed_Password, retrieved_user.Hashed_Password)
	}
}

func TestGetUserByUsername(t *testing.T) {
	// Create a new user
	user := &User{
		Username:        "test_user_5",
		Email:           "test_email_5@test.com",
		Hashed_Password: "hashed_password",
	}
	db, err := OpenConnection()
	if err != nil {
		t.Errorf("Error opening connection: %v", err)
	}
	defer CloseConnection(db)

	err = user.CreateUser(db)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	// Retrieve the user by username
	retrieved_user, err := GetUserByUsername(db, user.Username)
	if err != nil {
		t.Errorf("Error retrieving user: %v", err)
	}

	if retrieved_user.Username != user.Username {
		t.Errorf("Expected username %s, got %s", user.Username, retrieved_user.Username)
	}

	if retrieved_user.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, retrieved_user.Email)
	}

	if retrieved_user.Hashed_Password != user.Hashed_Password {
		t.Errorf("Expected hashed password %s, got %s", user.Hashed_Password, retrieved_user.Hashed_Password)
	}
}

func TestGetUserByEmail(t *testing.T) {
	// Create a new user
	user := &User{
		Username:        "test_user_6",
		Email:           "test_email_6@test.com",
		Hashed_Password: "hashed_password",
	}
	db, err := OpenConnection()
	if err != nil {
		t.Errorf("Error opening connection: %v", err)
	}
	defer CloseConnection(db)

	err = user.CreateUser(db)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	// Retrieve the user by email
	retrieved_user, err := GetUserByEmail(db, user.Email)
	if err != nil {
		t.Errorf("Error retrieving user: %v", err)
	}

	if retrieved_user.Username != user.Username {
		t.Errorf("Expected username %s, got %s", user.Username, retrieved_user.Username)
	}

	if retrieved_user.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, retrieved_user.Email)
	}

	if retrieved_user.Hashed_Password != user.Hashed_Password {
		t.Errorf("Expected hashed password %s, got %s", user.Hashed_Password, retrieved_user.Hashed_Password)
	}
}

func TestGetUserByUsernameOREmail(t *testing.T) {
	// Create a new user
	user := &User{
		Username:        "test_user_7",
		Email:           "test_email_7@test.com",
		Hashed_Password: "hashed_password",
	}
	db, err := OpenConnection()
	if err != nil {
		t.Errorf("Error opening connection: %v", err)
	}
	defer CloseConnection(db)

	err = user.CreateUser(db)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	// Retrieve the user by username or email
	retrieved_user, err := GetUserByUsernameOREmail(db, user.Username)
	if err != nil {
		t.Errorf("Error retrieving user by username: %v", err)
	}

	if retrieved_user.Username != user.Username {
		t.Errorf("Expected username %s, got %s", user.Username, retrieved_user.Username)
	}

	if retrieved_user.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, retrieved_user.Email)
	}

	if retrieved_user.Hashed_Password != user.Hashed_Password {
		t.Errorf("Expected hashed password %s, got %s", user.Hashed_Password, retrieved_user.Hashed_Password)
	}

	// Retrieve the user by username or email
	retrieved_user, err = GetUserByUsernameOREmail(db, user.Email)
	if err != nil {
		t.Errorf("Error retrieving user by email: %v", err)
	}

	if retrieved_user.Username != user.Username {
		t.Errorf("Expected username %s, got %s", user.Username, retrieved_user.Username)
	}

	if retrieved_user.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, retrieved_user.Email)
	}

	if retrieved_user.Hashed_Password != user.Hashed_Password {
		t.Errorf("Expected hashed password %s, got %s", user.Hashed_Password, retrieved_user.Hashed_Password)
	}
}

func TestGetUsersByUsername(t *testing.T) {
	// Create a new user
	user1 := &User{
		Username:        "test_user_8",
		Email:           "test_email_8@test.com",
		Hashed_Password: "hashed_password",
	}
	user2 := &User{
		Username:        "test_user_9",
		Email:           "test_email_9@test.com",
		Hashed_Password: "hashed_password",
	}

	db, err := OpenConnection()
	if err != nil {
		t.Errorf("Error opening connection: %v", err)
	}
	defer CloseConnection(db)

	err = CreateUsers(db, []*User{user1, user2})
	if err != nil {
		t.Errorf("Error creating users: %v", err)
	}

	// Retrieve the users by username
	retrieved_users, err := GetUsersByUsername(db, "test_user")
	if err != nil {
		t.Errorf("Error retrieving users: %v", err)
	}

	if len(retrieved_users) < 2 {
		t.Errorf("Expected at least 2 users, got %d", len(retrieved_users))
	}
}

func TestGetAllUsers(t *testing.T) {
	// Create a new user
	user1 := &User{
		Username:        "test_user_10",
		Email:           "test_email_10@test.com",
		Hashed_Password: "hashed_password",
	}
	user2 := &User{
		Username:        "test_user_11",
		Email:           "test_email_11@test.com",
		Hashed_Password: "hashed_password",
	}

	db, err := OpenConnection()
	if err != nil {
		t.Errorf("Error opening connection: %v", err)
	}
	defer CloseConnection(db)

	err = CreateUsers(db, []*User{user1, user2})
	if err != nil {
		t.Errorf("Error creating users: %v", err)
	}

	// Retrieve all users
	retrieved_users, err := GetAllUsers(db)
	if err != nil {
		t.Errorf("Error retrieving users: %v", err)
	}

	if len(retrieved_users) < 2 {
		t.Errorf("Expected at least 2 users, got %d", len(retrieved_users))
	}
}

func TestUserUpdate(t *testing.T) {
	// Create a new user
	user := &User{
		Username:        "test_user_12",
		Email:           "test_user_12@gmail.com",
		Hashed_Password: "hashed_password",
	}
	db, err := OpenConnection()
	if err != nil {
		t.Errorf("Error opening connection: %v", err)
	}
	defer CloseConnection(db)

	err = user.CreateUser(db)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	// Update the user
	user.Username = "test_user_12_updated"
	err = user.UpdateUser(db)
	if err != nil {
		t.Errorf("Error updating user: %v", err)
	}

	// Retrieve the user by ID
	retrieved_user, err := GetUserByID(db, user.ID)
	if err != nil {
		t.Errorf("Error retrieving user: %v", err)
	}

	if retrieved_user.Username != user.Username {
		t.Errorf("Expected username %s, got %s", user.Username, retrieved_user.Username)
	}
}

func TestUserDelete(t *testing.T) {
	// Create a new user
	user := &User{
		Username:        "test_user_13",
		Email:           "test_user_13@gmail.com",
		Hashed_Password: "hashed_password",
	}
	db, err := OpenConnection()
	if err != nil {
		t.Errorf("Error opening connection: %v", err)
	}
	defer CloseConnection(db)

	err = user.CreateUser(db)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	// Delete the user
	err = user.DeleteUser(db)
	if err != nil {
		t.Errorf("Error deleting user: %v", err)
	}

	// Retrieve the user by ID
	_, err = GetUserByID(db, user.ID)
	if err == nil {
		t.Errorf("Expected error retrieving user, got nil")
	}
}

func TestUsersDelete(t *testing.T) {
	// Create a new user
	user1 := &User{
		Username:        "test_user_14",
		Email:           "test_user_14@gmail.com",
		Hashed_Password: "hashed_password",
	}
	user2 := &User{
		Username:        "test_user_15",
		Email:           "test_user_15@gmail.com",
		Hashed_Password: "hashed_password",
	}

	db, err := OpenConnection()
	if err != nil {
		t.Errorf("Error opening connection: %v", err)
	}
	defer CloseConnection(db)

	err = CreateUsers(db, []*User{user1, user2})
	if err != nil {
		t.Errorf("Error creating users: %v", err)
	}

	// Delete the users
	err = DeleteUsers(db, []*User{user1, user2})
	if err != nil {
		t.Errorf("Error deleting users: %v", err)
	}

	// Retrieve the users by ID
	_, err = GetUserByID(db, user1.ID)
	if err == nil {
		t.Errorf("Expected error retrieving user, got nil")
	}

	_, err = GetUserByID(db, user2.ID)
	if err == nil {
		t.Errorf("Expected error retrieving user, got nil")
	}

}

func TestCheckPasswordMatches(t *testing.T) {

	password := "test_password"
	hashed_password, err := cryptutils.HashString(password)
	if err != nil {
		t.Errorf("Error hashing password: %v", err)
	}

	// Create a new user
	user := &User{
		Username:        "test_user_16",
		Email:           "test_user_16@gmail.com",
		Hashed_Password: hashed_password,
	}
	db, err := OpenConnection()
	if err != nil {
		t.Errorf("Error opening connection: %v", err)
	}
	defer CloseConnection(db)

	err = user.CreateUser(db)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	// Check if the password matches
	matches := user.CheckPasswordMatches(password)
	if !matches {
		t.Errorf("Expected password to match, got false")
	}

	// Check if the password matches
	matches = user.CheckPasswordMatches("wrong_password")
	if matches {
		t.Errorf("Expected password to not match, got true")
	}
}
