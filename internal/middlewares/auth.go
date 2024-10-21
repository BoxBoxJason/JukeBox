package middlewares

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/boxboxjason/jukebox/internal/constants"
	db_model "github.com/boxboxjason/jukebox/internal/model"
	"github.com/boxboxjason/jukebox/pkg/utils/httputils"
)

type contextKey string

const UserContextKey contextKey = "user"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the user ID and access token from the request
		user_id, access_token, err := getUserIDAndAccessToken(r)
		if err != nil {
			httputils.SendErrorToClient(w, err)
			return
		}

		// Open a connection to the database
		db, err := db_model.OpenConnection()
		if err != nil {
			httputils.SendErrorToClient(w, err)
			return
		}
		defer db_model.CloseConnection(db)

		// Check if the user exists and the access token is valid
		user, err := db_model.GetUserByID(db, user_id)
		if err != nil {
			httputils.SendErrorToClient(w, err)
			return
		}

		_, err = user.CheckAuthTokenMatchesByType(db, access_token, db_model.ACCESS_TOKEN)
		if err == nil {
			ctx := context.WithValue(r.Context(), UserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			httputils.SendErrorToClient(w, httputils.NewUnauthorizedError("Invalid access token"))
			return
		}
	})
}

func getUserIDAndAccessToken(r *http.Request) (int, string, error) {
	identity, err := httputils.RetrieveAuthorizationToken(r, constants.AUTH_SCHEME+" ")
	if err != nil {
		return -1, "", err
	}
	parts := strings.Split(identity, ":")
	if len(parts) != 2 {
		return -1, "", httputils.NewUnauthorizedError("Invalid identity format")
	}
	user_id, err := strconv.Atoi(parts[0])
	if err != nil || user_id < 0 {
		return -1, "", httputils.NewUnauthorizedError("Invalid user ID")
	}
	return user_id, parts[1], nil
}
