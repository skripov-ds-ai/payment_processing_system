package usecase

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"math"
	"payment_processing_system/internal/domain/entity"
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

type Producer interface {
	PublishMessage(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}

type BalanceUseCase struct {
	bs       BalanceService
	ts       TransactionService
	producer Producer
}

func NewBalanceUseCase(bs BalanceService, ts TransactionService, producer Producer) *BalanceUseCase {
	return &BalanceUseCase{bs: bs, ts: ts, producer: producer}
}

func (buc *BalanceUseCase) ChangeAmount(ctx context.Context, id string, amount float32) error {
	if math.Abs(float64(amount)) < 1e-9 {
		return fmt.Errorf("changing balance with id = %s by zero(amount = %f)", id, amount)
	}
	// TODO: change logic!
	var transactionID string
	var err error
	if amount > 0 {
		transactionID, err = buc.ts.CreateDefaultTransaction(ctx, nil, &id, amount, entity.TypeOuterIncreasing)
	} else {
		transactionID, err = buc.ts.CreateDefaultTransaction(ctx, &id, nil, -amount, entity.TypeOuterDecreasing)
	}
	if err != nil {
		// TODO: wrap err
		_ = buc.ts.CancelByID(ctx, transactionID)
		return err
	}
	_ = buc.ts.ProcessingByID(ctx, transactionID)
	err = buc.bs.ChangeAmount(ctx, id, amount)
	if err != nil {
		_ = buc.ts.ShouldRetryByID(ctx, id)
		return err
	}
	_ = buc.ts.CompletedByID(ctx, id)
	return nil
}
