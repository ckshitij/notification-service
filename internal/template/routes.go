package template

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()

	// Collection-level operations
	r.Post("/", h.CreateTemplate)

	r.Get("/summary", h.ListTemplatesSummary)

	// Resource-level operations
	r.Route("/{channel}/{name}", func(r chi.Router) {

		// Render / preview template (representation)
		r.Post("/", h.Render)

		// Versions sub-resource
		r.Route("/versions", func(r chi.Router) {
			r.Get("/", h.ListVersions)
			r.Post("/", h.AddVersion)
		})
	})

	return r
}

func NewTemplateRoutes(repo Repository) http.Handler {
	renderer := NewGoTemplateRenderer()
	service := NewService(repo, renderer)
	return NewHandler(service).Routes()
}
