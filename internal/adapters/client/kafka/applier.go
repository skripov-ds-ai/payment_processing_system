package kafka

import (
	"context"
	"github.com/Shopify/sarama"
	"payment_processing_system/internal/controller/kafka/messages"
	"payment_processing_system/internal/domain/entity"
	"time"

	jsoniter "github.com/json-iterator/go"
)

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
	message := sarama.ProducerMessage{Topic: a.topic, Value: sarama.ByteEncoder(dtoBs), Timestamp: time.Now()}
	return a.producer.PublishMessage(&message)
}
