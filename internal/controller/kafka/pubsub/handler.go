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
			// TODO: add constants
			ctx := session.Context()
			switch message.Topic {
			case "apply":
				var event messages.ApplyTransactionEvent
				err := jsoniter.Unmarshal(message.Value, &event)
				if err != nil {
					return err
				}
				err = c.applier.ApplyTransaction(ctx, event.Transaction)
				if err != nil {
					return err
				}
			case "complete":
				var event messages.CompleteTransactionEvent
				err := jsoniter.Unmarshal(message.Value, &event)
				if err != nil {
					return err
				}
				err = c.verifier.CompleteTransactionByID(ctx, event.TransactionID)
				if err != nil {
					return err
				}
			case "cancel":
				var event messages.CancelTransactionEvent
				err := jsoniter.Unmarshal(message.Value, &event)
				if err != nil {
					return err
				}
				err = c.verifier.CancelTransactionByID(ctx, event.TransactionID)
				if err != nil {
					return err
				}
			}
			//var (
			//	sender    = string(message.Value)
			//	offset    = message.Offset
			//	partition = message.Partition
			//)

			//log.Printf("offset: %d, partition: %d, value: %v", offset, partition, sender)
			//metric.IncMessagesConsumed(sender, offset, partition)
			// TODO
			session.MarkMessage(message, "")
		}
	}
}
