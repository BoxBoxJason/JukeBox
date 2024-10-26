package middlewares

import (
	"testing"
)

func TestTokenEncodeDecode(t *testing.T) {
	token := "test_token"
	user_id := 2205
	encoded_token := EncodeUserAndTokenToIdentityBearer(user_id, token)
	decoded_user_id, decoded_token, err := DecodeIdentityBearerToUserAndToken(encoded_token)
	if err != nil {
		t.Errorf("Error decoding token: %v", err)
	}

	if decoded_user_id != user_id {
		t.Errorf("Expected user_id %d, got %d", user_id, decoded_user_id)
	}

	if decoded_token != token {
		t.Errorf("Expected token %s, got %s", token, decoded_token)
	}
}

func TestTokenEncodeDecodeInvalidToken(t *testing.T) {
	encoded_token := "thiswontwork"
	_, _, err := DecodeIdentityBearerToUserAndToken(encoded_token)
	if err == nil {
		t.Errorf("Expected error decoding token")
	}
}
