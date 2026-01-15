package template

import (
	"time"

	"github.com/ckshitij/notify-srv/internal/shared"
)

type Template struct {
	ID            int64               `json:"id"`
	Name          string              `json:"name"`
	Description   string              `json:"description"`
	Channel       shared.Channel      `json:"channel"`
	Type          shared.TemplateType `json:"type"`
	ActiveVersion int                 `json:"active_version"`
	CreatedBy     int64               `json:"created_by"`
	UpdatedBy     int64               `json:"updated_by"`
	CreatedAt     time.Time           `json:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at"`
}

type TemplateVersion struct {
	ID         int64     `json:"id"`
	TemplateID int64     `json:"template_id"`
	Version    int       `json:"version"`
	Subject    string    `json:"subject"`
	Body       string    `json:"body"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
}

type ListTemplatesFilter struct {
	Channel *shared.Channel
	Type    *shared.TemplateType
}

type TemplateWithActiveVersion struct {
	ID            int64               `json:"id"`
	Name          string              `json:"name"`
	Description   string              `json:"description"`
	Channel       shared.Channel      `json:"channel"`
	Type          shared.TemplateType `json:"type"`
	ActiveVersion int                 `json:"active_version"`
	Subject       *string             `json:"subject,omitempty"`
	Body          string              `json:"body"`
	CreatedAt     time.Time           `json:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at"`
}
