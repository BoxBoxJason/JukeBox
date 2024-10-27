package api

import (
	"github.com/boxboxjason/jukebox/internal/constants"
	"github.com/go-chi/chi/v5"
)

const (
	ID_PARAM_ENDPOINT = "/{" + constants.ID_PARAM + "}"
)

// ApiRouter creates the API router with all the routes
func ApiRouter() chi.Router {
	api_router := chi.NewRouter()

	// Setup the subrouters
	SetupMessagesRoutes(api_router)
	SetUsersRoutes(api_router)
	SetupMiscRoutes(api_router)
	SetupAuthRoutes(api_router)

	return api_router
}
