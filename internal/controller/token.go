package db_controller

import (
	"time"

	"github.com/boxboxjason/jukebox/internal/constants"
	"github.com/boxboxjason/jukebox/internal/middlewares"
	db_model "github.com/boxboxjason/jukebox/internal/model"
	"github.com/boxboxjason/jukebox/pkg/logger"
	"github.com/boxboxjason/jukebox/pkg/utils/cryptutils"
	"github.com/boxboxjason/jukebox/pkg/utils/httputils"
	"gorm.io/gorm"
)

// ================= CRUD Operations =================

// ================= Create =================

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
	access_token, access_string, err := createUserToken(db, user, constants.ACCESS_TOKEN)
	if err != nil {
		return "", "", err
	}

	// Generate refresh token for the user
	refresh_token, refresh_string, err := createUserToken(db, user, constants.REFRESH_TOKEN)
	if err != nil {
		return "", "", err
	}

	// Link the refresh token to the access token
	refresh_token.LinkedToken = access_token
	err = refresh_token.UpdateAuthToken(db)
	if err != nil {
		logger.Error("Unable to link the refresh token to the access token for user", user.Username)
		return "", "", httputils.NewInternalServerError("Unable to link the refresh token to the access token")
	}

	// Link the access token to the refresh token
	access_token.LinkedToken = refresh_token
	err = access_token.UpdateAuthToken(db)
	if err != nil {
		logger.Error("Unable to link the access token to the refresh token for user", user.Username)
		return "", "", httputils.NewInternalServerError("Unable to link the access token to the refresh token")
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
		User:         user,
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

	return &token, middlewares.EncodeUserAndTokenToIdentityBearer(user.ID, string_token), nil
}

// calculateExpirationTime calculates the expiration time of the token
func calculateExpirationTime(token_type string) int64 {
	return time.Now().Add(time.Duration(constants.TOKEN_EXPIRATION_MAP[token_type]) * time.Hour).Unix()
}

// ================= Read =================

// ================= Update =================

// RefreshToken refreshes a token with a new hash value
func RefreshToken(db *gorm.DB, token *db_model.AuthToken) (string, error) {
	// Open db connection
	if db == nil {
		db, err := db_model.OpenConnection()
		if err != nil {
			return "", err
		}
		defer db_model.CloseConnection(db)
	}

	// Generate new token and hash
	new_token_string, new_token_string_hash, err := cryptutils.GenerateHashedToken()
	if err != nil {
		logger.Error("Unable to generate new token", err)
		return "", err
	}

	// Update the token with the new hash
	token.Hashed_Token = new_token_string_hash
	token.Expiration = calculateExpirationTime(token.Type)
	err = token.UpdateAuthToken(db)
	if err != nil {
		logger.Error("Unable to update token", err)
		return "", err
	}

	return middlewares.EncodeUserAndTokenToIdentityBearer(token.User.ID, new_token_string), nil
}

// ================= Delete =================

func DeleteToken(db *gorm.DB, token *db_model.AuthToken) error {
	// Open db connection
	if db == nil {
		db, err := db_model.OpenConnection()
		if err != nil {
			return err
		}
		defer db_model.CloseConnection(db)
	}

	linked_token, err := token.GetLinkedToken(db)
	if err == nil {
		err = linked_token.DeleteAuthToken(db)
		if err != nil {
			logger.Error("Unable to delete linked token")
			return err
		}
	}

	err = token.DeleteAuthToken(db)
	if err != nil {
		logger.Error("Unable to delete token")
		return err
	}

	return nil
}
