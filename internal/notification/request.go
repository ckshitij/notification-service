package notification

import (
	"time"

	"github.com/ckshitij/notify-srv/internal/shared"
)

type SendNowRequest struct {
	Channel shared.Channel `json:"channel"`

	TemplateVersionID int64 `json:"template_version_id"`

	Recipient map[string]string `json:"recipient"`
	Payload   map[string]any    `json:"payload"`
}

type ScheduleRequest struct {
	SendNowRequest
	ScheduledAt time.Time `json:"scheduled_at"`
}
