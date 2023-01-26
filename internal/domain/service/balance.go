package service

import (
	"context"
	"fmt"
	"payment_processing_system/internal/domain/entity"
	"payment_processing_system/internal/utils"
)

type BalanceStorage interface {
	GetByID(ctx context.Context, id string) (*entity.Balance, error)
	// GetAll(ctx context.Context, limit, offset int) ([]entity.Balance, error)
	// Create(ctx context.Context, balance entity.Balance) (entity.Balance, error)
	// Update(ctx context.Context, id string, amount int64) error
	IncreaseAmount(ctx context.Context, id string, amount float32) error
	DecreaseAmount(ctx context.Context, id string, amount float32) error
}

type BalanceService struct {
	storage BalanceStorage
}

func NewBalanceService(storage BalanceStorage) *BalanceService {
	return &BalanceService{storage: storage}
}

func (s *BalanceService) GetByID(ctx context.Context, id string) (*entity.Balance, error) {
	return s.storage.GetByID(ctx, id)
}

// func (s BalanceService) Create(ctx context.Context, balance entity.Balance) (entity.Balance, error) {
//	return s.storage.Create(ctx, balance)
// }

func (s *BalanceService) ChangeAmount(ctx context.Context, id string, amount float32) error {
	// TODO: fix check is zero
	if utils.IsZero(amount) {
		return fmt.Errorf("id = %q ; amount = %f ; %w", id, amount, ChangeBalanceByZeroAmountErr)
	} else if amount > 0 {
		return s.storage.IncreaseAmount(ctx, id, amount)
	}
	return s.storage.DecreaseAmount(ctx, id, -amount)
}
