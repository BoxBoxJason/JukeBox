package middlewares

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/boxboxjason/jukebox/internal/constants"
	db_model "github.com/boxboxjason/jukebox/internal/model"
	"github.com/boxboxjason/jukebox/pkg/utils/httputils"
)

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
		user, err := db_model.GetUserByID(db.Preload("Tokens").Preload("Bans"), user_id)
		if err != nil {
			httputils.SendErrorToClient(w, err)
			return
		}

		// Check if the user is under a current ban
		bans, err := user.GetActiveBans(db)
		if err != nil {
			httputils.SendErrorToClient(w, err)
			return
		} else if len(bans) > 0 {
			httputils.SendErrorToClient(w, httputils.NewForbiddenError(fmt.Sprintf("User is banned until %s for reason: %s", bans[0].EndsAt, bans[0].Reason)))
			return
		}

		// Check if the access token matches the one stored in the database
		db_access_token, err := user.CheckAuthTokenMatchesByType(db, access_token, constants.ACCESS_TOKEN)
		if err == nil {
			// Attach the user to the request context
			ctx := context.WithValue(r.Context(), constants.USER_CONTEXT_KEY, user)
			// Attach the access token to the request context
			ctx = context.WithValue(ctx, constants.ACCESS_TOKEN_CONTEXT_KEY, db_access_token)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			httputils.SendErrorToClient(w, httputils.NewUnauthorizedError("Invalid access token"))
			return
		}
	})
}

func AdminAuthMiddleware(next http.Handler) http.Handler {
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
		user, err := db_model.GetUserByID(db.Preload("Tokens").Preload("Bans"), user_id)
		if err != nil {
			httputils.SendErrorToClient(w, err)
			return
		}

		// Check if the user is under a current ban
		bans, err := user.GetActiveBans(db)
		if err != nil {
			httputils.SendErrorToClient(w, err)
			return
		} else if len(bans) > 0 {
			httputils.SendErrorToClient(w, httputils.NewForbiddenError(fmt.Sprintf("User is banned until %s for reason: %s", bans[0].EndsAt, bans[0].Reason)))
			return
		}

		// Check if the access token matches the one stored in the database
		db_access_token, err := user.CheckAuthTokenMatchesByType(db, access_token, constants.ACCESS_TOKEN)
		if err == nil {
			if user.Admin {
				// Attach the user to the request context
				ctx := context.WithValue(r.Context(), constants.USER_CONTEXT_KEY, user)
				// Attach the access token to the request context
				ctx = context.WithValue(ctx, constants.ACCESS_TOKEN_CONTEXT_KEY, db_access_token)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				httputils.SendErrorToClient(w, httputils.NewForbiddenError("User not authorized"))
				return
			}
		} else {
			httputils.SendErrorToClient(w, httputils.NewUnauthorizedError("Invalid access token"))
			return
		}
	})
}

func getUserIDAndAccessToken(r *http.Request) (int, string, error) {
	identity_bearer, err := readAccessCookie(r)
	if err != nil {
		identity_bearer, err = httputils.RetrieveAuthorizationToken(r, constants.AUTH_SCHEME)
		if err != nil {
			return -1, "", err
		}
	}

	user_id, access_token, err := DecodeIdentityBearerToUserAndToken(identity_bearer)
	if err != nil {
		return -1, "", err
	}
	return user_id, access_token, nil
}

func GetUserIDAndRefreshToken(r *http.Request) (int, string, error) {
	identity_bearer, err := readRefreshCookie(r)
	if err != nil {
		identity_bearer, err = httputils.RetrieveAuthorizationToken(r, constants.AUTH_SCHEME+" ")
		if err != nil {
			return -1, "", err
		}
	}

	user_id, refresh_token, err := DecodeIdentityBearerToUserAndToken(identity_bearer)
	if err != nil {
		return -1, "", err
	}
	return user_id, refresh_token, nil
}

func readAccessCookie(r *http.Request) (string, error) {
	cookie, err := httputils.ReadCookie(r, constants.ACCESS_TOKEN_COOKIE_NAME)
	if err != nil {
		return "", httputils.NewUnauthorizedError("access token not found")
	}
	return cookie, nil
}

func readRefreshCookie(r *http.Request) (string, error) {
	cookie, err := httputils.ReadCookie(r, constants.REFRESH_TOKEN_COOKIE_NAME)
	if err != nil {
		return "", httputils.NewUnauthorizedError("refresh token not found")
	}
	return cookie, nil
}

func DecodeIdentityBearerToUserAndToken(identity_bearer string) (int, string, error) {
	// Decode the base64 encoded identity bearer
	decoded_bearer, err := base64.RawStdEncoding.DecodeString(identity_bearer)

	if err != nil {
		return -1, "", httputils.NewUnauthorizedError("Invalid identity bearer")
	}

	parts := strings.Split(string(decoded_bearer), ":")
	if len(parts) != 2 {
		return -1, "", httputils.NewUnauthorizedError("Invalid identity format")
	}

	user_id, err := strconv.Atoi(parts[0])
	if err != nil || user_id < 0 {
		return -1, "", httputils.NewUnauthorizedError("Invalid user ID")
	}
	return user_id, parts[1], nil
}

func EncodeUserAndTokenToIdentityBearer(user_id int, access_token string) string {
	// Encode the user ID and access token to base64
	return base64.StdEncoding.EncodeToString([]byte(strconv.Itoa(user_id) + ":" + access_token))
}
