package main

import (
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/boxboxjason/jukebox/internal/api"
	"github.com/boxboxjason/jukebox/internal/websocket"
	"github.com/boxboxjason/jukebox/internal/constants"
	"github.com/boxboxjason/jukebox/internal/jobs"
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
	frontend_dir := path.Join(".", "frontend", "dist")
	frontend_fs := http.FileServer(http.Dir(frontend_dir))
	main_router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		requested_path := filepath.Join(frontend_dir, r.URL.Path)
		if _, err := os.Stat(requested_path); os.IsNotExist(err) || r.URL.Path == "/" {
			// Serve index.html for unmatched routes or root path
			http.ServeFile(w, r, filepath.Join(frontend_dir, "index.html"))
		} else {
			// Serve static file
			frontend_fs.ServeHTTP(w, r)
		}
	})
	logger.Info("Serving frontend at \"/\" from", frontend_dir)

	// Serve the API
	api_router := api.ApiRouter()
	main_router.Mount("/api", api_router)
	logger.Info("Serving API at /api")

	// Serve WebSocket
	main_router.HandleFunc("/ws/chat", chatwebsocket.ChatWebSocket)
	logger.Info("Serving WebSocket chat at /ws/chat")

	// Start jobs
	jobs.SetupJobs()
	logger.Info("Jobs started")

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