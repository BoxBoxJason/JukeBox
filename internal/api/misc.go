package api

import (
	"net/http"

	"github.com/boxboxjason/jukebox/pkg/utils/httputils"
	"github.com/go-chi/chi/v5"
)

func SetupMiscRoutes(r chi.Router) {
	r.Get("/health", HealthCheck)
	r.Get("/version", Version)
	r.Get("/metrics", Metrics)
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	httputils.SendErrorToClient(w, httputils.NewNotImplementedError("route not implemented yet"))
}

func Version(w http.ResponseWriter, r *http.Request) {
	httputils.SendErrorToClient(w, httputils.NewNotImplementedError("route not implemented yet"))
}

func Metrics(w http.ResponseWriter, r *http.Request) {
	httputils.SendErrorToClient(w, httputils.NewNotImplementedError("route not implemented yet"))
}
