package cryptutils

import "golang.org/x/crypto/bcrypt"

// HashString hashes a password using bcrypt
func HashString(to_hash string) (string, error) {
	hashed_string, err := bcrypt.GenerateFromPassword([]byte(to_hash), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed_string), nil
}

// CompareHashAndPassword compares a password with a hashed password
func CompareHashAndString(hashed_string string, raw_string string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed_string), []byte(raw_string)) == nil
}
