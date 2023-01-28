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
	CreateDefaultTransaction(ctx context.Context, sourceID, destinationID *string, amount float32, ttype entity.TransactionType) (*entity.Transaction, error)
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

func (buc *BalanceUseCase) GetByID(ctx context.Context, id string) (*entity.Balance, error) {
	return buc.bs.GetByID(ctx, id)
}

func (buc *BalanceUseCase) Transfer(ctx context.Context, idFrom, idTo *string, amount float32) (err error) {
	if idFrom == nil {
		return domain.TransactionNilSourceErr
	}
	if idTo == nil {
		return domain.TransactionNilDestinationErr
	}
	if utils.IsZero(amount) {
		return fmt.Errorf("idFrom = %q ; idFrom = %q ; amount = %f ; %w", *idFrom, *idTo, amount, domain.ChangeBalanceByZeroAmountErr)
	}
	if amount < 0 {
		return fmt.Errorf("idFrom = %q ; idFrom = %q ; amount = %f ; %w", *idFrom, *idTo, amount, domain.NegativeAmountTransactionErr)
	}
	if *idFrom == *idTo {
		return fmt.Errorf("idFrom = %q ; idFrom = %q ; %w", *idFrom, *idTo, domain.TransactionSourceDestinationAreEqualErr)
	}
	// Check existence of idFrom balance
	_, err = buc.bs.GetByID(ctx, *idFrom)
	if err != nil {
		return err
	}
	var transaction *entity.Transaction
	defer func() {
		// Cancel transaction on err
		if err != nil && transaction != nil {
			multierr.AppendInto(&err, buc.ts.CancelByID(ctx, transaction.ID))
		}
	}()
	// Create transaction
	transaction, err = buc.ts.CreateDefaultTransaction(ctx, idFrom, idTo, amount, entity.TypeTransfer)
	// Cancel transaction on err
	if err != nil {
		return err
	}
	// Change transaction status to "processing"
	err = buc.ts.ProcessingByID(ctx, transaction.ID)
	// Cancel transaction on err
	if err != nil {
		return err
	}
	// Change balance by amount
	err = buc.bs.ChangeAmount(ctx, *idTo, amount)
	if err != nil {
		// Change balance by -amount on err
		if errors.Is(err, domain.BalanceWasNotIncreased) {
			multierr.AppendInto(&err, buc.bs.ChangeAmount(ctx, *idTo, -amount))
		}
		return err
	}
	err = buc.bs.ChangeAmount(ctx, *idFrom, -amount)
	if err != nil {
		// Change balance by amount on err
		if errors.Is(err, domain.BalanceWasNotDecreased) {
			multierr.AppendInto(&err, buc.bs.ChangeAmount(ctx, *idFrom, amount))
		}
		// Change balance by -amount on err
		if errors.Is(err, domain.BalanceWasNotIncreased) {
			multierr.AppendInto(&err, buc.bs.ChangeAmount(ctx, *idTo, -amount))
		}
		return err
	}
	// Change transaction status to "completed"
	err = buc.ts.CompletedByID(ctx, transaction.ID)
	if err != nil {
		// Change balance by amount on err
		if errors.Is(err, domain.BalanceWasNotDecreased) {
			multierr.AppendInto(&err, buc.bs.ChangeAmount(ctx, *idFrom, amount))
		}
		// Change balance by -amount on err
		if errors.Is(err, domain.BalanceWasNotIncreased) {
			multierr.AppendInto(&err, buc.bs.ChangeAmount(ctx, *idTo, -amount))
		}
		return err
	}
	return nil
}

func (buc *BalanceUseCase) ChangeAmount(ctx context.Context, id *string, amount float32) (err error) {
	if id == nil {
		return domain.TransactionNilSourceOrDestinationErr
	}
	if utils.IsZero(amount) {
		return fmt.Errorf("idFrom = %q ; amount = %f ; %w", *id, amount, domain.ChangeBalanceByZeroAmountErr)
	}
	var transaction *entity.Transaction
	defer func() {
		// Cancel transaction on err
		if err != nil && transaction != nil {
			multierr.AppendInto(&err, buc.ts.CancelByID(ctx, transaction.ID))
		}
	}()
	// Create transaction
	if amount > 0 {
		transaction, err = buc.ts.CreateDefaultTransaction(ctx, nil, id, amount, entity.TypeOuterIncreasing)
	} else {
		transaction, err = buc.ts.CreateDefaultTransaction(ctx, id, nil, -amount, entity.TypeOuterDecreasing)
	}
	// Cancel transaction on err
	if err != nil {
		return err
	}
	// Change transaction status to "processing"
	err = buc.ts.ProcessingByID(ctx, transaction.ID)
	// Cancel transaction on err
	if err != nil {
		return err
	}
	// Change balance by amount
	err = buc.bs.ChangeAmount(ctx, *id, amount)
	if err != nil {
		// Change balance by -amount on err
		if errors.Is(err, domain.BalanceWasNotDecreased) || errors.Is(err, domain.BalanceWasNotIncreased) {
			multierr.AppendInto(&err, buc.bs.ChangeAmount(ctx, *id, -amount))
		}
		return err
	}
	// Change transaction status to "completed"
	err = buc.ts.CompletedByID(ctx, transaction.ID)
	if err != nil {
		// Change balance by -amount on err
		multierr.AppendInto(&err, buc.bs.ChangeAmount(ctx, *id, -amount))
		return err
	}
	return nil
}

func (buc *BalanceUseCase) PayForService(ctx context.Context, id *string, amount float32) (err error) {
	if id == nil {
		return domain.TransactionNilSourceErr
	}
	if utils.IsZero(amount) {
		return fmt.Errorf("idFrom = %q ; amount = %f ; %w", *id, amount, domain.ChangeBalanceByZeroAmountErr)
	}
	if amount < 0 {
		return fmt.Errorf("idFrom = %q ; amount = %f ; %w", *id, amount, domain.NegativeAmountTransactionErr)
	}
	// Check existence of idFrom balance
	_, err = buc.bs.GetByID(ctx, *id)
	if err != nil {
		return err
	}
	var transaction *entity.Transaction
	defer func() {
		// Cancel transaction on err
		if err != nil && transaction != nil {
			multierr.AppendInto(&err, buc.ts.CancelByID(ctx, transaction.ID))
		}
	}()
	// Create transaction
	transaction, err = buc.ts.CreateDefaultTransaction(ctx, id, nil, -amount, entity.TypePayment)
	// Cancel transaction on err
	if err != nil {
		return err
	}
	// Change transaction status to "processing"
	err = buc.ts.ProcessingByID(ctx, transaction.ID)
	// Cancel transaction on err
	if err != nil {
		return err
	}
	// Change balance by -amount
	err = buc.bs.ChangeAmount(ctx, *id, -amount)
	if err != nil {
		// Change balance by amount on err
		if errors.Is(err, domain.BalanceWasNotDecreased) {
			multierr.AppendInto(&err, buc.bs.ChangeAmount(ctx, *id, amount))
		}
		return err
	}
	// Change transaction status to "completed"
	err = buc.ts.CompletedByID(ctx, transaction.ID)
	if err != nil {
		// Change balance by amount on err
		multierr.AppendInto(&err, buc.bs.ChangeAmount(ctx, *id, amount))
		return err
	}
	return nil
}
