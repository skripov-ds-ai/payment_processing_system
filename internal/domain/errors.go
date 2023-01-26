package domain

import (
	"errors"
)

var (
	BalanceWasNotIncreased                  = errors.New("balance was not increased")
	BalanceWasNotDecreased                  = errors.New("balance was not decreased")
	ChangeBalanceByZeroAmountErr            = errors.New("changing balance by zero amount")
	TransactionSourceDestinationAreEqualErr = errors.New("source and destination of transaction are equal")
	TransactionNilSourceDestinationErr      = errors.New("source and destination of transaction are nil")
	ZeroAmountTransactionErr                = errors.New("transaction amount is zero")
	NegativeAmountTransactionErr            = errors.New("transaction amount is negative")
)
