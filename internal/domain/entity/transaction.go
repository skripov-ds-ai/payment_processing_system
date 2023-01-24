package entity

import "time"

// Transaction is entity to represent transaction in payment processing system
// Please, do not use float32 for money operations in production!
type Transaction struct {
	ID            string    `json:"id"`
	SourceID      *string   `json:"source_id"`
	DestinationID *string   `json:"destination_id"`
	Amount        float32   `json:"amount"`
	Type          string    `json:"type"`
	DateTime      time.Time `json:"date_time"`
	Status        string    `json:"status"`
}
