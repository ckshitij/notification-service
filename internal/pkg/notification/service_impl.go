package notification

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ckshitij/notify-srv/internal/config"
	"github.com/ckshitij/notify-srv/internal/kafka"
	"github.com/ckshitij/notify-srv/internal/logger"
	"github.com/ckshitij/notify-srv/internal/pkg/renderer"
	"github.com/ckshitij/notify-srv/internal/pkg/template"
	"github.com/ckshitij/notify-srv/internal/shared"
)

type serviceImpl struct {
	repo         Repository
	renderer     renderer.Renderer
	senders      map[shared.Channel]Sender
	templateRepo template.TemplateRepository
	log          logger.Logger
	producer     *kafka.Producer
	cfg          *config.Config
}

func NewNotificationService(
	repo Repository,
	renderer renderer.Renderer,
	senders map[shared.Channel]Sender,
	templateRepo template.TemplateRepository,
	log logger.Logger,
	producer *kafka.Producer,
	cfg *config.Config,
) Service {
	return &serviceImpl{repo, renderer, senders, templateRepo, log, producer, cfg}
}

func (s *serviceImpl) SendNow(ctx context.Context, n *Notification) (int64, error) {

	n.Status = StatusPending

	id, err := s.repo.Create(ctx, n)
	if err != nil {
		return -1, err
	}

	msg, err := json.Marshal(id)
	if err != nil {
		s.log.Error(ctx, "failed to marshal notification id", logger.Error(err))
		return -1, err
	}

	topic, ok := s.cfg.Kafka.Topics[string(n.Channel)]
	if !ok {
		return -1, fmt.Errorf("kafka topic not found for channel %s", n.Channel)
	}

	_, _, err = s.producer.SendMessage(topic, msg)
	if err != nil {
		s.log.Error(ctx, "failed to send message to kafka", logger.Error(err))
		return -1, err
	}

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
		s.log.Warn(ctx, "failed to aquire the record", logger.String("status", "sending"), logger.Int64("notificationID", notificationID))
		// return nil
	}

	n, err := s.repo.GetByID(ctx, notificationID)
	if err != nil {
		return err
	}

	// Load template version
	tplVersion, err := s.templateRepo.GetByID(ctx, n.TemplateID)
	if err != nil || tplVersion == nil {
		s.log.Warn(ctx, "failed to get template info",
			logger.Int64("templateID", n.TemplateID),
			logger.Int64("notificationID", notificationID),
			logger.Error(err),
		)
		s.repo.UpdateStatus(ctx, n.ID, StatusFailed)
		return err
	}

	s.log.Info(ctx, "received data", logger.Field{
		Key:   "data",
		Value: tplVersion,
	})

	// Render content
	content, err := s.renderer.Render(tplVersion.Subject, tplVersion.Body, n.TemplateKeyValue)
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
