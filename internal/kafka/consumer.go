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
	consumer sarama.Consumer
	topic    string
	service  NotificationService
	log      logger.Logger
}

func NewConsumer(brokers []string, topic string, service NotificationService, log logger.Logger) (*Consumer, error) {
	config := sarama.NewConfig()
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		consumer: consumer,
		topic:    topic,
		service:  service,
		log:      log,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) {
	partitionConsumer, err := c.consumer.ConsumePartition(c.topic, 0, sarama.OffsetNewest)
	if err != nil {
		c.log.Error(ctx, "failed to consume partition", logger.Error(err))
		return
	}
	defer partitionConsumer.Close()

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			var notificationID int64
			if err := json.Unmarshal(msg.Value, &notificationID); err != nil {
				c.log.Error(ctx, "failed to unmarshal message", logger.Error(err))
				continue
			}
			err := c.service.Process(ctx, notificationID)
			if err != nil {
				c.log.Error(ctx, "failed to process notification", logger.Error(err), logger.Int64("notification_id", notificationID))
			}
		case <-ctx.Done():
			return
		}
	}
}
