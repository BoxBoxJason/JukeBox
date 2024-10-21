package api

import (
	"net/http"

	"github.com/boxboxjason/jukebox/internal/middlewares"
	"github.com/boxboxjason/jukebox/pkg/utils/httputils"
	"github.com/go-chi/chi/v5"
)

const (
	MESSAGES_PREFIX = "/messages"
)

func SetupMessagesRoutes(r chi.Router) {
	messages_subrouter := chi.NewRouter()

	// Unauthenticated routes
	r.Get(MESSAGES_PREFIX, GetMessages)
	r.Get(MESSAGES_PREFIX+"/{id}", GetMessage)

	// Authenticated routes
	messages_subrouter.Use(middlewares.AuthMiddleware)
	messages_subrouter.Post("/", CreateMessage)
	messages_subrouter.Delete("/{id}", DeleteMessage)
	messages_subrouter.Delete("/", DeleteMessages)

	r.Mount(MESSAGES_PREFIX, messages_subrouter)
}

func GetMessages(w http.ResponseWriter, r *http.Request) {
	httputils.SendErrorToClient(w, httputils.NewNotImplementedError("route not implemented yet"))
}

func CreateMessage(w http.ResponseWriter, r *http.Request) {
	httputils.SendErrorToClient(w, httputils.NewNotImplementedError("route not implemented yet"))
}

func GetMessage(w http.ResponseWriter, r *http.Request) {
	httputils.SendErrorToClient(w, httputils.NewNotImplementedError("route not implemented yet"))
}

func DeleteMessage(w http.ResponseWriter, r *http.Request) {
	httputils.SendErrorToClient(w, httputils.NewNotImplementedError("route not implemented yet"))
}

func DeleteMessages(w http.ResponseWriter, r *http.Request) {
	httputils.SendErrorToClient(w, httputils.NewNotImplementedError("route not implemented yet"))
}
