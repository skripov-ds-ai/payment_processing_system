package entity

// Balance is balance entity
// Please, do not use float32 for money operations in production!
type Balance struct {
	ID     string  `json:"user_id"`
	Amount float32 `json:"amount"`
}
