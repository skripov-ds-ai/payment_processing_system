package usecase

import (
	"context"
	"fmt"
	"github.com/shopspring/decimal"
	"go.uber.org/multierr"
	"payment_processing_system/internal/domain"
	"payment_processing_system/internal/domain/entity"
)

// TODO: usecases будут 3 типов
// 1) Внешний(web) - получение balance, получение transactions, изменение balance(через внесение извне, списание, покупку, transfer);
// по сути получение(R) balance, transaction, создание(C) transaction
// 2) Внутренний(kafka) - получение событий, (R) balance, создание(C) balance, изменение(U) balance, отправка событий изменения transaction
// 3) Внутренний(kafka) - получение событий, (R) transaction, изменение(U) transaction

type ApplyTransactionProducer interface {
	ApplyTransaction(transaction entity.Transaction) error
}

type BalanceGetService interface {
	GetByID(ctx context.Context, id int64) (*entity.Balance, error)
}

type TransactionGetCreateService interface {
	GetByID(ctx context.Context, id uint64) (*entity.Transaction, error)
	CreateDefaultTransaction(ctx context.Context, sourceID, destinationID *int64, amount decimal.Decimal, ttype entity.TransactionType) (*entity.Transaction, error)
	CancelByID(ctx context.Context, id uint64) error
}

type ManagerUseCase struct {
	bs       BalanceGetService
	ts       TransactionGetCreateService
	producer ApplyTransactionProducer
}

func NewManagerUseCase(bs BalanceGetService, ts TransactionGetCreateService, producer ApplyTransactionProducer) *ManagerUseCase {
	return &ManagerUseCase{bs: bs, ts: ts, producer: producer}
}

// TODO: add GetBalanceTransactions(ctx context.Context, id string) ([]entity.Transaction, error)
func (buc *ManagerUseCase) GetBalanceTransactions(ctx context.Context, id string) ([]entity.Transaction, error) {
	return []entity.Transaction{}, nil
}

func (buc *ManagerUseCase) GetBalanceByID(ctx context.Context, id int64) (*entity.Balance, error) {
	return buc.bs.GetByID(ctx, id)
}

func (buc *ManagerUseCase) Transfer(ctx context.Context, idFrom, idTo *int64, amount decimal.Decimal) (transaction *entity.Transaction, err error) {
	defer func() {
		// Cancel transaction by service
		if err != nil && transaction != nil {
			multierr.AppendInto(&err, buc.ts.CancelByID(ctx, transaction.ID))
		}
	}()
	if idFrom == nil {
		return nil, domain.TransactionNilSourceErr
	}
	if idTo == nil {
		return nil, domain.TransactionNilDestinationErr
	}
	if amount.IsZero() {
		return nil, fmt.Errorf("idFrom = %q ; idFrom = %q ; amount = %s ; %w", *idFrom, *idTo, amount.String(), domain.ChangeBalanceByZeroAmountErr)
	}
	if amount.IsNegative() {
		return nil, fmt.Errorf("idFrom = %q ; idFrom = %q ; amount = %s ; %w", *idFrom, *idTo, amount.String(), domain.NegativeAmountTransactionErr)
	}
	if *idFrom == *idTo {
		return nil, fmt.Errorf("idFrom = %q ; idFrom = %q ; %w", *idFrom, *idTo, domain.TransactionSourceDestinationAreEqualErr)
	}
	// Check existence of idFrom balance
	_, err = buc.bs.GetByID(ctx, *idFrom)
	if err != nil {
		// TODO: wrap NotFoundErr!
		return nil, err
	}
	// Create transaction
	transaction, err = buc.ts.CreateDefaultTransaction(ctx, idFrom, idTo, amount, entity.TypeTransfer)
	// Apply transaction by producer
	if err == nil {
		err = buc.producer.ApplyTransaction(*transaction)
	}
	return transaction, err
}

func (buc *ManagerUseCase) ChangeAmount(ctx context.Context, id *int64, amount decimal.Decimal) (transaction *entity.Transaction, err error) {
	defer func() {
		// Cancel transaction by service
		if err != nil && transaction != nil {
			multierr.AppendInto(&err, buc.ts.CancelByID(ctx, transaction.ID))
		}
	}()
	if id == nil {
		return nil, domain.TransactionNilSourceOrDestinationErr
	}
	if amount.IsZero() {
		return nil, fmt.Errorf("idFrom = %q ; amount = %s ; %w", *id, amount.String(), domain.ChangeBalanceByZeroAmountErr)
	}
	// Create transaction
	if amount.IsPositive() {
		transaction, err = buc.ts.CreateDefaultTransaction(ctx, nil, id, amount, entity.TypeOuterIncreasing)
	} else {
		transaction, err = buc.ts.CreateDefaultTransaction(ctx, id, nil, amount.Neg(), entity.TypeOuterDecreasing)
	}
	// Apply transaction by producer
	if err == nil {
		err = buc.producer.ApplyTransaction(*transaction)
	}
	return transaction, err
}

// TODO: fix amount Sprintfs
func (buc *ManagerUseCase) PayForService(ctx context.Context, id *int64, amount decimal.Decimal) (transaction *entity.Transaction, err error) {
	defer func() {
		// Cancel transaction by service
		if err != nil && transaction != nil {
			multierr.AppendInto(&err, buc.ts.CancelByID(ctx, transaction.ID))
		}
	}()
	if id == nil {
		return nil, domain.TransactionNilSourceErr
	}
	if amount.IsZero() {
		return nil, fmt.Errorf("idFrom = %q ; amount = %s ; %w", *id, amount.String(), domain.ChangeBalanceByZeroAmountErr)
	}
	if amount.IsNegative() {
		return nil, fmt.Errorf("idFrom = %q ; amount = %s ; %w", *id, amount.String(), domain.NegativeAmountTransactionErr)
	}
	// Check existence of idFrom balance
	_, err = buc.bs.GetByID(ctx, *id)
	if err != nil {
		// TODO: wrap by NotFoundErr
		return nil, err
	}
	// Create transaction
	transaction, err = buc.ts.CreateDefaultTransaction(ctx, id, nil, amount.Neg(), entity.TypePayment)
	// Apply transaction by producer
	if err == nil {
		err = buc.producer.ApplyTransaction(*transaction)
	}
	return transaction, err
}
