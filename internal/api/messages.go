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
		auth_router.Patch(ID_PARAM_ENDPOINT, UpdateMessage)
	})

	// Admin routes
	messages_subrouter.Group(func(admin_router chi.Router) {
		admin_router.Use(middlewares.AdminAuthMiddleware)
		admin_router.Delete(ID_PARAM_ENDPOINT, DeleteMessage)
		admin_router.Delete("/", DeleteMessages)
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
	message_content, err := httputils.RetrieveStringParameter(r, constants.MESSAGE_PARAMETER, false)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Create the message
	message, err := db_controller.CreateMessage(nil, message_content, user)
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
	ids, err := httputils.RetrieveIntListValueParameter(r, constants.ID_PARAMETER, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	sender_ids, err := httputils.RetrieveIntListValueParameter(r, constants.SENDER_ID_PARAMETER, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	flagged, err := httputils.RetrieveBoolParameter(r, constants.FLAGGED_PARAMETER, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	censored, err := httputils.RetrieveBoolParameter(r, constants.CENSORED_PARAMETER, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	removed, err := httputils.RetrieveBoolParameter(r, constants.REMOVED_PARAMETER, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	contains, err := httputils.RetrieveStringListValueParameter(r, constants.CONTAINS_PARAMETER, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Retrieve the messages
	messages, err := db_controller.GetMessages(nil, ids, sender_ids, flagged, censored, removed, contains)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Send the messages to the client
	httputils.SendJSONResponse(w, messages)
}

// GetMessage retrieves a message by its ID
func GetMessage(w http.ResponseWriter, r *http.Request) {
	message_id, err := httputils.RetrieveChiIntArgument(r, constants.ID_PARAMETER)
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

func UpdateMessage(w http.ResponseWriter, r *http.Request) {
	message_id, err := httputils.RetrieveChiIntArgument(r, constants.ID_PARAMETER)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Retrieve the user from the context
	user, ok := r.Context().Value(constants.USER_CONTEXT_KEY).(*db_model.User)
	if !ok {
		httputils.SendErrorToClient(w, httputils.NewUnauthorizedError("user not found"))
		return
	}

	// Retrieve the message content
	message_content, err := httputils.RetrieveStringParameter(r, constants.MESSAGE_PARAMETER, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Retrieve the possibly updated censor status
	censored, err := httputils.RetrieveBoolParameter(r, constants.CENSORED_PARAMETER, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Retrieve the possibly updated flagged status
	flagged, err := httputils.RetrieveBoolParameter(r, constants.FLAGGED_PARAMETER, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Retrieve the possibly updated removed status
	removed, err := httputils.RetrieveBoolParameter(r, constants.REMOVED_PARAMETER, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	if (len(censored) > 0 || len(flagged) > 0 || len(removed) > 0) && !user.Admin {
		httputils.SendErrorToClient(w, httputils.NewForbiddenError("only admins can update the censor, flagged and removed status"))
		return
	}

	// Open db connection
	db, err := db_model.OpenConnection()
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}
	defer db_model.CloseConnection(db)

	// Retrieve the message
	message, err := db_controller.GetMessage(message_id)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	if message.SenderID != user.ID && len(message_content) > 0 {
		httputils.SendErrorToClient(w, httputils.NewForbiddenError("only the sender can update the message content"))
		return
	}

	// Update the message
	message, err = db_controller.UpdateMessage(nil, message, message_content, censored, flagged, removed)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Send the message to the client
	httputils.SendJSONResponse(w, message)
}

// ==================== Delete ====================

// DeleteMessage deletes a message by its ID
func DeleteMessage(w http.ResponseWriter, r *http.Request) {
	message_id, err := httputils.RetrieveChiIntArgument(r, constants.ID_PARAMETER)
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

	// Delete the message
	err = message.DeleteMessage(db)
	if err != nil {
		logger.Error("Failed to delete message", err)
		httputils.SendErrorToClient(w, httputils.NewInternalServerError("Failed to delete message"))
		return
	}

	// Send a success response
	httputils.SendSuccessResponse(w, "deleted message successfully")
}

// DeleteMessages deletes messages depending on the query parameters
func DeleteMessages(w http.ResponseWriter, r *http.Request) {
	ids, err := httputils.RetrieveIntListValueParameter(r, constants.ID_PARAMETER, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	sender_ids, err := httputils.RetrieveIntListValueParameter(r, constants.SENDER_ID_PARAMETER, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	flagged, err := httputils.RetrieveBoolParameter(r, constants.FLAGGED_PARAMETER, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	censored, err := httputils.RetrieveBoolParameter(r, constants.CENSORED_PARAMETER, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	removed, err := httputils.RetrieveBoolParameter(r, constants.REMOVED_PARAMETER, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	contains, err := httputils.RetrieveStringListValueParameter(r, constants.CONTAINS_PARAMETER, true)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Delete the messages
	err = db_controller.DeleteMessages(nil, ids, sender_ids, flagged, censored, removed, contains)
	if err != nil {
		httputils.SendErrorToClient(w, err)
		return
	}

	// Send a success response
	httputils.SendSuccessResponse(w, "deleted messages successfully")
}
