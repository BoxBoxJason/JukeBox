package api

import (
	"net/http"

	"github.com/boxboxjason/jukebox/internal/middlewares"
	"github.com/boxboxjason/jukebox/pkg/utils/httputils"
	"github.com/go-chi/chi/v5"
)

const (
	USERS_PREFIX = "/users"
)

func SetUsersRoutes(r chi.Router) {
	users_subrouter := chi.NewRouter()
	users_subrouter.Use(middlewares.AuthMiddleware)

	// Unauthenticated routes
	r.Post(USERS_PREFIX+"/login", Login)

	// Authenticated routes
	users_subrouter.Get("/", GetUsers)
	users_subrouter.Post("/", CreateUser)
	users_subrouter.Get("/{id}", GetUser)
	users_subrouter.Put("/{id}", UpdateUser)
	users_subrouter.Delete("/{id}", DeleteUser)
	users_subrouter.Post("/logout", Logout)

	r.Mount(USERS_PREFIX, users_subrouter)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	httputils.SendErrorToClient(w, httputils.NewNotImplementedError("route not implemented yet"))
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	httputils.SendErrorToClient(w, httputils.NewNotImplementedError("route not implemented yet"))
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	httputils.SendErrorToClient(w, httputils.NewNotImplementedError("route not implemented yet"))
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	httputils.SendErrorToClient(w, httputils.NewNotImplementedError("route not implemented yet"))
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	httputils.SendErrorToClient(w, httputils.NewNotImplementedError("route not implemented yet"))
}

func Login(w http.ResponseWriter, r *http.Request) {
	httputils.SendErrorToClient(w, httputils.NewNotImplementedError("route not implemented yet"))
}

func Logout(w http.ResponseWriter, r *http.Request) {
	httputils.SendErrorToClient(w, httputils.NewNotImplementedError("route not implemented yet"))
}
