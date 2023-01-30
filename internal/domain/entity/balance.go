package entity

// Balance is balance entity
// Please, do not use float32 for money operations in production!
type Balance struct {
	ID     int64 `json:"user_id"`
	Amount int64 `json:"amount"`
}
