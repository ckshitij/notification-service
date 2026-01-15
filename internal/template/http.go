package template

import (
	"encoding/json"
	"net/http"

	"github.com/ckshitij/notify-srv/internal/shared"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

// Only for user template, system templates are created via migrations
func (h *Handler) CreateTemplate(w http.ResponseWriter, r *http.Request) {
	var req CreateTemplateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	tpl := Template{
		Name:        req.Name,
		Description: req.Description,
		Channel:     req.Channel,
		Type:        shared.UserTemplate,
		CreatedBy:   1, // TODO: from auth
		UpdatedBy:   1,
	}

	if err := h.service.CreateTemplate(r.Context(), tpl); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) Render(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	channel := chi.URLParam(r, "channel")

	var req RenderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	out, err := h.service.Render(
		r.Context(),
		name,
		shared.UserTemplate,
		shared.Channel(channel),
		req.TemplateKeyValue,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shared.WriteJSON(w, http.StatusOK, out)
}

func (h *Handler) ListVersions(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	channel := chi.URLParam(r, "channel")

	templates, err := h.service.ListVersionsByName(
		r.Context(),
		name,
		shared.Channel(channel),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	shared.WriteJSON(w, http.StatusOK, templates)
}

func (h *Handler) AddVersion(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	channel := chi.URLParam(r, "channel")

	var req AddVersionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.AddVersionByName(
		r.Context(),
		name,
		shared.Channel(channel),
		req.Subject,
		req.Body,
	); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shared.WriteJSON(w, http.StatusNoContent, nil)
}

func (h *Handler) ListTemplatesSummary(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	var (
		channel *shared.Channel
		tplType *shared.TemplateType
	)

	if c := q.Get("channel"); c != "" {
		ch := shared.Channel(c)
		channel = &ch
	}

	if t := q.Get("type"); t != "" {
		tt := shared.TemplateType(t)
		tplType = &tt
	}

	result, err := h.service.ListTemplatesWithActiveVersion(
		r.Context(),
		channel,
		tplType,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	shared.WriteJSON(w, http.StatusOK, result)
}
