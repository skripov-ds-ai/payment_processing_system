package service

import "errors"

var (
	ChangeBalanceByZeroAmountErr            = errors.New("changing balance by zero amount")
	TransactionSourceDestinationAreEqualErr = errors.New("source and destination of transaction are equal")
	TransactionNilSourceDestinationErr      = errors.New("source and destination of transaction are nil")
	ZeroAmountTransactionErr                = errors.New("transaction amount is zero")
)
