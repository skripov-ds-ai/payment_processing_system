package kafka

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"github.com/segmentio/kafka-go"
	"payment_processing_system/internal/domain/entity"
	"payment_processing_system/pkg/kafka/pubsub"
	"time"
)

type ApplyTransactionProducer struct {
	topic    string
	producer pubsub.Producer
}

func NewApplyTransactionProducer(topic string, producer pubsub.Producer) *ApplyTransactionProducer {
	return &ApplyTransactionProducer{topic: topic, producer: producer}
}

// TODO: implement
func (a *ApplyTransactionProducer) ApplyTransaction(ctx context.Context, transaction entity.Transaction) error {
	dto := ApplyTransactionEvent{Transaction: transaction}
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
