package template

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()

	// Collection-level operations
	r.Post("/", h.Create)

	r.Get("/", h.List)
	r.Get("/{id}", h.GetByID)
	r.Post("/{id}/render", h.Render)

	return r
}

func NewTemplateRoutes(service TemplateService) http.Handler {
	return NewHandler(service).Routes()
}
