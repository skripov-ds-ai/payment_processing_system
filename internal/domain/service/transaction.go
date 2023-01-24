package service

import (
	"context"
	"payment_processing_system/internal/domain/entity"
)

type TransactionStorage interface {
	GetByID(ctx context.Context, id string) (*entity.Transaction, error)
	Create(ctx context.Context, transaction entity.Transaction) error
	UpdateStatusByID(ctx context.Context, id, status string) error
}

type TransactionService struct {
	storage TransactionStorage
}

func (t TransactionService) GetByID(ctx context.Context, id string) (*entity.Transaction, error) {
	return t.storage.GetByID(ctx, id)
}
