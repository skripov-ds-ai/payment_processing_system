package service

import (
	"context"
	"payment_processing_system/internal/domain"
	"payment_processing_system/internal/domain/entity"
	"payment_processing_system/internal/utils"
	"time"
)

type TransactionStorage interface {
	GetByID(ctx context.Context, id string) (*entity.Transaction, error)
	Create(ctx context.Context, transaction entity.Transaction) (string, error)
	UpdateStatusByID(ctx context.Context, id string, status entity.TransactionStatus) error
}

type TransactionService struct {
	storage TransactionStorage
}

func NewTransactionService(storage TransactionStorage) *TransactionService {
	return &TransactionService{storage: storage}
}

func (t *TransactionService) GetByID(ctx context.Context, id string) (*entity.Transaction, error) {
	return t.storage.GetByID(ctx, id)
}

func (t *TransactionService) CancelByID(ctx context.Context, id string) error {
	return t.storage.UpdateStatusByID(ctx, id, entity.StatusCancelled)
}

func (t *TransactionService) ProcessingByID(ctx context.Context, id string) error {
	return t.storage.UpdateStatusByID(ctx, id, entity.StatusProcessing)
}

func (t *TransactionService) CompletedByID(ctx context.Context, id string) error {
	return t.storage.UpdateStatusByID(ctx, id, entity.StatusCompleted)
}

func (t *TransactionService) ShouldRetryByID(ctx context.Context, id string) error {
	return t.storage.UpdateStatusByID(ctx, id, entity.StatusShouldRetry)
}

func (t *TransactionService) CannotApplyByID(ctx context.Context, id string) error {
	return t.storage.UpdateStatusByID(ctx, id, entity.StatusCannotApply)
}

func (t *TransactionService) CreateDefaultTransaction(ctx context.Context, sourceID, destinationID *string, amount float32, ttype entity.TransactionType) (string, error) {
	if utils.IsZero(amount) {
		return "", domain.ZeroAmountTransactionErr
	}
	if amount < 0 {
		return "", domain.NegativeAmountTransactionErr
	}
	if sourceID == nil && destinationID == nil {
		return "", domain.TransactionNilSourceAndDestinationErr
	}
	if sourceID == destinationID || sourceID != nil && destinationID != nil && *sourceID == *destinationID {
		return "", domain.TransactionSourceDestinationAreEqualErr
	}
	now := time.Now()
	transaction := entity.Transaction{
		Amount: amount, SourceID: sourceID, DestinationID: destinationID,
		Status: entity.StatusCreated, DateTime: now, Type: ttype}
	return t.storage.Create(ctx, transaction)
}
