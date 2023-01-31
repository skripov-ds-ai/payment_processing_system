package kafka

import "payment_processing_system/internal/domain/entity"

type ApplyTransactionProducer struct {
}

func NewApplyTransactionProducer() *ApplyTransactionProducer {
	return &ApplyTransactionProducer{}
}

// TODO: implement
func (a *ApplyTransactionProducer) ApplyTransaction(transaction entity.Transaction) error {
	return nil
}
