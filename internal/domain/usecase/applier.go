package usecase

import (
	"context"
	"errors"
	"fmt"
	"payment_processing_system/internal/domain"
	"payment_processing_system/internal/domain/entity"

	"github.com/shopspring/decimal"
	"go.uber.org/multierr"
)

type ConfirmTransactionProducer interface {
	CancelTransaction(ctx context.Context, id uint64) error
	CompleteTransaction(ctx context.Context, id uint64) error
}

type BalanceGetChangeService interface {
	GetByID(ctx context.Context, id int64) (*entity.Balance, error)
	ChangeAmount(ctx context.Context, id int64, amount decimal.Decimal) error
}

type ApplierUseCase struct {
	bs       BalanceGetChangeService
	producer ConfirmTransactionProducer
}

func NewApplierUseCase(bs BalanceGetChangeService, producer ConfirmTransactionProducer) *ApplierUseCase {
	return &ApplierUseCase{bs: bs, producer: producer}
}

func (a *ApplierUseCase) ApplyTransaction(ctx context.Context, transaction entity.Transaction) error {
	switch transaction.TType {
	case entity.TypeTransfer:
		return a.applyTransfer(ctx, transaction)
	case entity.TypePayment:
		return a.applyPayForService(ctx, transaction)
	case entity.TypeOuterIncreasing:
		return a.applyIncrease(ctx, transaction)
	case entity.TypeOuterDecreasing:
		return a.applyDecrease(ctx, transaction)
	}
	return fmt.Errorf("type = %s ; %w", transaction.TType, domain.UnknownTransactionTypeErr)
}

func (a *ApplierUseCase) applyTransfer(ctx context.Context, transaction entity.Transaction) error {
	// Increase destination balance
	err := a.bs.ChangeAmount(ctx, *transaction.DestinationID, transaction.Amount)
	if err != nil {
		// If destination balance was increased with error then destination balance should be decreased
		if !errors.Is(err, domain.BalanceWasNotIncreasedErr) {
			multierr.AppendInto(&err, a.bs.ChangeAmount(ctx, *transaction.DestinationID, transaction.Amount.Neg()))
		}
		// Cancel transaction by producer
		multierr.AppendInto(&err, a.producer.CancelTransaction(ctx, transaction.ID))
		return err
	}
	// Decrease source balance
	err = a.bs.ChangeAmount(ctx, *transaction.SourceID, transaction.Amount.Neg())
	if err != nil {
		// If source balance was decreased with error then source balance should be increased
		if !errors.Is(err, domain.BalanceWasNotDecreasedErr) {
			multierr.AppendInto(&err, a.bs.ChangeAmount(ctx, *transaction.SourceID, transaction.Amount))
		}
		// Decrease destination balance
		multierr.AppendInto(&err, a.bs.ChangeAmount(ctx, *transaction.DestinationID, transaction.Amount.Neg()))
		// Cancel transaction by producer
		multierr.AppendInto(&err, a.producer.CancelTransaction(ctx, transaction.ID))
		return err
	}
	err = a.producer.CompleteTransaction(ctx, transaction.ID)
	return err
}

func (a *ApplierUseCase) applyIncrease(ctx context.Context, transaction entity.Transaction) error {
	err := a.bs.ChangeAmount(ctx, *transaction.DestinationID, transaction.Amount)
	if err != nil {
		if !errors.Is(err, domain.BalanceWasNotIncreasedErr) {
			multierr.AppendInto(&err, a.bs.ChangeAmount(ctx, *transaction.DestinationID, transaction.Amount.Neg()))
		}
		multierr.AppendInto(&err, a.producer.CancelTransaction(ctx, transaction.ID))
		return err
	}
	err = a.producer.CompleteTransaction(ctx, transaction.ID)
	return err
}

func (a *ApplierUseCase) applyDecrease(ctx context.Context, transaction entity.Transaction) error {
	err := a.bs.ChangeAmount(ctx, *transaction.SourceID, transaction.Amount.Neg())
	if err != nil {
		if !errors.Is(err, domain.BalanceWasNotDecreasedErr) {
			multierr.AppendInto(&err, a.bs.ChangeAmount(ctx, *transaction.DestinationID, transaction.Amount))
		}
		multierr.AppendInto(&err, a.producer.CancelTransaction(ctx, transaction.ID))
		return err
	}
	err = a.producer.CompleteTransaction(ctx, transaction.ID)
	return err
}

func (a *ApplierUseCase) applyPayForService(ctx context.Context, transaction entity.Transaction) error {
	err := a.bs.ChangeAmount(ctx, *transaction.SourceID, transaction.Amount.Neg())
	if err != nil {
		if !errors.Is(err, domain.BalanceWasNotDecreasedErr) {
			multierr.AppendInto(&err, a.bs.ChangeAmount(ctx, *transaction.DestinationID, transaction.Amount))
		}
		multierr.AppendInto(&err, a.producer.CancelTransaction(ctx, transaction.ID))
		return err
	}
	err = a.producer.CompleteTransaction(ctx, transaction.ID)
	return err
}
