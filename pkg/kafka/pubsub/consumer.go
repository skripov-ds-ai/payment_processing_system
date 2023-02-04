package pubsub

import (
	"context"
	"payment_processing_system/pkg/logger"
	"time"

	"github.com/segmentio/kafka-go"
)

type consumer struct {
	r       *kafka.Reader
	log     *logger.Logger
	address []string
	topic   string
}

func NewConsumer(address []string, topic, groupID string, log *logger.Logger) *consumer {
	// TODO: check NewConsumerGroup - https://github.com/segmentio/kafka-go/blob/294fbdbf43ce5c2bc6aa89d8a35db949e1f95367/example_consumergroup_test.go#L11
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        address,
		Topic:          topic,
		GroupID:        groupID,
		CommitInterval: time.Second,
		Logger:         kafka.LoggerFunc(log.Infof),
		ErrorLogger:    kafka.LoggerFunc(log.Errorf),
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
	})
	return &consumer{r: r, log: log, address: address, topic: topic}
}

func (c *consumer) ReadMessage(ctx context.Context) (kafka.Message, error) {
	return c.r.ReadMessage(ctx)
}

func (c *consumer) Close() error {
	return c.r.Close()
}
