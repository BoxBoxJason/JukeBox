package api

import (
	"net/http"

	"github.com/boxboxjason/jukebox/internal/constants"
	db_model "github.com/boxboxjason/jukebox/internal/model"
	"github.com/boxboxjason/jukebox/pkg/utils/httputils"
	"github.com/go-chi/chi/v5"
)

func SetupMiscRoutes(r chi.Router) {
	r.Get("/health", HealthCheck)
	r.Get("/version", Version)
	r.Get("/metrics", Metrics)
}

// HealthCheck is the handler for the health check route
// Returns the status of the API, of the database, and of the services
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	err := db_model.CheckDatabase()
	if err != nil {
		httputils.SendJSONResponse(w, map[string]string{
			"database": "unhealthy",
			"api":      "healthy",
			"service":  "healthy",
		})
	} else {
		httputils.SendJSONResponse(w, map[string]string{
			"database": "healthy",
			"api":      "healthy",
			"service":  "healthy",
		})
	}
}

func Version(w http.ResponseWriter, r *http.Request) {
	httputils.SendJSONResponse(w, map[string]string{
		"version":  constants.JUKEBOX_VERSION,
		"database": db_model.DatabaseVersion(),
	})
}

func Metrics(w http.ResponseWriter, r *http.Request) {
	httputils.SendErrorToClient(w, httputils.NewNotImplementedError("route not implemented yet"))
}
