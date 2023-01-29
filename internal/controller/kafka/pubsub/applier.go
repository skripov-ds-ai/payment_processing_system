package pubsub

import "payment_processing_system/internal/domain/entity"

type ApplierUseCase interface {
	ApplyTransaction(transaction entity.Transaction) error
}

type ApplierHandler struct {
	applier ApplierUseCase
}
