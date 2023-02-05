package messages

import "payment_processing_system/internal/domain/entity"

// ApplyTransactionEvent is event to apply transaction on changing balance/balances
type ApplyTransactionEvent struct {
	Transaction entity.Transaction `json:"transaction"`
}

// CancelTransactionEvent is event to mark transaction status "cancelled"
type CancelTransactionEvent struct {
	TransactionID uint64 `json:"transaction_id"`
}

// CompleteTransactionEvent is event to mark transaction
// status "completed" and finish transaction
type CompleteTransactionEvent struct {
	TransactionID uint64 `json:"transaction_id"`
}
