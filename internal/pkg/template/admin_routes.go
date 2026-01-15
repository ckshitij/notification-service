package template

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) AdminRoutes() http.Handler {
	r := chi.NewRouter()

	r.Get("/system/cache/reload", h.CacheReloadSystemTemplates)
	r.Get("/{id}/cache/invalidate", h.InvalidateTemplateCache)

	return r
}

func NewAdminTemplateRoutes(service TemplateService) http.Handler {
	return NewHandler(service).AdminRoutes()
}
