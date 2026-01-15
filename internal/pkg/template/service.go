package template

import (
	"context"
)

type TemplateService interface {
	Create(ctx context.Context, tpl Template) (int64, error)
	GetByID(ctx context.Context, templateID int64) (*Template, error)
	List(ctx context.Context, filter TemplateFilter) ([]*Template, error)
	CacheReloadSystemTemplates(ctx context.Context) error
	InvalidateTemplateCache(ctx context.Context, templateID int64) error
	Render(ctx context.Context, templateID int64, data map[string]any) (*Template, error)
}
