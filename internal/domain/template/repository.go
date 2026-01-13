package template

import (
	"context"

	"github.com/ckshitij/notification-srv/internal/domain/shared"
)

type Repository interface {
	CreateTemplate(ctx context.Context, tpl Template) error
	GetTemplate(ctx context.Context, name string, tplType shared.TemplateType, channel shared.Channel) (*Template, error)

	CreateVersion(ctx context.Context, version TemplateVersion) error
	GetActiveVersion(ctx context.Context, templateID int64) (*TemplateVersion, error)
	GetVersion(ctx context.Context, templateID int64, version int) (*TemplateVersion, error)
	ListVersions(ctx context.Context, templateID int64) ([]TemplateVersion, error)
}
