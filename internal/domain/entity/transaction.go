package entity

import (
	decimal "github.com/jackc/pgx-shopspring-decimal"
	"time"
)

type (
	TransactionStatus string
	TransactionType   string
)

const (
	StatusCancelled   TransactionStatus = "cancelled"
	StatusCreated     TransactionStatus = "created"
	StatusCompleted   TransactionStatus = "completed"
	StatusProcessing  TransactionStatus = "processing"
	StatusShouldRetry TransactionStatus = "should_retry"
	// TODO: add this to logic!
	StatusCannotApply TransactionStatus = "cannot_apply"

	TypeOuterIncreasing TransactionType = "increasing"
	TypeOuterDecreasing TransactionType = "decreasing"
	TypeTransfer        TransactionType = "transfer"
	TypePayment         TransactionType = "payment"
)

// Transaction is entity to represent transaction in payment processing system
// Please, do not use float32 for money operations in production!
type Transaction struct {
	ID              int64             `json:"id"`
	SourceID        *int64            `json:"source_id"`
	DestinationID   *int64            `json:"destination_id"`
	Amount          decimal.Decimal   `json:"amount"`
	Type            TransactionType   `json:"type"`
	DateTimeCreated time.Time         `json:"date_time_created"`
	DateTimeUpdated time.Time         `json:"date_time_updated"`
	Status          TransactionStatus `json:"status"`
}
