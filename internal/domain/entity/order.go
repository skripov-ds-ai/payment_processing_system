package entity

type Order struct {
	TransactionID string `json:"transaction_id"`
	OrderID       string `json:"order_id"`
	ServiceID     string `json:"service_id"`
}
