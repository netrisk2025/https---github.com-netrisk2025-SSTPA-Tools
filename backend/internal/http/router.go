package apihttp

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(version string) http.Handler {
	router := chi.NewRouter()
	router.Get("/healthz", healthHandler(version))
	router.Route("/api/v1", func(group chi.Router) {
		group.Get("/health", healthHandler(version))
	})

	return router
}
