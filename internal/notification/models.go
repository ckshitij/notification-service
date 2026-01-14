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
	SlackUser *string `json:"slack_user,omitempty"`
	SlackChan *string `json:"slack_channel,omitempty"`
	InAppUser *string `json:"in_app_user,omitempty"`
}

type Notification struct {
	ID                int64                 `json:"id"`
	Channel           shared.Channel        `json:"channel"`
	TemplateVersionID int64                 `json:"template_version_id"`
	Recipient         NotificationRecipient `json:"recipient"`
	Payload           map[string]any        `json:"payload"`
	Status            NotificationStatus    `json:"status"`
	ScheduledAt       *time.Time            `json:"scheduled_at,omitempty"`
	SentAt            *time.Time            `json:"sent_at,omitempty"`
	CreatedAt         time.Time             `json:"created_at"`
	UpdatedAt         time.Time             `json:"updated_at"`
}

type NotificationFilter struct {
	Channel *shared.Channel
	Status  *NotificationStatus
}
