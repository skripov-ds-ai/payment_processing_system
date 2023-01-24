package pgx

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type transactionStorage struct {
	tableScheme  string
	queryBuilder sq.StatementBuilderType
	pool         *pgxpool.Pool
	logger       *zap.Logger
}

func NewTransactionStorage(pool *pgxpool.Pool, logger *zap.Logger) *transactionStorage {
	tableScheme := "public.balance"
	return &transactionStorage{
		tableScheme:  tableScheme,
		queryBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		pool:         pool,
		logger:       logger,
	}
}
