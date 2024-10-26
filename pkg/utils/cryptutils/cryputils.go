package cryptutils

import (
	"crypto/rand"

	"github.com/boxboxjason/jukebox/pkg/logger"
	"github.com/boxboxjason/jukebox/pkg/utils/httputils"
	"golang.org/x/crypto/bcrypt"
)

// HashString hashes a password using bcrypt
func HashString(to_hash string) (string, error) {
	hashed_string, err := bcrypt.GenerateFromPassword([]byte(to_hash), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed_string), nil
}

// CompareHashAndPassword compares a string with a hashed string
func CompareHashAndString(hashed_string string, raw_string string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed_string), []byte(raw_string)) == nil
}

// GenerateToken generates a cryptographically secure token
func GenerateToken() ([]byte, error) {
	token := make([]byte, 64)
	_, err := rand.Read(token)
	if err != nil {
		return []byte{}, err
	}
	return token, nil
}

func GenerateHashedToken() (string, string, error) {
	token, err := GenerateToken()
	if err != nil {
		logger.Error("Unable to generate token", err)
		return "", "", httputils.NewInternalServerError("Unable to generate token")
	}
	hashed_token, err := HashString(string(token))
	if err != nil {
		logger.Error("Unable to hash token", err)
		return "", "", httputils.NewInternalServerError("Unable to hash token")
	}
	return string(token), hashed_token, nil
}
