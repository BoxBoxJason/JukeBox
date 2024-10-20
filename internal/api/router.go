package api

import "github.com/go-chi/chi/v5"

// ApiRouter creates the API router with all the routes
func ApiRouter() chi.Router {
	api_router := chi.NewRouter()
	return api_router
}
