package converter

import "github.com/shopspring/decimal"

type apiResult struct {
	Date       string `json:"date"`
	Historical string `json:"historical"`
	Info       struct {
		Rate      decimal.Decimal `json:"rate"`
		Timestamp int             `json:"timestamp"`
	} `json:"info"`
	Query struct {
		Amount int    `json:"amount"`
		From   string `json:"from"`
		To     string `json:"to"`
	} `json:"query"`
	Result  decimal.Decimal `json:"result"`
	Success bool            `json:"success"`
}
