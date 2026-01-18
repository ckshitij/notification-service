package kafka

import (
	"time"

	"github.com/IBM/sarama"
)

type Producer struct {
	producer sarama.SyncProducer
}

func NewProducer(brokers []string) (*Producer, error) {
	config := sarama.NewConfig()

	// Kafka version
	config.Version = sarama.V2_8_0_0

	// Required for SyncProducer
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	// Durability guarantees
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Retry.Backoff = 200 * time.Millisecond

	// Ordering & idempotency
	config.Producer.Idempotent = true
	config.Net.MaxOpenRequests = 1

	// Performance tuning (safe defaults)
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.Frequency = 5 * time.Millisecond

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &Producer{producer: producer}, nil
}

func (p *Producer) SendMessage(topic string, key string, message []byte) (int32, int64, error) {

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(message),
	}

	return p.producer.SendMessage(msg)
}

func (p *Producer) Close() error {
	return p.producer.Close()
}
