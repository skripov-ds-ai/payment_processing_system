package entity

type Transaction struct {
	SourceID      *string `json:"source_id"`
	DestinationID *string `json:"destination_id"`
	Amount        int64   `json:"amount"`
	Type          string  `json:"type"`
}
