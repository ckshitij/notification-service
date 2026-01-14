package template

import (
	"context"
	"errors"

	"github.com/ckshitij/notify-srv/internal/shared"
)

type Service struct {
	repo     Repository
	renderer Renderer
}

func NewService(repo Repository, renderer Renderer) *Service {
	return &Service{
		repo:     repo,
		renderer: renderer,
	}
}

func (s *Service) CreateTemplate(ctx context.Context, tpl Template) error {
	if tpl.Type == shared.SystemTemplate {
		return errors.New("system templates cannot be created via API")
	}
	return s.repo.CreateTemplate(ctx, tpl)
}

func (s *Service) AddVersion(ctx context.Context, templateID int64, subject string, body string) error {
	active, err := s.repo.GetActiveVersion(ctx, templateID)
	if err != nil {
		return err
	}

	version := TemplateVersion{
		TemplateID: templateID,
		Version:    active.Version + 1,
		Subject:    subject,
		Body:       body,
		IsActive:   true,
	}

	return s.repo.CreateVersion(ctx, version)
}

func (s *Service) AddVersionByName(
	ctx context.Context,
	name string,
	channel shared.Channel,
	subject string,
	body string,
) error {

	if body == "" {
		return errors.New("template body cannot be empty")
	}

	tpl, err := s.repo.GetTemplate(ctx, name, shared.UserTemplate, channel)
	if err != nil {
		return err
	}
	if tpl == nil {
		return errors.New("template not found")
	}

	if tpl.Type == shared.SystemTemplate {
		return errors.New("system templates cannot be modified")
	}

	version := TemplateVersion{
		TemplateID: tpl.ID,
		Version:    tpl.ActiveVersion + 1,
		Subject:    subject,
		Body:       body,
	}

	return s.repo.CreateVersion(ctx, version)
}

func (s *Service) Render(ctx context.Context, templateName string, tplType shared.TemplateType, channel shared.Channel, data map[string]any) (*RenderedTemplate, error) {

	tpl, err := s.repo.GetTemplate(ctx, templateName, tplType, channel)
	if err != nil {
		return nil, err
	}

	version, err := s.repo.GetActiveVersion(ctx, tpl.ID)
	if err != nil {
		return nil, err
	}

	rendered, err := s.renderer.Render(*version, data)
	if err != nil {
		return nil, err
	}

	return &rendered, nil
}

func (s *Service) ListVersionsByName(
	ctx context.Context,
	name string,
	channel shared.Channel,
) ([]TemplateVersion, error) {

	tpl, err := s.repo.GetTemplate(ctx, name, shared.UserTemplate, channel)
	if err != nil {
		return nil, err
	}
	if tpl == nil {
		return nil, errors.New("template not found")
	}

	return s.repo.ListVersions(ctx, tpl.ID)
}

func (s *Service) ListTemplatesWithActiveVersion(
	ctx context.Context,
	channel *shared.Channel,
	tplType *shared.TemplateType,
) ([]TemplateWithActiveVersion, error) {

	filter := ListTemplatesFilter{
		Channel: channel,
		Type:    tplType,
	}

	return s.repo.ListTemplatesWithActiveVersion(ctx, filter)
}
