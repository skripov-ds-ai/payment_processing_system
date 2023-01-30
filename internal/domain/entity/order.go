package entity

type Order struct {
	TransactionID int64 `json:"transaction_id"`
	OrderID       int64 `json:"order_id"`
	ServiceID     int64 `json:"service_id"`
}
