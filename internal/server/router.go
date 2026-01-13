package server

import (
	"net/http"

	"github.com/ckshitij/notification-srv/internal/logger"
	"github.com/ckshitij/notification-srv/internal/metrics"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	BasePath     = "/api/notifications"
	InternalPath = "/internal"
)

func NewRouter(log logger.Logger) http.Handler {
	r := chi.NewRouter()

	r.Use(RequestIDMiddleware(log))
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(AccessLogMiddleware(log))

	// Base path group
	r.Route(BasePath, func(r chi.Router) {

		// Health
		r.Get("/hell", HealthHandler)

	})

	r.Route(InternalPath, func(r chi.Router) {
		// Health
		r.Get("/health", HealthHandler)
		r.Get("/metrics", metrics.PromHandler().ServeHTTP)
	})

	return r
}
