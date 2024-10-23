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
	USERS_PREFIX = "/users"
)

func SetUsersRoutes(r chi.Router) {
	users_subrouter := chi.NewRouter()

	// Unauthenticated routes
	users_subrouter.Post("/", CreateUser)

	// Authenticated routes
	users_subrouter.Group(func(auth_router chi.Router) {
		auth_router.Use(middlewares.AuthMiddleware)
		auth_router.Get("/", GetUsers)
		auth_router.Get(ID_PARAM_ENDPOINT, GetUser)
		auth_router.Put(ID_PARAM_ENDPOINT, UpdateUser)
		auth_router.Delete(ID_PARAM_ENDPOINT, DeleteUser)
	})

	r.Mount(USERS_PREFIX, users_subrouter)
}

// ==================== CRUD operations ====================

// ==================== Create ====================
func CreateUser(w http.ResponseWriter, r *http.Request) {
	// Check that the user is not already authenticated
	_, ok := r.Context().Value(constants.USER_CONTEXT_KEY).(*db_model.User)
	if ok {
		httputils.SendErrorToClient(w, httputils.NewForbiddenError("user already authenticated"))
		return
	}

	// Retrieve the user input
	username, err := httputils.RetrievePostFormStringParameter(r, "username", false)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}
	email, err := httputils.RetrievePostFormStringParameter(r, "email", false)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}
	password, err := httputils.RetrievePostFormStringParameter(r, "password", false)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Create the user
	db_user, err := db_controller.CreateUser(username, email, password)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Send the user back to the client
	httputils.SendJSONResponse(w, db_user)
}

// ==================== Read ====================
func GetUsers(w http.ResponseWriter, r *http.Request) {
	httputils.SendErrorToClient(w, httputils.NewNotImplementedError("route not implemented yet"))
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	user_id, err := httputils.RetrieveChiIntArgument(r, ID_PARAM)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	user, err := db_controller.GetUser(user_id)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	httputils.SendJSONResponse(w, user)
}

// ==================== Update ====================
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	httputils.SendErrorToClient(w, httputils.NewNotImplementedError("route not implemented yet"))
}

// ==================== Delete ====================
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	httputils.SendErrorToClient(w, httputils.NewNotImplementedError("route not implemented yet"))
}
