package usecase

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/multierr"
	"payment_processing_system/internal/domain"
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

func (a *ApplierUseCase) ApplyTransaction(ctx context.Context, transaction entity.Transaction) error {
	switch transaction.Type {
	case entity.TypeTransfer:
		return a.applyTransfer(ctx, transaction)
	case entity.TypePayment:
		return a.applyPayForService(ctx, transaction)
	case entity.TypeOuterIncreasing:
		return a.applyIncrease(ctx, transaction)
	case entity.TypeOuterDecreasing:
		return a.applyDecrease(ctx, transaction)
	}
	return fmt.Errorf("type = %s ; %w", transaction.Type, domain.UnknownTransactionTypeErr)
}

func (a *ApplierUseCase) applyTransfer(ctx context.Context, transaction entity.Transaction) error {
	// Increase destination balance
	err := a.bs.ChangeAmount(ctx, *transaction.DestinationID, transaction.Amount)
	if err != nil {
		// If destination balance was increased with error then destination balance should be decreased
		if !errors.Is(err, domain.BalanceWasNotIncreased) {
			multierr.AppendInto(&err, a.bs.ChangeAmount(ctx, *transaction.DestinationID, -transaction.Amount))
		}
		// Cancel transaction by producer
		multierr.AppendInto(&err, a.producer.CancelTransaction(transaction.ID))
		return err
	}
	// Decrease source balance
	err = a.bs.ChangeAmount(ctx, *transaction.SourceID, -transaction.Amount)
	if err != nil {
		// If source balance was decreased with error then source balance should be increased
		if !errors.Is(err, domain.BalanceWasNotDecreased) {
			multierr.AppendInto(&err, a.bs.ChangeAmount(ctx, *transaction.SourceID, transaction.Amount))
		}
		// Decrease destination balance
		multierr.AppendInto(&err, a.bs.ChangeAmount(ctx, *transaction.DestinationID, -transaction.Amount))
		// Cancel transaction by producer
		multierr.AppendInto(&err, a.producer.CancelTransaction(transaction.ID))
		return err
	}
	err = a.producer.CompleteTransaction(transaction.ID)
	return err
}

func (a *ApplierUseCase) applyIncrease(ctx context.Context, transaction entity.Transaction) error {
	err := a.bs.ChangeAmount(ctx, *transaction.DestinationID, transaction.Amount)
	if err != nil {
		multierr.AppendInto(&err, a.producer.CancelTransaction(transaction.ID))
		return err
	}
	err = a.producer.CompleteTransaction(transaction.ID)
	return err
}

func (a *ApplierUseCase) applyDecrease(ctx context.Context, transaction entity.Transaction) error {
	err := a.bs.ChangeAmount(ctx, *transaction.SourceID, -transaction.Amount)
	if err != nil {
		multierr.AppendInto(&err, a.producer.CancelTransaction(transaction.ID))
		return err
	}
	err = a.producer.CompleteTransaction(transaction.ID)
	return err
}

func (a *ApplierUseCase) applyPayForService(ctx context.Context, transaction entity.Transaction) error {
	err := a.bs.ChangeAmount(ctx, *transaction.SourceID, -transaction.Amount)
	if err != nil {
		multierr.AppendInto(&err, a.producer.CancelTransaction(transaction.ID))
		return err
	}
	err = a.producer.CompleteTransaction(transaction.ID)
	return err
}
