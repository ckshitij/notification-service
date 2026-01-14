package template

import (
	"time"

	"github.com/ckshitij/notify-srv/internal/shared"
)

type Template struct {
	ID          int64
	Name        string
	Description string

	Channel shared.Channel
	Type    shared.TemplateType

	ActiveVersion int

	CreatedBy int64
	UpdatedBy int64

	CreatedAt time.Time
	UpdatedAt time.Time
}

type TemplateVersion struct {
	ID         int64
	TemplateID int64

	Version int

	Subject string
	Body    string

	IsActive bool

	CreatedAt time.Time
}

type ListTemplatesFilter struct {
	Channel *shared.Channel
	Type    *shared.TemplateType
}

type TemplateWithActiveVersion struct {
	ID            int64
	Name          string
	Description   string
	Channel       shared.Channel
	Type          shared.TemplateType
	ActiveVersion int

	Subject *string
	Body    string

	CreatedAt time.Time
	UpdatedAt time.Time
}
