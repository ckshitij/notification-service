package notification

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/ckshitij/notify-srv/internal/config"
	"github.com/ckshitij/notify-srv/internal/kafka"
	"github.com/ckshitij/notify-srv/internal/logger"
)

type Scheduler struct {
	repo     Repository
	log      logger.Logger
	interval time.Duration
	batch    int
	workers  int
	producer *kafka.Producer
	kafkaCfg *config.KafkaConfig
}

func NewSchedular(repo Repository, log logger.Logger, interval time.Duration, batch int, workers int, producer *kafka.Producer, kafkaCfg *config.KafkaConfig) *Scheduler {
	return &Scheduler{repo, log, interval, batch, workers, producer, kafkaCfg}
}

func (s *Scheduler) Run(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	s.log.Info(ctx, "scheduler started")

	for {
		select {
		case <-ctx.Done():
			s.log.Info(ctx, "scheduler stopped")
			return

		case <-ticker.C:
			s.tick(ctx)
		}
	}
}

func (s *Scheduler) fetchCandidates(ctx context.Context) ([]NotificationScheduled, error) {
	scheduled, err := s.repo.FindDue(ctx, s.batch)
	if err != nil {
		s.log.Error(ctx, "failed to fetch scheduled notifications", logger.Error(err))
		return nil, err
	}

	stuck, err := s.repo.FindStuckSending(ctx, 10*time.Minute, s.batch)
	if err != nil {
		s.log.Error(ctx, "failed to fetch stuck notifications", logger.Error(err))
		return nil, err
	}

	return append(scheduled, stuck...), nil
}

func (s *Scheduler) tick(ctx context.Context) {
	notificationIds, err := s.fetchCandidates(ctx)
	if err != nil {
		return
	}

	if len(notificationIds) == 0 {
		return
	}

	workerPool := make(chan struct{}, s.workers)

	for _, n := range notificationIds {
		select {
		case <-ctx.Done():
			return
		case workerPool <- struct{}{}:
			go func(n NotificationScheduled) {
				defer func() { <-workerPool }()

				topic, ok := s.kafkaCfg.Topics[string(n.Channel)]
				if !ok {
					s.log.Error(ctx, "missing kafka topic for channel",
						logger.String("channel", string(n.Channel)),
					)
					return
				}

				key := strconv.FormatInt(n.ID, 10)
				msg, _ := json.Marshal(n.ID)

				if _, _, err := s.producer.SendMessage(topic, key, msg); err != nil {
					s.log.Error(
						ctx,
						"failed to publish scheduled notification",
						logger.Int64("notification_id", n.ID),
						logger.Error(err),
					)
					return
				}

				// Mark as dispatched so scheduler won't pick it again
				if err := s.repo.UpdateStatus(ctx, n.ID, StatusDispatched); err != nil {
					s.log.Error(
						ctx,
						"failed to mark notification dispatched",
						logger.Int64("notification_id", n.ID),
						logger.Error(err),
					)
				}
			}(n)
		}
	}
}
