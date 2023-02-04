package pubsub

import (
	"context"
	"payment_processing_system/pkg/logger"

	"github.com/segmentio/kafka-go"
)

type producer struct {
	log     *logger.Logger
	address []string
	w       *kafka.Writer
	topic   string
}

// NewProducer create new kafka producer
func NewProducer(address []string, topic string, log *logger.Logger) *producer {
	writer := kafka.Writer{
		Addr:        kafka.TCP(address...),
		Topic:       topic,
		Logger:      kafka.LoggerFunc(log.Infof),
		ErrorLogger: kafka.LoggerFunc(log.Errorf),
	}
	return &producer{log: log, address: address, w: &writer, topic: topic}
}

func (p *producer) PublishMessage(ctx context.Context, msgs ...kafka.Message) error {
	return p.w.WriteMessages(ctx, msgs...)
}

func (p *producer) Close() error {
	return p.w.Close()
}
