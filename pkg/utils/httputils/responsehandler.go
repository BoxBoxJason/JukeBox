package httputils

import (
	"encoding/json"
	"net/http"
)

func SendJSONResponse(w http.ResponseWriter, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		SendErrorToClient(w, NewInternalServerError("Failed to encode response"))
	}
}

func SendSuccessResponse(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusOK)
	if message != "" {
		SendJSONResponse(w, map[string]interface{}{
			"message": message,
			"status":  "success",
		})
	}
}

func SetSecureCookie(w http.ResponseWriter, name string, value string, cookie_path string, expires_hours int) {
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		HttpOnly: true,
		Secure:   true,
		Path:     cookie_path,
		SameSite: http.SameSiteNoneMode,
	}
	if expires_hours > 0 {
		cookie.MaxAge = expires_hours * 60 * 60
	} else if expires_hours < 0 {
		cookie.MaxAge = -1
	}

	http.SetCookie(w, &cookie)
}
