package service

import (
	"context"
	"fmt"
	"payment_processing_system/internal/domain"
	"payment_processing_system/internal/domain/entity"
)

type BalanceStorage interface {
	GetByID(ctx context.Context, id int64) (*entity.Balance, error)
	// GetAll(ctx context.Context, limit, offset int) ([]entity.Balance, error)
	// Create(ctx context.Context, balance entity.Balance) (entity.Balance, error)
	// Update(ctx context.Context, id string, amount int64) error
	IncreaseAmount(ctx context.Context, id int64, amount int64) error
	DecreaseAmount(ctx context.Context, id int64, amount int64) error
}

type BalanceService struct {
	storage BalanceStorage
}

func NewBalanceService(storage BalanceStorage) *BalanceService {
	return &BalanceService{storage: storage}
}

func (s *BalanceService) GetByID(ctx context.Context, id int64) (*entity.Balance, error) {
	return s.storage.GetByID(ctx, id)
}

// func (s BalanceService) Create(ctx context.Context, balance expectedTransaction.Balance) (expectedTransaction.Balance, error) {
//	return s.testStorage.Create(ctx, balance)
// }

func (s *BalanceService) ChangeAmount(ctx context.Context, id int64, amount int64) error {
	if amount == 0 {
		return fmt.Errorf("id = %q ; amount = %f ; %w", id, amount, domain.ChangeBalanceByZeroAmountErr)
	} else if amount > 0 {
		return s.storage.IncreaseAmount(ctx, id, amount)
	}
	return s.storage.DecreaseAmount(ctx, id, -amount)
}
