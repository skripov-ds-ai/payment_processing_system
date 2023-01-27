package pgx

import (
	"context"
	"fmt"
	"payment_processing_system/internal/domain/entity"
	"payment_processing_system/pkg/logger"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type transactionStorage struct {
	tableScheme  string
	queryBuilder sq.StatementBuilderType
	pool         *pgxpool.Pool
	logger       *logger.Logger
}

func NewTransactionStorage(pool *pgxpool.Pool, logger *logger.Logger) *transactionStorage {
	tableScheme := "public.transaction"
	return &transactionStorage{
		tableScheme:  tableScheme,
		queryBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		pool:         pool,
		logger:       logger,
	}
}

func (ts *transactionStorage) UpdateStatusByID(ctx context.Context, id, status string) error {
	sql, args, buildErr := ts.queryBuilder.
		Update(ts.tableScheme).Set("status", status).
		Set("date_time_updated", time.Now()).
		Where(sq.Eq{"id": id}).ToSql()
	ts.logger.Info("update sql",
		zap.String("table", ts.tableScheme),
		zap.String("sql", sql),
		zap.String("args", fmt.Sprintf("%v", args)))
	if buildErr != nil {
		// TODO: add wrapping
		// buildErr
		return buildErr
	}
	if exec, execErr := ts.pool.Exec(ctx, sql, args...); execErr != nil {
		// TODO: wrap
		return execErr
	} else if exec.RowsAffected() == 0 || !exec.Update() {
		// TODO: create err
		return fmt.Errorf("transaction status was not updated")
	}
	return nil
}

func (ts *transactionStorage) Create(ctx context.Context, transaction entity.Transaction) (*string, error) {
	sql, args, buildErr := ts.queryBuilder.
		Insert(ts.tableScheme).Columns(
		"amount", "source_id", "destination_id",
		"type", "date_time_created", "datetime_updated", "status").
		Values(
			transaction.Amount, transaction.SourceID, transaction.DestinationID,
			transaction.Type, transaction.DateTimeCreated, transaction.DateTimeUpdated, transaction.Status).
		Suffix("RETURNING \"id\"").
		ToSql()
	ts.logger.Info("insert sql",
		zap.String("table", ts.tableScheme),
		zap.String("sql", sql),
		zap.String("args", fmt.Sprintf("%v", args)))
	if buildErr != nil {
		// TODO: add wrapping
		// buildErr
		return nil, buildErr
	}
	err := ts.pool.QueryRow(ctx, sql, args...).Scan(&transaction.ID)
	if err != nil {
		return &transaction.ID, err
	}
	return &transaction.ID, nil
}

func (ts *transactionStorage) GetByID(ctx context.Context, id string) (*entity.Transaction, error) {
	sql, args, buildErr := ts.queryBuilder.
		Select("id", "source_id", "destination_id", "amount",
			"type", "date_time_created", "date_time_updated", "status").
		From(ts.tableScheme).Where(sq.Eq{"id": id}).ToSql()
	ts.logger.Info("select sql",
		zap.String("table", ts.tableScheme),
		zap.String("sql", sql),
		zap.String("args", fmt.Sprintf("%v", args)))
	if buildErr != nil {
		// TODO: add wrapping
		// buildErr
		return nil, buildErr
	}
	var obj entity.Transaction
	err := ts.pool.QueryRow(ctx, sql, args...).Scan(
		&obj.ID,
		&obj.SourceID,
		&obj.DestinationID,
		&obj.Amount,
		&obj.Type,
		&obj.DateTimeCreated,
		&obj.DateTimeUpdated,
		&obj.Status)
	if err != nil {
		// TODO: wrap error
		return nil, err
	}
	// bs.pool.Query()
	return &obj, nil
}
