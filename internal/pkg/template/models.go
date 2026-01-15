package template

import (
	"time"

	"github.com/ckshitij/notify-srv/internal/shared"
)

type Template struct {
	ID          int64               `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Channel     shared.Channel      `json:"channel"`
	Type        shared.TemplateType `json:"type"`
	IsActive    bool                `json:"is_active"`
	Subject     string              `json:"subject"`
	Body        string              `json:"body"`
	CreatedBy   int64               `json:"created_by,omitempty"`
	UpdatedBy   int64               `json:"updated_by,omitempty"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}

type CreateTemplateRequest struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Channel     shared.Channel `json:"channel"`
	Subject     string         `json:"subject"`
	Body        string         `json:"body"`
}

func (r CreateTemplateRequest) Validate() error {
	if r.Name == "" {
		return shared.ErrRequiredFieldName
	}
	if r.Channel == "" {
		return shared.ErrRequiredFieldChannel
	}
	if r.Body == "" {
		return shared.ErrRequiredFieldBody
	}
	if r.Channel == shared.ChannelEmail && r.Subject == "" {
		return shared.ErrRequiredFieldSubject
	}
	return nil
}

type RenderRequest struct {
	TemplateKeyValue map[string]any `json:"template_key_value"`
}

type TemplateFilter struct {
	Name     *string
	Channel  *shared.Channel
	Type     *shared.TemplateType
	IsActive *bool
	Limit    int
	Offset   int
}
