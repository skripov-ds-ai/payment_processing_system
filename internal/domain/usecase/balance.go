package usecase

import "payment_processing_system/internal/domain/entity"

type BalanceService interface {
	GetByUUID(uuid string) (entity.Balance, error)
	ChangeBalanceByUUID(uuid string, amount int64) error
}

type BalanceUseCase struct {
}
