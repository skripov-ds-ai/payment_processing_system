package service

import (
	"context"
	"fmt"
	"payment_processing_system/internal/domain/entity"
)

type BalanceStorage interface {
	GetByID(ctx context.Context, id string) (*entity.Balance, error)
	// GetAll(ctx context.Context, limit, offset int) ([]entity.Balance, error)
	// Create(ctx context.Context, balance entity.Balance) (entity.Balance, error)
	// Update(ctx context.Context, id string, amount int64) error
	IncreaseAmount(ctx context.Context, id string, amount int64) error
	DecreaseAmount(ctx context.Context, id string, amount int64) error
}

type BalanceService struct {
	storage BalanceStorage
}

func (s BalanceService) GetByID(ctx context.Context, id string) (*entity.Balance, error) {
	return s.storage.GetByID(ctx, id)
}

// func (s BalanceService) Create(ctx context.Context, balance entity.Balance) (entity.Balance, error) {
//	return s.storage.Create(ctx, balance)
// }

func (s BalanceService) ChangeAmount(ctx context.Context, id string, amount int64) error {
	if amount == 0 {
		return fmt.Errorf("changing balance with id = %s by zero(amount = %d)", id, amount)
	} else if amount > 0 {
		return s.storage.IncreaseAmount(ctx, id, amount)
	}
	return s.storage.DecreaseAmount(ctx, id, -amount)
}
