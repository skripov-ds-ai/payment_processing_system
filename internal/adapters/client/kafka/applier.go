package kafka

import (
	"context"
	"payment_processing_system/internal/domain/entity"
	"payment_processing_system/pkg/kafka/messages"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/segmentio/kafka-go"
)

type Producer interface {
	PublishMessage(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}

type ApplyTransactionProducer struct {
	topic    string
	producer Producer
}

func NewApplyTransactionProducer(topic string, producer Producer) *ApplyTransactionProducer {
	return &ApplyTransactionProducer{topic: topic, producer: producer}
}

func (a *ApplyTransactionProducer) ApplyTransaction(ctx context.Context, transaction entity.Transaction) error {
	dto := messages.ApplyTransactionEvent{Transaction: transaction}
	dtoBs, err := jsoniter.Marshal(dto)
	if err != nil {
		return err
	}
	message := kafka.Message{
		Topic: a.topic,
		Value: dtoBs,
		Time:  time.Now(),
	}
	return a.producer.PublishMessage(ctx, message)
}
