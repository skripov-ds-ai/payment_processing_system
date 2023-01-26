package entity

import "fmt"

var (
	BalanceWasNotIncreased = fmt.Errorf("balance was not increased")
	BalanceWasNotDecreased = fmt.Errorf("balance was not decreased")
)
