package server

import (
	"net/http"

	"github.com/ckshitij/notification-srv/internal/db"
	"github.com/ckshitij/notification-srv/internal/logger"
	"github.com/ckshitij/notification-srv/internal/metrics"
	"github.com/ckshitij/notification-srv/internal/transport/http/template"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
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

	// Swagger UI
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/docs"),
	))

	r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./swagger/swagger.yaml")
	})

	// Internal / infra APIs
	r.Route(InternalPath, func(r chi.Router) {
		r.Get("/health", LivenessHandler)
		r.Get("/ready", ReadinessHandler(database))
		r.Get("/metrics", metrics.PromHandler().ServeHTTP)
	})

	return r
}
