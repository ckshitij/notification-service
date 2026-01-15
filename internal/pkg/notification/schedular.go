package notification

import (
	"context"
	"time"

	"github.com/ckshitij/notify-srv/internal/logger"
)

type Scheduler struct {
	service Service
	repo    Repository
	log     logger.Logger

	interval time.Duration
	batch    int
}

func NewSchedular(service Service, repo Repository, log logger.Logger, interval time.Duration, batch int) *Scheduler {
	return &Scheduler{
		service:  service,
		repo:     repo,
		log:      log,
		interval: interval,
		batch:    batch,
	}
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

func (s *Scheduler) tick(ctx context.Context) {
	scheduled, err := s.repo.FindDue(ctx, s.batch)
	if err != nil {
		s.log.Error(ctx, "failed to fetch scheduled notifications",
			logger.Error(err),
		)
		return
	}

	// Stuck sending notifications
	stuck, err := s.repo.FindStuckSending(ctx, 10*time.Minute, s.batch)
	if err != nil {
		s.log.Error(ctx, "failed to fetch stucked notifications",
			logger.Error(err),
		)
		return
	}

	ids := append(scheduled, stuck...)

	for _, id := range ids {
		go func(notificationID int64) {
			if err := s.service.Process(ctx, notificationID); err != nil {
				s.log.Error(ctx, "failed to process notification",
					logger.Int64("notification_id", notificationID),
					logger.Error(err),
				)
			}
		}(id)
	}
}
