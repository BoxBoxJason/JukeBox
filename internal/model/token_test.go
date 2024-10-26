package db_model

import (
	"testing"

	"github.com/boxboxjason/jukebox/internal/constants"
	"github.com/boxboxjason/jukebox/pkg/utils/cryptutils"
)

func TestCreateAuthToken(t *testing.T) {

	db, err := OpenConnection()
	if err != nil {
		t.Errorf("Error opening connection to the database: %v", err)
	}
	defer CloseConnection(db)

	// Create a user for the auth token
	user := &User{
		Email:           "test_email_100@test.com",
		Hashed_Password: "hashed_password",
		Username:        "test_user_100",
	}

	err = user.CreateUser(db)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	// Create a new auth token
	auth_token := &AuthToken{
		Hashed_Token: "test_hashed_token_1",
		User:         *user,
	}

	err = auth_token.CreateAuthToken(db)
	if err != nil {
		t.Errorf("Error creating auth token: %v", err)
	}
}

func TestCreateAuthTokens(t *testing.T) {

	db, err := OpenConnection()
	if err != nil {
		t.Errorf("Error opening connection to the database: %v", err)
	}
	defer CloseConnection(db)

	// Create a user for the auth token
	user := &User{
		Email:           "test_email_101@test.com",
		Hashed_Password: "hashed_password",
		Username:        "test_user_101",
	}

	err = user.CreateUser(db)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	// Create multiple auth tokens
	auth_tokens := []*AuthToken{
		{
			Hashed_Token: "test_hashed_token_2",
			User:         *user,
		},
		{
			Hashed_Token: "test_hashed_token_3",
			User:         *user,
		},
	}

	err = CreateAuthTokens(db, auth_tokens)
	if err != nil {
		t.Errorf("Error creating auth tokens: %v", err)
	}
}

func TestGetAuthTokenByID(t *testing.T) {

	db, err := OpenConnection()
	if err != nil {
		t.Errorf("Error opening connection to the database: %v", err)
	}
	defer CloseConnection(db)

	// Create a user for the auth token
	user := &User{
		Email:           "test_email_102@test.com",
		Hashed_Password: "hashed_password",
		Username:        "test_user_102",
	}

	err = user.CreateUser(db)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	// Create an auth token
	auth_token := &AuthToken{
		Hashed_Token: "test_hashed_token_4",
		User:         *user,
	}

	err = auth_token.CreateAuthToken(db)
	if err != nil {
		t.Errorf("Error creating auth token: %v", err)
	}

	// Retrieve the auth token by ID
	retrieved_auth_token, err := GetAuthTokenByID(db, auth_token.ID)
	if err != nil {
		t.Errorf("Error retrieving auth token by ID: %v", err)
	}

	if retrieved_auth_token.Hashed_Token != auth_token.Hashed_Token {
		t.Errorf("Retrieved auth token does not match created auth token")
	}
}

func TestGetUserTokensByType(t *testing.T) {

	db, err := OpenConnection()
	if err != nil {
		t.Errorf("Error opening connection to the database: %v", err)
	}
	defer CloseConnection(db)

	// Create a user for the auth tokens
	user := &User{
		Email:           "test_email_103@test.com",
		Hashed_Password: "hashed_password",
		Username:        "test_user_103",
	}

	err = user.CreateUser(db)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	// Create access auth token
	access_auth_token := &AuthToken{
		Hashed_Token: "test_hashed_token_5",
		User:         *user,
		Type:         constants.ACCESS_TOKEN,
	}

	err = access_auth_token.CreateAuthToken(db)
	if err != nil {
		t.Errorf("Error creating access auth token: %v", err)
	}

	// Create refresh auth token
	refresh_auth_token := &AuthToken{
		Hashed_Token: "test_hashed_token_6",
		User:         *user,
		Type:         constants.REFRESH_TOKEN,
	}

	err = refresh_auth_token.CreateAuthToken(db)
	if err != nil {
		t.Errorf("Error creating refresh auth token: %v", err)
	}

	// Retrieve all access tokens for the user
	access_tokens, err := user.GetUserTokensByType(db, constants.ACCESS_TOKEN)
	if err != nil {
		t.Errorf("Error retrieving access tokens: %v", err)
	}

	if len(access_tokens) < 1 {
		t.Errorf("No access tokens retrieved")
	}

	// Retrieve all refresh tokens for the user
	refresh_tokens, err := user.GetUserTokensByType(db, constants.REFRESH_TOKEN)
	if err != nil {
		t.Errorf("Error retrieving refresh tokens: %v", err)
	}

	if len(refresh_tokens) < 1 {
		t.Errorf("No refresh tokens retrieved")
	}
}

func TestGetLinkedToken(t *testing.T) {

	db, err := OpenConnection()
	if err != nil {
		t.Errorf("Error opening connection to the database: %v", err)
	}
	defer CloseConnection(db)

	// Create a user for the auth tokens
	user := &User{
		Email:           "test_email_104@test.com",
		Hashed_Password: "hashed_password",
		Username:        "test_user_104",
	}

	err = user.CreateUser(db)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	// Create an access auth token
	access_auth_token := &AuthToken{
		Hashed_Token: "test_hashed_token_7",
		User:         *user,
		Type:         constants.ACCESS_TOKEN,
	}

	err = access_auth_token.CreateAuthToken(db)
	if err != nil {
		t.Errorf("Error creating access auth token: %v", err)
	}

	// Create a refresh auth token
	refresh_auth_token := &AuthToken{
		Hashed_Token: "test_hashed_token_8",
		User:         *user,
		Type:         constants.REFRESH_TOKEN,
		LinkedToken:  access_auth_token.ID,
	}

	err = refresh_auth_token.CreateAuthToken(db)
	if err != nil {
		t.Errorf("Error creating refresh auth token: %v", err)
	}

	// Retrieve the linked token
	linked_token, err := refresh_auth_token.GetLinkedToken(db)
	if err != nil {
		t.Errorf("Error retrieving linked token: %v", err)
	}

	if linked_token.ID != access_auth_token.ID {
		t.Errorf("Linked token does not match access token")
	}
}

