package server

import (
	"net/http"

	db "github.com/ckshitij/notification-srv/internal/db/mysql"
	"github.com/ckshitij/notification-srv/internal/logger"
	"github.com/ckshitij/notification-srv/internal/metrics"
	"github.com/ckshitij/notification-srv/internal/transport/http/template"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	BasePath     = "/notifications-srv/api"
	InternalPath = "/internal"
)

func NewRouter(log logger.Logger, database *db.DB) http.Handler {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(AccessLogMiddleware(log))

	// Public APIs
	r.Route(BasePath, func(r chi.Router) {
		r.Mount("/v1/templates", template.NewTemplateRoutes(database))
	})

	// Internal / infra APIs
	r.Route(InternalPath, func(r chi.Router) {
		r.Get("/health", LivenessHandler)
		r.Get("/ready", ReadinessHandler(database))
		r.Get("/metrics", metrics.PromHandler().ServeHTTP)
	})

	return r
}
