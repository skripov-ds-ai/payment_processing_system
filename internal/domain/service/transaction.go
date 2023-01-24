package service

type TransactionStorage interface {
	//GetByID(ctx context.Context, id string) (*entity.Transaction, error)
	//IncreaseAmount(ctx context.Context, id string, amount float32) error
	//DecreaseAmount(ctx context.Context, id string, amount float32) error
}

type TransactionService struct {
	storage TransactionStorage
}
