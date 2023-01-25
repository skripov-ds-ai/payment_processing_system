package service

import "errors"

var (
	ChangeBalanceByZeroAmountErr = errors.New("changing balance by zero amount")
)
