package pubsub

import (
	"context"
	"github.com/Shopify/sarama"
	jsoniter "github.com/json-iterator/go"
	"payment_processing_system/internal/controller/kafka/messages"
	"payment_processing_system/internal/domain/entity"
)

type ApplierUseCase interface {
	ApplyTransaction(ctx context.Context, transaction entity.Transaction) error
}

type VerifierUseCase interface {
	CancelTransactionByID(ctx context.Context, id uint64) error
	CompleteTransactionByID(ctx context.Context, id uint64) error
}

type ConsumerHandler struct {
	applier  ApplierUseCase
	verifier VerifierUseCase
}

func NewConsumerHandler(applier ApplierUseCase, verifier VerifierUseCase) *ConsumerHandler {
	return &ConsumerHandler{applier: applier, verifier: verifier}
}

func (c *ConsumerHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *ConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case <-session.Context().Done():
			return nil
		case message, ok := <-claim.Messages():
			if !ok {
				return nil
			}
			err := c.processMessage(session.Context(), message)
			if err != nil {
				// TODO: log err
				continue
			}
			// TODO
			session.MarkMessage(message, "")
		}
	}
}

func (c *ConsumerHandler) processMessage(ctx context.Context, message *sarama.ConsumerMessage) (err error) {
	switch message.Topic {
	case "apply":
		var event messages.ApplyTransactionEvent
		err = jsoniter.Unmarshal(message.Value, &event)
		if err != nil {
			return err
		}
		err = c.applier.ApplyTransaction(ctx, event.Transaction)
	case "complete":
		var event messages.CompleteTransactionEvent
		err = jsoniter.Unmarshal(message.Value, &event)
		if err != nil {
			return err
		}
		err = c.verifier.CompleteTransactionByID(ctx, event.TransactionID)
	case "cancel":
		var event messages.CancelTransactionEvent
		err = jsoniter.Unmarshal(message.Value, &event)
		if err != nil {
			return err
		}
		err = c.verifier.CancelTransactionByID(ctx, event.TransactionID)
	}
	if err != nil {
		return err
	}
	return nil
}
