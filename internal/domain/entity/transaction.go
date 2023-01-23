package entity

type Transaction struct {
	UserFromUUID string  `json:"user_from_uuid"`
	UserToUUID   *string `json:"user_to_uuid"`
	Amount       int64   `json:"amount"`
	ServiceID    *int64  `json:"service_id"`
	OrderID      *int64  `json:"order_id"`
}