func TestCheckAuthTokenMatchesByType(t *testing.T) {

	db, err := OpenConnection()
	if err != nil {
		t.Errorf("Error opening connection to the database: %v", err)
	}
	defer CloseConnection(db)

	// Create a user for the auth tokens
	user := &User{
		Email:           "test_email_105@test.com",
		Hashed_Password: "hashed_password",
		Username:        "test_user_105",
	}

	err = user.CreateUser(db)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	raw_token, hashed_token, err := cryptutils.GenerateHashedToken()
	if err != nil {
		t.Errorf("Error generating hashed token: %v", err)
	}

	// Create an access auth token
	access_auth_token := &AuthToken{
		Hashed_Token: hashed_token,
		User:         *user,
		Type:         constants.ACCESS_TOKEN,
	}

	err = access_auth_token.CreateAuthToken(db)
	if err != nil {
		t.Errorf("Error creating access auth token: %v", err)
	}

	// Check if the access token matches
	_, err = user.CheckAuthTokenMatchesByType(db, raw_token, constants.ACCESS_TOKEN)
	if err != nil {
		t.Errorf("Error checking if access token matches: %v", err)
	}
}

func TestUpdateAuthToken(t *testing.T) {

	db, err := OpenConnection()
	if err != nil {
		t.Errorf("Error opening connection to the database: %v", err)
	}
	defer CloseConnection(db)

	// Create a user for the auth token
	user := &User{
		Email:           "test_email_106@test.com",
		Hashed_Password: "hashed_password",
		Username:        "test_user_106",
	}

	err = user.CreateUser(db)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	// Create an auth token
	auth_token := &AuthToken{
		Hashed_Token: "test_hashed_token_9",
		User:         *user,
	}

	err = auth_token.CreateAuthToken(db)
	if err != nil {
		t.Errorf("Error creating auth token: %v", err)
	}

	// Update the auth token
	auth_token.Hashed_Token = "test_hashed_token_9_updated"

	err = auth_token.UpdateAuthToken(db)
	if err != nil {
		t.Errorf("Error updating auth token: %v", err)
	}

	// Retrieve the updated auth token
	retrieved_auth_token, err := GetAuthTokenByID(db, auth_token.ID)
	if err != nil {
		t.Errorf("Error retrieving auth token by ID: %v", err)
	}

	if retrieved_auth_token.Hashed_Token != auth_token.Hashed_Token {
		t.Errorf("Updated auth token does not match")
	}
}

func TestDeleteAuthToken(t *testing.T) {

	db, err := OpenConnection()
	if err != nil {
		t.Errorf("Error opening connection to the database: %v", err)
	}
	defer CloseConnection(db)

	// Create a user for the auth token
	user := &User{
		Email:           "test_email_107@test.com",
		Hashed_Password: "hashed_password",
		Username:        "test_user_107",
	}

	err = user.CreateUser(db)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	// Create an auth token
	auth_token := &AuthToken{
		Hashed_Token: "test_hashed_token_10",
		User:         *user,
	}

	err = auth_token.CreateAuthToken(db)
	if err != nil {
		t.Errorf("Error creating auth token: %v", err)
	}

	// Delete the auth token
	err = auth_token.DeleteAuthToken(db)
	if err != nil {
		t.Errorf("Error deleting auth token: %v", err)
	}

	// Retrieve the deleted auth token
	_, err = GetAuthTokenByID(db, auth_token.ID)
	if err == nil {
		t.Errorf("Auth token not deleted")
	}
}

func TestDeleteAuthTokens(t *testing.T) {

	db, err := OpenConnection()
	if err != nil {
		t.Errorf("Error opening connection to the database: %v", err)
	}
	defer CloseConnection(db)

	// Create a user for the auth tokens
	user := &User{
		Email:           "test_email_108@test.com",
		Hashed_Password: "hashed_password",
		Username:        "test_user_108",
	}

	err = user.CreateUser(db)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}

	// Create multiple auth tokens
	auth_tokens := []*AuthToken{
		{
			Hashed_Token: "test_hashed_token_11",
			User:         *user,
		},
		{
			Hashed_Token: "test_hashed_token_12",
			User:         *user,
		},
	}

	err = CreateAuthTokens(db, auth_tokens)
	if err != nil {
		t.Errorf("Error creating auth tokens: %v", err)
	}

	// Delete the auth tokens
	err = DeleteAuthTokens(db, auth_tokens)
	if err != nil {
		t.Errorf("Error deleting auth tokens: %v", err)
	}

	// Retrieve the deleted auth tokens
	for _, auth_token := range auth_tokens {
		_, err = GetAuthTokenByID(db, auth_token.ID)
		if err == nil {
			t.Errorf("Auth token not deleted")
		}
	}
}
