package main

import (
	"net/http"
	"path"

	"github.com/boxboxjason/jukebox/pkg/logger"
	"github.com/gorilla/mux"
)

func main() {
	main_router := mux.NewRouter()
	// Serve the frontend
	fs := http.FileServer(http.Dir(path.Join(".", "frontend", "dist")))
	main_router.PathPrefix("/").Handler(fs)
	logger.Debug("Frontend routing initialized")

	// Start the server (attempt to use TLS first)
	logger.Info("Starting JukeBox server at https://localhost:3000")
	err := http.ListenAndServeTLS(":3000", "secret/cert.pem", "secret/key.pem", main_router)
	if err != nil {
		logger.Critical("Unable to start the server using TLS: ", err)
		logger.Debug("Starting the server without TLS at http://localhost:3000")
		err = http.ListenAndServe(":3000", main_router)
		if err != nil {
			logger.Fatal("Unable to start the server", err)
		}
	}
}
