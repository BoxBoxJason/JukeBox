package api

import (
	"net/http"

	"github.com/boxboxjason/jukebox/internal/constants"
	db_controller "github.com/boxboxjason/jukebox/internal/controller"
	"github.com/boxboxjason/jukebox/internal/middlewares"
	db_model "github.com/boxboxjason/jukebox/internal/model"
	"github.com/boxboxjason/jukebox/pkg/utils/httputils"
	"github.com/go-chi/chi/v5"
)

const (
	AUTH_PREFIX = "/auth"
)

func SetupAuthRoutes(r chi.Router) {
	auth_subrouter := chi.NewRouter()

	// Unauthenticated routes
	auth_subrouter.Post("/login", Login)
	auth_subrouter.Post("/refresh", Refresh)

	// Authenticated routes
	auth_subrouter.Group(func(auth_router chi.Router) {
		auth_router.Use(middlewares.AuthMiddleware)
		auth_router.Post("/logout", Logout)
	})

	r.Mount(AUTH_PREFIX, auth_subrouter)
}

// ==================== Login/Logout ====================

// Login logs in a user by checking the validity of the input fields
func Login(w http.ResponseWriter, r *http.Request) {

	auth_from_token, err := LoginFromToken(w, r)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	} else if auth_from_token {
		return
	}

	_, err = LoginFromPassword(w, r)
	if err != nil {
		httputils.SendErrorToClient(w, err)
	}
}

func LoginFromToken(w http.ResponseWriter, r *http.Request) (bool, error) {
	success := false
	// Attempt to retrieve the access token from the request cookies
	identity_bearer, err := httputils.ReadCookie(r, constants.ACCESS_TOKEN_COOKIE_NAME)
	if err != nil {
		identity_bearer, err = httputils.RetrieveAuthorizationToken(r, constants.AUTH_SCHEME+" ")
	}

	if identity_bearer != "" {
		// Attempt to login from the token
		user_id, access_token, err := middlewares.DecodeIdentityBearerToUserAndToken(identity_bearer)
		if err != nil {
			return success, err
		} else {
			user_id, username, access_token, refresh_token, err := db_controller.LoginFromToken(user_id, access_token)
			if err != nil {
				httputils.SendErrorToClient(w, err)
				return success, err
			} else {
				setAuthCookies(w, access_token, "")
				httputils.SendJSONResponse(w, map[string]interface{}{
					"username":      username,
					"user_id":       user_id,
					"access_token":  access_token,
					"refresh_token": refresh_token,
				})
				return true, nil
			}
		}
	}
	return false, nil
}

func LoginFromPassword(w http.ResponseWriter, r *http.Request) (bool, error) {
	// Cannot login from token, login from password
	username_or_email, err := httputils.RetrievePostFormStringParameter(r, "username_or_email", false)
	if err != nil {
		return false, err
	}
	password, err := httputils.RetrievePostFormStringParameter(r, "password", false)
	if err != nil {
		return false, err
	}

	user_id, username, access_token, refresh_token, err := db_controller.LoginUserFromPassword(username_or_email, password)
	if err != nil {
		return false, err
	}

	setAuthCookies(w, access_token, refresh_token)
	httputils.SendJSONResponse(w, map[string]interface{}{
		"username":      username,
		"user_id":       user_id,
		"access_token":  access_token,
		"refresh_token": refresh_token,
	})

	return true, nil
}

// Logout logs out a user by deleting the access token
func Logout(w http.ResponseWriter, r *http.Request) {
	access_token, ok := r.Context().Value(constants.ACCESS_TOKEN_CONTEXT_KEY).(*db_model.AuthToken)
	if !ok {
		httputils.SendErrorToClient(w, httputils.NewUnauthorizedError("access token not found"))
		return
	}

	err := db_controller.DeleteToken(nil, access_token)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	deleteAuthCookies(w)
	httputils.SendSuccessResponse(w, "")
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	// Retrieve the refresh token
	identity_bearer, err := httputils.ReadCookie(r, constants.REFRESH_TOKEN_COOKIE_NAME)
	if err != nil {
		identity_bearer, err = httputils.RetrieveAuthorizationToken(r, constants.AUTH_SCHEME+" ")
		if err != nil {
			httputils.SendErrorToClient(w, httputils.NewUnauthorizedError("refresh token not found"))
			return
		}
	}

	user_id, username, access_token, new_refresh_token, err := db_controller.RefreshTokens(identity_bearer)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	setAuthCookies(w, access_token, new_refresh_token)
	httputils.SendJSONResponse(w, map[string]interface{}{
		"access_token":  access_token,
		"refresh_token": new_refresh_token,
		"username":      username,
		"user_id":       user_id,
	})
}

func setAuthCookies(w http.ResponseWriter, access_token string, refresh_token string) {
	httputils.SetSecureCookie(w, constants.ACCESS_TOKEN_COOKIE_NAME, access_token, constants.ACCESS_TOKEN_COOKIE_PATH, 0)
	httputils.SetSecureCookie(w, constants.REFRESH_TOKEN_COOKIE_NAME, refresh_token, constants.REFRESH_TOKEN_COOKIE_PATH, constants.REFRESH_TOKEN_EXPIRATION)
}

func deleteAuthCookies(w http.ResponseWriter) {
	httputils.SetSecureCookie(w, constants.ACCESS_TOKEN_COOKIE_NAME, "", constants.ACCESS_TOKEN_COOKIE_PATH, -1)
	httputils.SetSecureCookie(w, constants.REFRESH_TOKEN_COOKIE_NAME, "", constants.REFRESH_TOKEN_COOKIE_PATH, -1)
}
