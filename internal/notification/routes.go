package notification

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Post("/", h.SendNow)
	r.Post("/schedule", h.Schedule)
	r.Get("/{id}", h.GetByID)
	r.Get("/", h.List)

	return r
}

func NewNotificationRoutes(service Service) http.Handler {
	handler := NewHandler(service)
	return handler.Routes()
}
