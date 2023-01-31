package service

import (
	"context"
	"payment_processing_system/internal/domain"
	"payment_processing_system/internal/domain/entity"
	"time"

	"github.com/shopspring/decimal"
)

type TransactionStorage interface {
	GetByID(ctx context.Context, id int64) (*entity.Transaction, error)
	Create(ctx context.Context, transaction entity.Transaction) (*entity.Transaction, error)
	UpdateStatusByID(ctx context.Context, id int64, status entity.TransactionStatus) error
}

type TransactionService struct {
	storage TransactionStorage
}

func NewTransactionService(storage TransactionStorage) *TransactionService {
	return &TransactionService{storage: storage}
}

func (t *TransactionService) GetByID(ctx context.Context, id int64) (*entity.Transaction, error) {
	return t.storage.GetByID(ctx, id)
}

func (t *TransactionService) CancelByID(ctx context.Context, id int64) error {
	return t.storage.UpdateStatusByID(ctx, id, entity.StatusCancelled)
}

func (t *TransactionService) ProcessingByID(ctx context.Context, id int64) error {
	return t.storage.UpdateStatusByID(ctx, id, entity.StatusProcessing)
}

func (t *TransactionService) CompletedByID(ctx context.Context, id int64) error {
	return t.storage.UpdateStatusByID(ctx, id, entity.StatusCompleted)
}

func (t *TransactionService) ShouldRetryByID(ctx context.Context, id int64) error {
	return t.storage.UpdateStatusByID(ctx, id, entity.StatusShouldRetry)
}

func (t *TransactionService) CannotApplyByID(ctx context.Context, id int64) error {
	return t.storage.UpdateStatusByID(ctx, id, entity.StatusCannotApply)
}

func (t *TransactionService) CreateDefaultTransaction(ctx context.Context, sourceID, destinationID *int64, amount decimal.Decimal, ttype entity.TransactionType) (*entity.Transaction, error) {
	if amount.IsZero() {
		return nil, domain.ZeroAmountTransactionErr
	}
	if amount.IsNegative() {
		return nil, domain.NegativeAmountTransactionErr
	}
	if sourceID == nil && destinationID == nil {
		return nil, domain.TransactionNilSourceAndDestinationErr
	}
	if sourceID == destinationID || sourceID != nil && destinationID != nil && *sourceID == *destinationID {
		return nil, domain.TransactionSourceDestinationAreEqualErr
	}
	now := time.Now()
	transaction := entity.Transaction{
		Amount: amount, SourceID: sourceID, DestinationID: destinationID,
		Status: entity.StatusCreated, DateTimeCreated: now, DateTimeUpdated: now, TType: ttype}
	return t.storage.Create(ctx, transaction)
}
