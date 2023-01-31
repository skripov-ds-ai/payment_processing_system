package entity

import "github.com/shopspring/decimal"

// Balance is balance entity
// Please, do not use float32 for money operations in production!
type Balance struct {
	ID     int64           `json:"user_id"`
	Amount decimal.Decimal `json:"amount"`
}
