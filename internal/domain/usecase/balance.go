package usecase

import (
	"context"
	"errors"
	"fmt"
	"payment_processing_system/internal/domain"
	"payment_processing_system/internal/domain/entity"
	"payment_processing_system/internal/utils"

	"go.uber.org/multierr"
)

type BalanceService interface {
	GetByID(ctx context.Context, id string) (*entity.Balance, error)
	ChangeAmount(ctx context.Context, id string, amount float32) error
}

type TransactionService interface {
	GetByID(ctx context.Context, id string) (*entity.Transaction, error)
	CreateDefaultTransaction(ctx context.Context, sourceID, destinationID *string, amount float32, ttype entity.TransactionType) (string, error)
	CancelByID(ctx context.Context, id string) error
	ProcessingByID(ctx context.Context, id string) error
	CompletedByID(ctx context.Context, id string) error
	ShouldRetryByID(ctx context.Context, id string) error
	CannotApplyByID(ctx context.Context, id string) error
}

type BalanceUseCase struct {
	bs BalanceService
	ts TransactionService
}

func NewBalanceUseCase(bs BalanceService, ts TransactionService) *BalanceUseCase {
	return &BalanceUseCase{bs: bs, ts: ts}
}

func (buc *BalanceUseCase) ChangeAmount(ctx context.Context, id string, amount float32) error {
	if utils.IsZero(amount) {
		return fmt.Errorf("id = %q ; amount = %f ; %w", id, amount, domain.ChangeBalanceByZeroAmountErr)
	}
	var transactionID string
	var err error
	// Create transaction
	if amount > 0 {
		transactionID, err = buc.ts.CreateDefaultTransaction(ctx, nil, &id, amount, entity.TypeOuterIncreasing)
	} else {
		transactionID, err = buc.ts.CreateDefaultTransaction(ctx, &id, nil, -amount, entity.TypeOuterDecreasing)
	}
	// Cancel transaction on err
	if err != nil {
		multierr.AppendInto(&err, buc.ts.CancelByID(ctx, transactionID))
		return err
	}
	// Change transaction status to "processing"
	err = buc.ts.ProcessingByID(ctx, transactionID)
	// Cancel transaction on err
	if err != nil {
		multierr.AppendInto(&err, buc.ts.CancelByID(ctx, transactionID))
		return err
	}
	// Change balance by amount
	err = buc.bs.ChangeAmount(ctx, id, amount)
	if err != nil {
		// Change balance by -amount on err
		if errors.Is(err, domain.BalanceWasNotDecreased) || errors.Is(err, domain.BalanceWasNotIncreased) {
			multierr.AppendInto(&err, buc.bs.ChangeAmount(ctx, id, -amount))
		}
		// Cancel transaction on err
		multierr.AppendInto(&err, buc.ts.CancelByID(ctx, id))
		return err
	}
	// Change transaction status to "completed"
	err = buc.ts.CompletedByID(ctx, id)
	if err != nil {
		// Change balance by -amount on err
		multierr.AppendInto(&err, buc.bs.ChangeAmount(ctx, id, -amount))
		// Cancel transaction on err
		multierr.AppendInto(&err, buc.ts.CancelByID(ctx, transactionID))
		return err
	}
	return nil
}
