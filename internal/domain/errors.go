package domain

import (
	"errors"
	"fmt"
)

var (
	BalanceWasNotIncreased = fmt.Errorf("balance was not increased")
	BalanceWasNotDecreased = fmt.Errorf("balance was not decreased")
)

var TransactionSourceDestinationAreEqualErr = errors.New("source and destination of transaction are equal")

var (
	ChangeBalanceByZeroAmountErr = errors.New("changing balance by zero amount")

	TransactionNilSourceDestinationErr = errors.New("source and destination of transaction are nil")
	ZeroAmountTransactionErr           = errors.New("transaction amount is zero")
)
