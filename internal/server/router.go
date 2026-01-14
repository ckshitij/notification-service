package server

import (
	"net/http"

	"github.com/ckshitij/notify-srv/internal/logger"
	"github.com/ckshitij/notify-srv/internal/metrics"
	"github.com/ckshitij/notify-srv/internal/repository/mysql"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	BasePath     = "/notify-srv/api"
	InternalPath = "/internal"
)

func NewRouter(log logger.Logger, database *mysql.DB, openAPIPath string, modRoutes map[string]http.Handler) http.Handler {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(AccessLogMiddleware(log))

	// Public APIs
	r.Route(BasePath, func(r chi.Router) {
		for pathStr, routes := range modRoutes {
			r.Mount(pathStr, routes)
		}
	})

	// Swagger UI
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/docs"),
	))

	r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, openAPIPath)
	})

	// Internal / infra APIs
	r.Route(InternalPath, func(r chi.Router) {
		r.Get("/health", LivenessHandler)
		r.Get("/ready", ReadinessHandler(database))
		r.Get("/metrics", metrics.PromHandler().ServeHTTP)
	})

	return r
}
