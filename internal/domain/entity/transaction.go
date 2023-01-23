package entity

// Transaction is entity to represent transaction in payment processing system
// Please, do not use float32 for money operations in production!
type Transaction struct {
	SourceID      *string `json:"source_id"`
	DestinationID *string `json:"destination_id"`
	Amount        float32 `json:"amount"`
	Type          string  `json:"type"`
}
