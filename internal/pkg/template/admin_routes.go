package template

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ckshitij/notify-srv/internal/shared"
	"github.com/go-chi/chi/v5"
)

func (h *Handler) AdminRoutes() http.Handler {
	r := chi.NewRouter()

	r.Get("/system/cache/reload", h.CacheReloadSystemTemplates)
	r.Get("/{id}/cache/invalidate", h.InvalidateTemplateCache)
	r.Post("/", h.CreateSystemTemplate)

	return r
}

func NewAdminTemplateRoutes(service TemplateService) http.Handler {
	return NewHandler(service).AdminRoutes()
}

// Only for user template, system templates are created via migrations
func (h *Handler) CreateSystemTemplate(w http.ResponseWriter, r *http.Request) {
	var req CreateTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		http.Error(w, err.Error(), shared.ErrorHttpMapper(err))
		return
	}

	tpl := Template{
		Name:        req.Name,
		Description: req.Description,
		Channel:     req.Channel,
		Type:        shared.SystemTemplate,
		Subject:     req.Subject,
		Body:        req.Body,
	}

	id, err := h.service.Create(r.Context(), tpl)
	if err != nil {
		http.Error(w, err.Error(), shared.ErrorHttpMapper(err))
		return
	}

	shared.WriteJSON(w, http.StatusCreated, id)
}

func (h *Handler) CacheReloadSystemTemplates(w http.ResponseWriter, r *http.Request) {
	err := h.service.CacheReloadSystemTemplates(r.Context())
	if err != nil {
		http.Error(w, err.Error(), shared.ErrorHttpMapper(err))
		return
	}

	shared.WriteJSON(w, http.StatusNoContent, nil)
}

func (h *Handler) InvalidateTemplateCache(w http.ResponseWriter, r *http.Request) {
	templateIDStr := chi.URLParam(r, "id")
	templateID, err := strconv.ParseInt(templateIDStr, 10, 64)
	if err != nil || templateID <= 0 {
		http.Error(w, "invalid template ID ", http.StatusBadRequest)
		return
	}

	err = h.service.InvalidateTemplateCache(r.Context(), templateID)
	if err != nil {
		http.Error(w, err.Error(), shared.ErrorHttpMapper(err))
		return
	}

	shared.WriteJSON(w, http.StatusNoContent, nil)
}
