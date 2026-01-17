package kafka

import (
	"context"
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/ckshitij/notify-srv/internal/logger"
)

type NotificationService interface {
	Process(ctx context.Context, notificationID int64) error
}

type Consumer struct {
	group   sarama.ConsumerGroup
	topic   string
	service NotificationService
	log     logger.Logger
	workers int

	// internal
	workerPool chan struct{}
}

func NewConsumer(
	brokers []string,
	groupID string,
	topic string,
	service NotificationService,
	log logger.Logger,
	workers int,
) (*Consumer, error) {

	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0

	// Consumer Group configs
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRange()
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Return.Errors = true

	group, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		group:   group,
		topic:   topic,
		service: service,
		log:     log,
		workers: workers,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) {
	for {
		if err := c.group.Consume(ctx, []string{c.topic}, c); err != nil {
			c.log.Error(ctx, "consumer group consume error", logger.Error(err))
		}

		if ctx.Err() != nil {
			return
		}
	}
}

func (c *Consumer) Setup(session sarama.ConsumerGroupSession) error {
	c.log.Info(
		session.Context(),
		"kafka consumer group setup",
		logger.String("topic", c.topic),
	)

	// Initialize worker pool per rebalance
	c.workerPool = make(chan struct{}, c.workers)

	return nil
}

func (c *Consumer) Cleanup(session sarama.ConsumerGroupSession) error {
	c.log.Info(
		session.Context(),
		"kafka consumer group cleanup: waiting for in-flight messages",
	)

	// Drain worker pool → wait for all workers to finish
	for i := 0; i < c.workers; i++ {
		c.workerPool <- struct{}{}
	}

	close(c.workerPool)

	c.log.Info(
		session.Context(),
		"kafka consumer group cleanup complete",
	)

	return nil
}

func (c *Consumer) ConsumeClaim(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {

	for msg := range claim.Messages() {
		c.workerPool <- struct{}{} // acquire slot

		go func(m *sarama.ConsumerMessage) {
			defer func() { <-c.workerPool }()

			var notificationID int64
			if err := json.Unmarshal(m.Value, &notificationID); err != nil {
				c.log.Error(session.Context(), "failed to unmarshal message", logger.Error(err))
				return
			}

			if err := c.service.Process(session.Context(), notificationID); err != nil {
				c.log.Error(
					session.Context(),
					"failed to process notification",
					logger.Error(err),
					logger.Int64("notification_id", notificationID),
				)
				return
			}

			session.MarkMessage(m, "")
		}(msg)
	}

	return nil
}
