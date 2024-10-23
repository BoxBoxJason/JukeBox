package api

import (
	"net/http"

	"github.com/boxboxjason/jukebox/internal/constants"
	db_controller "github.com/boxboxjason/jukebox/internal/controller"
	"github.com/boxboxjason/jukebox/internal/middlewares"
	db_model "github.com/boxboxjason/jukebox/internal/model"
	"github.com/boxboxjason/jukebox/pkg/logger"
	"github.com/boxboxjason/jukebox/pkg/utils/httputils"
	"github.com/go-chi/chi/v5"
)

const (
	MESSAGES_PREFIX = "/messages"
)

func SetupMessagesRoutes(r chi.Router) {
	messages_subrouter := chi.NewRouter()

	// Unauthenticated routes
	messages_subrouter.Get("/", GetMessages)
	messages_subrouter.Get(ID_PARAM_ENDPOINT, GetMessage)

	// Authenticated routes
	messages_subrouter.Group(func(auth_router chi.Router) {
		auth_router.Use(middlewares.AuthMiddleware)
		auth_router.Post("/", CreateMessage)
		auth_router.Delete(ID_PARAM_ENDPOINT, DeleteMessage)
		auth_router.Delete("/", DeleteMessages)
	})

	r.Mount(MESSAGES_PREFIX, messages_subrouter)
}

// ==================== CRUD operations ====================

// ==================== Create ====================

// CreateMessage creates a message and adds it to the database
func CreateMessage(w http.ResponseWriter, r *http.Request) {
	// Retrieve the user from the context
	user, ok := r.Context().Value(constants.USER_CONTEXT_KEY).(*db_model.User)
	if !ok {
		httputils.SendErrorToClient(w, httputils.NewUnauthorizedError("user not found"))
		return
	}

	// Retrieve the message content
	message_content, err := httputils.RetrieveStringParameter(r, "message", false)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Create the message
	message, err := db_controller.CreateMessage(message_content, user)
	if err != nil {
		logger.Error("Failed to create message", err)
		httputils.SendErrorToClient(w, httputils.NewInternalServerError("Failed to create message"))
		return
	}

	// Send the message to the client
	httputils.SendJSONResponse(w, message)
}

// ==================== Read ====================

// GetMessages retrieves messages depending on the query parameters
func GetMessages(w http.ResponseWriter, r *http.Request) {
	httputils.SendErrorToClient(w, httputils.NewNotImplementedError("route not implemented yet"))
}

// GetMessage retrieves a message by its ID
func GetMessage(w http.ResponseWriter, r *http.Request) {
	message_id, err := httputils.RetrieveChiIntArgument(r, ID_PARAM)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Open db connection
	db, err := db_model.OpenConnection()
	if err != nil {
		logger.Error("Failed to open database connection", err)
		httputils.SendErrorToClient(w, httputils.NewInternalServerError("Failed to open database connection"))
		return
	}
	defer db_model.CloseConnection(db)

	// Retrieve the message
	message, err := db_model.GetMessageByID(db, message_id)
	if err != nil {
		logger.Error("Failed to retrieve message", err)
		httputils.SendErrorToClient(w, httputils.NewNotFoundError("Message not found"))
		return
	}

	// Send the message to the client
	httputils.SendJSONResponse(w, message)
}

// ==================== Update ====================

// ==================== Delete ====================

// DeleteMessage deletes a message by its ID
func DeleteMessage(w http.ResponseWriter, r *http.Request) {
	message_id, err := httputils.RetrieveChiIntArgument(r, ID_PARAM)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Open db connection
	db, err := db_model.OpenConnection()
	if err != nil {
		logger.Error("Failed to open database connection", err)
		httputils.SendErrorToClient(w, httputils.NewInternalServerError("Failed to open database connection"))
		return
	}
	defer db_model.CloseConnection(db)

	// Retrieve the message
	message, err := db_model.GetMessageByID(db, message_id)
	if err != nil {
		logger.Error("Failed to retrieve message", err)
		httputils.SendErrorToClient(w, httputils.NewNotFoundError("Message not found"))
		return
	}

	// Check if the user has permission to delete the message
	user, ok := r.Context().Value(constants.USER_CONTEXT_KEY).(*db_model.User)
	if !ok {
		httputils.SendErrorToClient(w, httputils.NewUnauthorizedError("user not found"))
		return
	}

	if !db_controller.UserHasPermissionToDeleteMessage(user, message) {
		httputils.SendErrorToClient(w, httputils.NewUnauthorizedError("User does not have permission to delete message"))
		return
	}

	// Delete the message
	err = message.DeleteMessage(db)
	if err != nil {
		logger.Error("Failed to delete message", err)
		httputils.SendErrorToClient(w, httputils.NewInternalServerError("Failed to delete message"))
		return
	}

	// Send a success response
	httputils.SendSuccessResponse(w, "")
}

// DeleteMessages deletes messages depending on the query parameters
func DeleteMessages(w http.ResponseWriter, r *http.Request) {
	httputils.SendErrorToClient(w, httputils.NewNotImplementedError("route not implemented yet"))
}
