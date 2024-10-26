package cryptutils

import "testing"

func TestHashString(t *testing.T) {
	// Test the HashString function
	hashed_string, err := HashString("password")
	if err != nil {
		t.Errorf("Error hashing string: %v", err)
	}
	if len(hashed_string) == 0 {
		t.Error("Hashed string is empty")
	}
}

func TestCompareHashAndString(t *testing.T) {
	// Test the CompareHashAndString function
	hashed_string, err := HashString("password")
	if err != nil {
		t.Errorf("Error hashing string: %v", err)
	}
	if !CompareHashAndString(hashed_string, "password") {
		t.Error("Hashed string does not match original string")
	}
}

func TestGenerateToken(t *testing.T) {
	// Test the GenerateToken function
	token, err := GenerateToken()
	if err != nil {
		t.Errorf("Error generating token: %v", err)
	}
	if len(token) == 0 {
		t.Error("Token is empty")
	}
}

func TestGenerateHashedToken(t *testing.T) {
	// Test the GenerateHashedToken function
	token, hashed_token, err := GenerateHashedToken()
	if err != nil {
		t.Errorf("Error generating hashed token: %v", err)
	}
	if len(token) == 0 {
		t.Error("Token is empty")
	}
	if len(hashed_token) == 0 {
		t.Error("Hashed token is empty")
	}

	if !CompareHashAndString(hashed_token, token) {
		t.Error("Hashed token does not match original token")
	}
}
