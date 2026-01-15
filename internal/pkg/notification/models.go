package notification

import (
	"time"

	"github.com/ckshitij/notify-srv/internal/shared"
)

/*
pending → sending → sent
pending → scheduled → sending → sent
sending → failed
*/

type NotificationStatus string

const (
	StatusPending   NotificationStatus = "pending"
	StatusScheduled NotificationStatus = "scheduled"
	StatusSending   NotificationStatus = "sending"
	StatusSent      NotificationStatus = "sent"
	StatusFailed    NotificationStatus = "failed"
)

type NotificationRecipient struct {
	Email     *string `json:"email,omitempty"`
	SlackUser *string `json:"slack,omitempty"`
	InAppUser *string `json:"in_app,omitempty"`
}

type Notification struct {
	ID               int64                 `json:"id"`
	Channel          shared.Channel        `json:"channel"`
	TemplateID       int64                 `json:"template_id"`
	Recipient        NotificationRecipient `json:"recipient"`
	TemplateKeyValue map[string]any        `json:"template_key_value"`
	Status           NotificationStatus    `json:"status"`
	ScheduledAt      *time.Time            `json:"scheduled_at,omitempty"`
	SentAt           *time.Time            `json:"sent_at,omitempty"`
	CreatedAt        time.Time             `json:"created_at"`
	UpdatedAt        time.Time             `json:"updated_at"`
}

type NotificationFilter struct {
	Channel *shared.Channel
	Status  *NotificationStatus
}

type SendNowRequest struct {
	Channel shared.Channel `json:"channel"`

	TemplateID int64 `json:"template_id"`

	Recipient        map[string]string `json:"recipient"`
	TemplateKeyValue map[string]any    `json:"template_key_value"`
}

type ScheduleRequest struct {
	SendNowRequest
	ScheduledAt time.Time `json:"scheduled_at"`
}

type NotificationResponse struct {
	ID     int64  `json:"id"`
	Status string `json:"status"`
}
