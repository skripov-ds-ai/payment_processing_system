package pubsub

import (
	"context"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Producer interface {
	PublishMessage(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}

type producer struct {
	log     *zap.Logger
	address []string
	w       *kafka.Writer
}

// NewProducer create new kafka producer
func NewProducer(address []string, log *zap.Logger) *producer {
	// TODO: add Infof, Errorf to logger/logger interface
	writer := kafka.Writer{
		Addr: kafka.TCP(address...),
		//Logger: kafka.Logger(),
	}
	return &producer{log: log, address: address, w: &writer}
}

func (p *producer) PublishMessage(ctx context.Context, msgs ...kafka.Message) error {
	return p.w.WriteMessages(ctx, msgs...)
}

func (p *producer) Close() error {
	return p.w.Close()
}
