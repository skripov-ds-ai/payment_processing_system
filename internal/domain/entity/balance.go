package entity

type Balance struct {
	UserUUID string `json:"user_uuid"`
	Amount   int64  `json:"amount"`
}
