package entity

type Balance struct {
	ID     string `json:"user_id"`
	Amount int64  `json:"amount"`
}
