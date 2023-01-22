package transaction

type Transaction struct {
	UserFromUUID string
	UserToUUID   *string
	Amount       int64
	ServiceID    *int64
	OrderID      *int64
}
