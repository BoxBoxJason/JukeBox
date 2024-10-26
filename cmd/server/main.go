package main

import (
	"net/http"
	"path"

	"github.com/boxboxjason/jukebox/internal/api"
	"github.com/boxboxjason/jukebox/internal/constants"
	"github.com/boxboxjason/jukebox/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Setup the logger
	logger.SetupLogger(constants.LOG_DIR, "DEBUG")

	// Create new main router
	main_router := chi.NewRouter()

	// Setup middlewares
	main_router.Use(middleware.Logger)    // Log every HTTP request
	main_router.Use(middleware.Recoverer) // Recover from panics
	main_router.Use(middleware.RealIP)    // Get the real IP address of the client
	main_router.Use(middleware.RequestID) // Generate a request ID for every request (might delete that later)

	// Serve the frontend
	frontend_fs := http.FileServer(http.Dir(path.Join(".", "frontend", "dist")))
	main_router.Handle("/*", frontend_fs)
	logger.Debug("Serving frontend from ", path.Join(".", "frontend", "dist"))

	// Serve the API
	api_router := api.ApiRouter()
	main_router.Mount("/api", api_router)

	// Start the server (attempt to use TLS first)
	logger.Info("Starting JukeBox server at https://localhost:3000")
	err := http.ListenAndServeTLS(":3000", "secret/cert.pem", "secret/key.pem", main_router)
	if err != nil {
		logger.Critical("Unable to start the server using TLS: ", err)
		logger.Info("Starting the server without TLS at http://localhost:3000")
		err = http.ListenAndServe(":3000", main_router)
		if err != nil {
			logger.Fatal("Unable to start the server", err)
		}
	}
}
