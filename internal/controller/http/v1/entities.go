package v1

import (
	"payment_processing_system/internal/domain/entity"

	"github.com/shopspring/decimal"
)

func (b *Balance) ToDomain() (*entity.Balance, error) {
	amount, err := decimal.NewFromString(b.Amount)
	if err != nil {
		return nil, err
	}
	return &entity.Balance{ID: b.Id, Amount: amount}, nil
}
