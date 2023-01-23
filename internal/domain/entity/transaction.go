package entity

type Transaction struct {
	Source      *string `json:"source"`
	Destination *string `json:"destination"`
	Amount      int64   `json:"amount"`
	Type        string  `json:"type"`
}
