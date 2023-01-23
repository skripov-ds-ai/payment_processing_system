package service

import "payment_processing_system/internal/domain/entity"

type BalanceStorage interface {
	GetByUUID(uuid string) (entity.Balance, error)
	F() // TODO: remove
}

type BalanceService struct {
	storage BalanceStorage
}

func (s BalanceService) GetByUUID(uuid string) (entity.Balance, error) {
	return s.storage.GetByUUID(uuid)
}
