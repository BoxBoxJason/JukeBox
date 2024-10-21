package api

import "github.com/go-chi/chi/v5"

// ApiRouter creates the API router with all the routes
func ApiRouter() chi.Router {
	api_router := chi.NewRouter()

	// Setup the subrouters
	SetupMessagesRoutes(api_router)
	SetUsersRoutes(api_router)
	SetupMiscRoutes(api_router)

	return api_router
}
