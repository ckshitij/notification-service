package notification

import (
	"context"
	"errors"
	"time"

	"github.com/ckshitij/notify-srv/internal/logger"
	"github.com/ckshitij/notify-srv/internal/renderer"
	"github.com/ckshitij/notify-srv/internal/shared"
	"github.com/ckshitij/notify-srv/internal/template"
)

type Service interface {
	SendNow(ctx context.Context, n *Notification) (int64, error)
	Schedule(ctx context.Context, n *Notification, when time.Time) (int64, error)
	Process(ctx context.Context, notificationID int64) error
	GetByID(ctx context.Context, notificationID int64) (*Notification, error)
	List(ctx context.Context, filter NotificationFilter) ([]*Notification, error)
}

type serviceImpl struct {
	repo         Repository
	renderer     renderer.Renderer
	senders      map[shared.Channel]Sender
	templateRepo template.Repository
	log          logger.Logger
}

func New(
	repo Repository,
	renderer renderer.Renderer,
	senders map[shared.Channel]Sender,
	templateRepo template.Repository,
	log logger.Logger,
) Service {
	return &serviceImpl{repo, renderer, senders, templateRepo, log}
}

func (s *serviceImpl) SendNow(ctx context.Context, n *Notification) (int64, error) {

	n.Status = StatusPending

	id, err := s.repo.Create(ctx, n)
	if err != nil {
		return -1, err
	}

	go func() {
		ct := context.Background()
		err := s.Process(ct, id)
		if err != nil {
			s.log.Error(ct, "failed to process", logger.Error(err))
		}
	}()

	return id, nil
}

func (s *serviceImpl) Schedule(ctx context.Context, n *Notification, when time.Time) (int64, error) {

	n.Status = StatusScheduled
	n.ScheduledAt = &when

	return s.repo.Create(ctx, n)
}

func (s *serviceImpl) Process(
	ctx context.Context,
	notificationID int64,
) error {

	acquired, err := s.repo.AcquireForSending(ctx, notificationID)
	if err != nil {
		return err
	}
	if !acquired {
		// Someone else is processing or it's already done
		return nil
	}

	n, err := s.repo.GetByID(ctx, notificationID)
	if err != nil {
		return err
	}

	// Load template version
	tplVersion, err := s.templateRepo.GetActiveVersion(
		ctx,
		n.TemplateVersionID,
	)
	if err != nil {
		s.repo.UpdateStatus(ctx, n.ID, StatusFailed)
		return err
	}

	s.log.Info(ctx, "received data", logger.Field{
		Key:   "data",
		Value: tplVersion,
	})

	// Render content
	content, err := s.renderer.Render(tplVersion.Body, tplVersion.Subject, n.TemplateKeyValue)
	if err != nil {
		s.repo.UpdateStatus(ctx, n.ID, StatusFailed)
		return err
	}

	// Resolve sender
	sender, ok := s.senders[n.Channel]
	if !ok {
		s.repo.UpdateStatus(ctx, n.ID, StatusFailed)
		return errors.New("sender not configured")
	}

	// Send
	if err := sender.Send(ctx, *n, content); err != nil {
		s.repo.UpdateStatus(ctx, n.ID, StatusFailed)
		return err
	}

	// Mark sent
	now := time.Now()
	if err := s.repo.MarkSent(ctx, n.ID, now); err != nil {
		return err
	}

	return nil
}

func (s *serviceImpl) GetByID(ctx context.Context, notificationID int64) (*Notification, error) {
	return s.repo.GetByID(ctx, notificationID)
}

func (s *serviceImpl) List(ctx context.Context, filter NotificationFilter) ([]*Notification, error) {
	return s.repo.List(ctx, filter)
}
