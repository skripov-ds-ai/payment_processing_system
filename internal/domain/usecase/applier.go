package usecase

import (
	"context"
	"payment_processing_system/internal/domain/entity"
)

type ConfirmTransactionProducer interface {
	CancelTransaction(id string) error
	CompleteTransaction(id string) error
}

type BalanceService interface {
	GetByID(ctx context.Context, id string) (*entity.Balance, error)
	ChangeAmount(ctx context.Context, id string, amount float32) error
}

type ApplierUseCase struct {
	bs       BalanceService
	producer ConfirmTransactionProducer
}

func NewApplierUseCase(bs BalanceService, producer ConfirmTransactionProducer) *ApplierUseCase {
	return &ApplierUseCase{bs: bs, producer: producer}
}

// TODO
func (a *ApplierUseCase) ApplyTransaction(transaction entity.Transaction) error {
	return nil
}

func (a *ApplierUseCase) applyTransfer(transaction entity.Transaction) error {
	return nil
}

func (a *ApplierUseCase) applyIncrease(transaction entity.Transaction) error {
	return nil
}

func (a *ApplierUseCase) applyDecrease(transaction entity.Transaction) error {
	return nil
}

func (a *ApplierUseCase) applyPayForService(transaction entity.Transaction) error {
	return nil
}
