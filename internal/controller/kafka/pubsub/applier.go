package pubsub

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"github.com/segmentio/kafka-go"
	"go.uber.org/multierr"
	"payment_processing_system/internal/domain/entity"
	"payment_processing_system/pkg/kafka/messages"
)

type Consumer interface {
	ReadMessage(ctx context.Context) (kafka.Message, error)
	Close() error
}

type ApplierUseCase interface {
	ApplyTransaction(transaction entity.Transaction) error
}

type ApplierHandler struct {
	applier  ApplierUseCase
	consumer Consumer
}

// TODO: add graceful shutdown
func (a *ApplierHandler) Apply(ctx context.Context) error {
	// ctx := context.Background()
	var err0 error
	for {
		m, err := a.consumer.ReadMessage(ctx)
		if err != nil {
			err0 = err
			break
		}
		var event messages.ApplyTransactionEvent
		err = jsoniter.Unmarshal(m.Value, &event)
		if err != nil {
			err0 = err
			break
		}
		err = a.applier.ApplyTransaction(event.Transaction)
		if err != nil {
			err0 = err
			break
		}
	}
	if err1 := a.consumer.Close(); err1 != nil {
		return multierr.Append(err0, err1)
	}
	if err0 != nil {
		return err0
	}
	return nil
}
