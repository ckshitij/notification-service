package notification

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/ckshitij/notify-srv/internal/shared"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) SendNow(w http.ResponseWriter, r *http.Request) {
	var req SendNowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	n, err := mapRequestToNotification(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.service.SendNow(r.Context(), n)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	shared.WriteJSON(w, http.StatusAccepted, NotificationResponse{
		ID:     id,
		Status: string(n.Status),
	})
}

func (h *Handler) Schedule(w http.ResponseWriter, r *http.Request) {
	var req ScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	n, err := mapRequestToNotification(req.SendNowRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.service.Schedule(r.Context(), n, req.ScheduledAt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	shared.WriteJSON(w, http.StatusAccepted, NotificationResponse{
		ID:     id,
		Status: string(n.Status),
	})
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {

	notificationIDStr := chi.URLParam(r, "id")
	notificationID, err := strconv.ParseInt(notificationIDStr, 10, 64)
	if err != nil || notificationID <= 0 {
		http.Error(w, "invalid notification ID ", http.StatusBadRequest)
		return
	}

	notification, err := h.service.GetByID(r.Context(), notificationID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	shared.WriteJSON(w, http.StatusOK, notification)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {

	filter := NotificationFilter{}
	channel := r.URL.Query().Get("channel")
	if channel != "" {
		ch := shared.Channel(channel)
		filter.Channel = &ch
	}
	status := r.URL.Query().Get("status")
	if status != "" {
		st := NotificationStatus(status)
		filter.Status = &st
	}

	notifications, err := h.service.List(r.Context(), filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	shared.WriteJSON(w, http.StatusOK, notifications)
}

func (h *Handler) Process(w http.ResponseWriter, r *http.Request) {

	notificationIDStr := chi.URLParam(r, "id")
	notificationID, err := strconv.ParseInt(notificationIDStr, 10, 64)
	if err != nil || notificationID <= 0 {
		http.Error(w, "invalid notification ID ", http.StatusBadRequest)
		return
	}

	err = h.service.Process(r.Context(), notificationID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	shared.WriteJSON(w, http.StatusAccepted, nil)
}

func mapRequestToNotification(req SendNowRequest) (*Notification, error) {

	n := &Notification{
		Channel:           shared.Channel(req.Channel),
		TemplateVersionID: req.TemplateVersionID,
		TemplateKeyValue:  req.TemplateKeyValue,
	}

	switch req.Channel {
	case "email":
		email := req.Recipient["email"]
		if email == "" {
			return nil, errors.New("email recipient required")
		}
		n.Recipient.Email = &email

	case "slack":
		if v := req.Recipient["user"]; v != "" {
			n.Recipient.SlackUser = &v
		}
		if n.Recipient.SlackUser == nil {
			return nil, errors.New("slack user required")
		}

	case "in_app":
		user := req.Recipient["user"]
		if user == "" {
			return nil, errors.New("in_app user required")
		}
		n.Recipient.InAppUser = &user

	default:
		return nil, errors.New("unsupported channel")
	}

	return n, nil
}
