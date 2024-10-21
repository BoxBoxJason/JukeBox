package token_controller

import (
	"time"

	db_model "github.com/boxboxjason/jukebox/internal/model"
	"github.com/boxboxjason/jukebox/pkg/logger"
	"github.com/boxboxjason/jukebox/pkg/utils/cryptutils"
	"github.com/boxboxjason/jukebox/pkg/utils/httputils"
	"gorm.io/gorm"
)

// GenerateUserAuthTokens generates an access token and a refresh token for the user
func GenerateUserAuthTokens(db *gorm.DB, user *db_model.User) (string, string, error) {
	// Open db connection
	if db == nil {
		db, err := db_model.OpenConnection()
		if err != nil {
			return "", "", err
		}
		defer db_model.CloseConnection(db)
	}

	// Generate the access token for the user
	access_token, access_string, err := createUserToken(db, user, db_model.ACCESS_TOKEN)
	if err != nil {
		return "", "", err
	}

	// Generate refresh token for the user
	refresh_token, refresh_string, err := createUserToken(db, user, db_model.REFRESH_TOKEN)
	if err != nil {
		return "", "", err
	}

	// Link the refresh token to the access token
	refresh_token.LinkedToken = access_token.ID
	err = refresh_token.UpdateAuthToken(db)
	if err != nil {
		logger.Error("Unable to link the refresh token to the access token for user", user.Username)
		return "", "", httputils.NewInternalServerError("Unable to link the refresh token to the access token")
	}

	return access_string, refresh_string, nil
}

// createUserToken creates a token for the user depending on the token type
func createUserToken(db *gorm.DB, user *db_model.User, token_type string) (*db_model.AuthToken, string, error) {
	// Open db connection
	if db == nil {
		db, err := db_model.OpenConnection()
		if err != nil {
			return &db_model.AuthToken{}, "", err
		}
		defer db_model.CloseConnection(db)
	}

	// Generate an auth token for the user

	string_token, hashed_string_token, err := cryptutils.GenerateHashedToken()
	if err != nil {
		logger.Error(err)
		return &db_model.AuthToken{}, "", err
	}

	token := db_model.AuthToken{
		User:         *user,
		Hashed_Token: hashed_string_token,
		Type:         token_type,
		Expiration:   calculateExpirationTime(token_type),
	}

	// Create the token in the database
	err = token.CreateAuthToken(db)
	if err != nil {
		logger.Error("Unable to create", token_type, "token for user", user.Username)
		return &db_model.AuthToken{}, "", err
	}

	return &token, string_token, nil
}

// calculateExpirationTime calculates the expiration time of the token
func calculateExpirationTime(token_type string) int64 {
	return time.Now().Add(time.Duration(db_model.TOKEN_EXPIRATION_MAP[token_type]) * time.Second).Unix()
}

// RefreshUserAccessToken refreshes the user's access token and refresh token
func RefreshUserAccessToken(db *gorm.DB, user *db_model.User, refresh_token_string string) (string, string, error) {
	// Open db connection
	if db == nil {
		db, err := db_model.OpenConnection()
		if err != nil {
			return "", "", err
		}
		defer db_model.CloseConnection(db)
	}

	// Check if the refresh token is valid
	refresh_token, err := user.CheckAuthTokenMatchesByType(db, refresh_token_string, db_model.REFRESH_TOKEN)
	if err != nil {
		return "", "", httputils.NewUnauthorizedError("Invalid refresh token")
	}

	// Check if the refresh token is expired
	if refresh_token.Expiration < time.Now().Unix() {
		return "", "", httputils.NewUnauthorizedError("Refresh token expired")
	}

	// Generate a new access token and refresh token for the user
	new_access_token_string, new_hashed_access_token_string, err := cryptutils.GenerateHashedToken()
	if err != nil {
		logger.Error(err)
		return "", "", err
	}

	// Retrieve the access token linked to the refresh token
	access_token, err := db_model.GetAuthTokenByID(db, refresh_token.LinkedToken)
	if err != nil {
		logger.Error("Unable to retrieve the access token linked to the refresh token")
		return "", "", httputils.NewNotFoundError("Unable to retrieve the access token linked to the refresh token")
	}

	// Update the access token
	access_token.Hashed_Token = new_hashed_access_token_string
	access_token.Expiration = calculateExpirationTime(db_model.ACCESS_TOKEN)
	err = access_token.UpdateAuthToken(db)
	if err != nil {
		logger.Error("Unable to update the access token")
		return "", "", httputils.NewInternalServerError("Unable to update the access token")
	}

	// Generate a new refresh token for the user
	new_refresh_token_string, new_hashed_refresh_token_string, err := cryptutils.GenerateHashedToken()
	if err != nil {
		logger.Error(err)
		return "", "", err
	}

	// Update the refresh token
	refresh_token.Hashed_Token = new_hashed_refresh_token_string
	refresh_token.Expiration = calculateExpirationTime(db_model.REFRESH_TOKEN)
	err = refresh_token.UpdateAuthToken(db)
	if err != nil {
		logger.Error("Unable to update the refresh token")
		return "", "", httputils.NewInternalServerError("Unable to update the refresh token")
	}

	return new_access_token_string, new_refresh_token_string, nil
}
