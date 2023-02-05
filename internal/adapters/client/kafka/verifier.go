package kafka

import (
	"context"
	"github.com/Shopify/sarama"
	jsoniter "github.com/json-iterator/go"
	"payment_processing_system/internal/controller/kafka/messages"
	"time"
)

type VerifyTransactionProducer struct {
	completeTopic string
	cancelTopic   string
	producer      Producer
}

func NewVerifyTransactionProducer(cancelTopic, confirmTopic string, producer Producer) *VerifyTransactionProducer {
	return &VerifyTransactionProducer{cancelTopic: cancelTopic, completeTopic: confirmTopic, producer: producer}
}

func (v *VerifyTransactionProducer) CancelTransaction(ctx context.Context, id uint64) error {
	dto := messages.CancelTransactionEvent{TransactionID: id}
	dtoBs, err := jsoniter.Marshal(dto)
	if err != nil {
		return err
	}
	message := sarama.ProducerMessage{Topic: v.cancelTopic, Value: sarama.ByteEncoder(dtoBs), Timestamp: time.Now()}
	return v.producer.PublishMessage(&message)
}

func (v *VerifyTransactionProducer) CompleteTransaction(ctx context.Context, id uint64) error {
	dto := messages.CancelTransactionEvent{TransactionID: id}
	dtoBs, err := jsoniter.Marshal(dto)
	if err != nil {
		return err
	}
	message := sarama.ProducerMessage{Topic: v.completeTopic, Value: sarama.ByteEncoder(dtoBs), Timestamp: time.Now()}
	return v.producer.PublishMessage(&message)
}
