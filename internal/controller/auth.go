package db_controller

import (
	"github.com/boxboxjason/jukebox/internal/constants"
	"github.com/boxboxjason/jukebox/internal/middlewares"
	db_model "github.com/boxboxjason/jukebox/internal/model"
	"github.com/boxboxjason/jukebox/pkg/utils/httputils"
)

// LoginUserFromPassword logs in a user by checking the validity of the input fields
// And the correctness of the username and password
func LoginUserFromPassword(username_or_email string, password string) (int, string, string, string, error) {
	// Open db connection
	db, err := db_model.OpenConnection()
	if err != nil {
		return -1, "", "", "", err
	}
	defer db_model.CloseConnection(db)

	// Retrieve the user (if it exists)
	user, err := db_model.GetUserByUsernameOREmail(db, username_or_email)
	if err != nil {
		return -1, "", "", "", httputils.NewUnauthorizedError("Invalid credentials combination")
	}

	// Check if the password matches
	if !user.CheckPasswordMatches(password) {
		return -1, "", "", "", httputils.NewUnauthorizedError("Invalid credentials combination")
	}

	// Generate the user's auth token
	access_token, refresh_token, err := GenerateUserAuthTokens(db, user)
	if err != nil {
		return -1, "", "", "", err
	}

	return user.ID, user.Username, access_token, refresh_token, nil
}

// LoginFromToken logs in a user by checking the validity of the token
// Refreshing the access token if it is valid and returning the new access token
// Returns the user id, username, and the new access token
func LoginFromToken(user_id int, token_string string) (int, string, string, string, error) {
	// Open db connection
	db, err := db_model.OpenConnection()
	if err != nil {
		return -1, "", "", "", err
	}
	defer db_model.CloseConnection(db)

	// Retrieve the user
	user, err := db_model.GetUserByID(db.Preload("AuthToken"), user_id)
	if err != nil {
		return -1, "", "", "", httputils.NewUnauthorizedError("Invalid token")
	}

	// Check if the token matches
	access_token, err := user.CheckAuthTokenMatchesByType(db, token_string, constants.ACCESS_TOKEN)
	if err != nil {
		return -1, "", "", "", err
	}

	// Refresh the access token
	access_token_string, err := RefreshToken(db, access_token)
	if err != nil {
		return -1, "", "", "", err
	}

	// Generate a refresh token
	_, refresh_token_string, err := createUserToken(db, user, constants.REFRESH_TOKEN)
	if err != nil {
		return -1, "", "", "", err
	}

	return user.ID, user.Username, access_token_string, refresh_token_string, nil
}

// RefreshTokens refreshes the access token and the refresh token
// Returns the user id, username, access token and refresh token
func RefreshTokens(identity_bearer string) (int, string, string, string, error) {
	// Check if the identity bearer is valid
	user_id, token_string, err := middlewares.DecodeIdentityBearerToUserAndToken(identity_bearer)
	if err != nil {
		return -1, "", "", "", err
	}

	// Open db connection
	db, err := db_model.OpenConnection()
	if err != nil {
		return -1, "", "", "", err
	}
	defer db_model.CloseConnection(db)

	// Retrieve the user
	user, err := db_model.GetUserByID(db, user_id)
	if err != nil {
		return -1, "", "", "", httputils.NewUnauthorizedError("Invalid token")
	}

	// Check if the token matches
	refresh_token, err := user.CheckAuthTokenMatchesByType(db, token_string, constants.REFRESH_TOKEN)
	if err != nil {
		return -1, "", "", "", err
	}

	// Update the refresh token
	refresh_token_string, err := RefreshToken(db, refresh_token)
	if err != nil {
		return -1, "", "", "", err
	}

	// Retrieve the linked access token if it exists OR create it
	var access_token_string string
	access_token, err := refresh_token.GetLinkedToken(db)
	if err != nil {
		access_token, access_token_string, err = createUserToken(db, refresh_token.User, constants.ACCESS_TOKEN)
		if err != nil {
			return -1, "", "", "", err
		}
		refresh_token.LinkedToken = access_token
	} else {
		access_token_string, err = RefreshToken(db, access_token)
		if err != nil {
			return -1, "", "", "", err
		}
	}

	return user.ID, user.Username, access_token_string, refresh_token_string, nil
}
