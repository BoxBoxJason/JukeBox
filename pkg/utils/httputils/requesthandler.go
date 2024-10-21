package httputils

import (
	"net/http"
	"strings"
)

func RetrieveAuthorizationToken(r *http.Request, authorization_scheme string) (string, error) {
	auth_header := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth_header, authorization_scheme) {
		return "", NewUnauthorizedError("Invalid authorization scheme")
	}
	return strings.TrimPrefix(auth_header, authorization_scheme), nil
}
