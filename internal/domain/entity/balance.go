package entity

type Balance struct {
	UserID string `json:"user_id"`
	Amount int64  `json:"amount"`
}
