package template

import (
	"context"

	"github.com/ckshitij/notify-srv/internal/pkg/renderer"
	"github.com/ckshitij/notify-srv/internal/shared"
)

type ServiceImpl struct {
	repo     TemplateRepository
	renderer renderer.Renderer
}

func NewTemplateService(repo TemplateRepository, renderer renderer.Renderer) TemplateService {
	return &ServiceImpl{
		repo:     repo,
		renderer: renderer,
	}
}

func (s *ServiceImpl) Create(ctx context.Context, tpl Template) (int64, error) {
	if tpl.Type == shared.SystemTemplate {
		return -1, shared.ErrSystemTemplateNotPermitted
	}
	return s.repo.Create(ctx, tpl)
}

func (s *ServiceImpl) GetByID(ctx context.Context, templateID int64) (*Template, error) {
	return s.repo.GetByID(ctx, templateID)
}

func (s *ServiceImpl) CacheReloadSystemTemplates(ctx context.Context) error {
	return s.repo.CacheReloadSystemTemplates(ctx)
}

func (s *ServiceImpl) InvalidateTemplateCache(ctx context.Context, templateID int64) error {
	return s.repo.InvalidateTemplateCache(ctx, templateID)
}

func (s *ServiceImpl) Render(ctx context.Context, templateID int64, data map[string]any) (*Template, error) {

	tpl, err := s.repo.GetByID(ctx, templateID)
	if err != nil {
		return nil, err
	}

	rendered, err := s.renderer.Render(tpl.Subject, tpl.Body, data)
	if err != nil {
		return nil, err
	}

	tpl.Subject = rendered.Subject
	tpl.Body = rendered.Body
	return tpl, nil
}

func (s *ServiceImpl) List(ctx context.Context, filter TemplateFilter) ([]*Template, error) {
	return s.repo.List(ctx, filter)
}
