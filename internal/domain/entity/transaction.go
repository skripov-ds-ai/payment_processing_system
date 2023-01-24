package entity

import "time"

type (
	TransactionStatus string
	TransactionType   string
)

const (
	StatusCancelled  TransactionStatus = "cancelled"
	StatusCreated    TransactionStatus = "created"
	StatusCompleted  TransactionStatus = "completed"
	StatusProcessing TransactionStatus = "processing"

	TypeOuterIncreasing TransactionType = "increasing"
	TypeOuterDecreasing TransactionType = "decreasing"
	TypeTransfer TransactionType = "transfer"
	TypePayment TransactionType = "payment"
)

// Transaction is entity to represent transaction in payment processing system
// Please, do not use float32 for money operations in production!
type Transaction struct {
	ID            string            `json:"id"`
	SourceID      *string           `json:"source_id"`
	DestinationID *string           `json:"destination_id"`
	Amount        float32           `json:"amount"`
	Type          TransactionType   `json:"type"`
	DateTime time.Time         `json:"date_time"`
	Status   TransactionStatus `json:"status"`
}
