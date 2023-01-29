package domain

import (
	"errors"
)

var (
	BalanceWasNotIncreased                  = errors.New("balance was not increased")
	BalanceWasNotDecreased                  = errors.New("balance was not decreased")
	ChangeBalanceByZeroAmountErr            = errors.New("changing balance by zero amount")
	TransactionSourceDestinationAreEqualErr = errors.New("source and destination of transaction are equal")
	TransactionNilSourceAndDestinationErr   = errors.New("source and destination of transaction are nil")
	TransactionNilSourceErr                 = errors.New("source of transaction are nil")
	TransactionNilDestinationErr            = errors.New("destination of transaction are nil")
	TransactionNilSourceOrDestinationErr    = errors.New("source or destination of transaction are nil")
	ZeroAmountTransactionErr                = errors.New("transaction amount is zero")
	NegativeAmountTransactionErr            = errors.New("transaction amount is negative")
	UnknownTransactionTypeErr               = errors.New("unknown transaction type")
)
