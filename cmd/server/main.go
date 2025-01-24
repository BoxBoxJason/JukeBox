package main

import (
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/boxboxjason/jukebox/internal/api"
	"github.com/boxboxjason/jukebox/internal/constants"
	"github.com/boxboxjason/jukebox/internal/jobs"
	"github.com/boxboxjason/jukebox/internal/middlewares"
	db_model "github.com/boxboxjason/jukebox/internal/model"
	"github.com/boxboxjason/jukebox/internal/websocket"
	"github.com/boxboxjason/jukebox/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	// Setup the logger
	logger.SetupLogger(constants.LOG_DIR, "DEBUG")

	// Create the tables in the database
	db_model.CreateTables()

	// Create new main router
	main_router := chi.NewRouter()

	// Setup middlewares
	main_router.Use(middleware.Logger)         // Log every HTTP request
	main_router.Use(middleware.Recoverer)      // Recover from panics
	main_router.Use(middleware.RealIP)         // Get the real IP address of the client
	main_router.Use(middleware.RequestID)      // Generate a request ID for every request
	main_router.Use(cors.Handler(cors.Options{ // Setup CORS
		AllowedOrigins:   []string{"http://localhost:3000", "https://localhost:3000"}, // Allow only the frontend to access the API (for now dev URL)
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
	}))

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

	// Setup WebSocket connection route (with authentication)
	main_router.Group(func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware)
		r.HandleFunc("/chat/ws", websocket.EstablishConnection)
	})
	logger.Info("Serving Chat WebSocket at /chat/ws")

	// Start jobs
	jobs.SetupJobs()
	logger.Info("Jobs started")

	// Start the server (attempt to use TLS first)
	logger.Info("Starting JukeBox server at https://localhost:3000")
	err := http.ListenAndServeTLS(":3000", "secret/cert.pem", "secret/key.pem", main_router)
	if err != nil {
		logger.Fatal("Unable to start the server using TLS: ", err)
	}
}
