package template

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/ckshitij/notify-srv/internal/shared"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service TemplateService
}

func NewHandler(s TemplateService) *Handler {
	return &Handler{service: s}
}

// Only for user template, system templates are created via migrations
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
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
		Type:        shared.UserTemplate,
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

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	templateIDStr := chi.URLParam(r, "id")
	templateID, err := strconv.ParseInt(templateIDStr, 10, 64)
	if err != nil || templateID <= 0 {
		http.Error(w, "invalid template ID ", http.StatusBadRequest)
		return
	}

	out, err := h.service.GetByID(r.Context(), templateID)
	if err != nil {
		http.Error(w, err.Error(), shared.ErrorHttpMapper(err))
		return
	}

	shared.WriteJSON(w, http.StatusOK, out)
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

func (h *Handler) Render(w http.ResponseWriter, r *http.Request) {
	templateIDStr := chi.URLParam(r, "id")
	templateID, err := strconv.ParseInt(templateIDStr, 10, 64)
	if err != nil || templateID <= 0 {
		http.Error(w, "invalid template ID ", http.StatusBadRequest)
		return
	}

	var req RenderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	out, err := h.service.Render(r.Context(), templateID, req.TemplateKeyValue)
	if err != nil {
		http.Error(w, err.Error(), shared.ErrorHttpMapper(err))
		return
	}

	shared.WriteJSON(w, http.StatusOK, out)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	var templateFilter = parseTemplateFilters(q)

	result, err := h.service.List(r.Context(), templateFilter)
	if err != nil {
		http.Error(w, err.Error(), shared.ErrorHttpMapper(err))
		return
	}

	shared.WriteJSON(w, http.StatusOK, result)
}

func (h *Handler) CacheReloadSystemTemplates(w http.ResponseWriter, r *http.Request) {
	err := h.service.CacheReloadSystemTemplates(r.Context())
	if err != nil {
		http.Error(w, err.Error(), shared.ErrorHttpMapper(err))
		return
	}

	shared.WriteJSON(w, http.StatusNoContent, nil)
}

func parseTemplateFilters(q url.Values) TemplateFilter {
	var filter = TemplateFilter{}

	if c := q.Get("channel"); c != "" {
		ch := shared.Channel(c)
		filter.Channel = &ch
	}

	if t := q.Get("type"); t != "" {
		tt := shared.TemplateType(t)
		filter.Type = &tt
	}

	if t := q.Get("name"); t != "" {
		filter.Name = &t
	}

	var err error
	// pagination
	if l := q.Get("limit"); l != "" {
		filter.Limit, err = strconv.Atoi(l)
		if err != nil || filter.Limit < 1 {
			filter.Limit = 0
		}
	}

	if o := q.Get("offset"); o != "" {
		filter.Offset, err = strconv.Atoi(o)
		if err != nil || filter.Offset < 0 {
			filter.Offset = 0
		}
	}

	return filter
}
