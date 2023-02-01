package pgx

import (
	"context"
	"fmt"
	"payment_processing_system/internal/domain/entity"
	"payment_processing_system/pkg/db/relational/pgx"
	"payment_processing_system/pkg/logger"
	"time"

	sq "github.com/Masterminds/squirrel"
	"go.uber.org/zap"
)

type transactionStorage struct {
	tableScheme  string
	queryBuilder sq.StatementBuilderType
	pool         pgx.PgxIface
	logger       *logger.Logger
}

func NewTransactionStorage(pool pgx.PgxIface, logger *logger.Logger) *transactionStorage {
	tableScheme := "public.transaction"
	return &transactionStorage{
		tableScheme:  tableScheme,
		queryBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		pool:         pool,
		logger:       logger,
	}
}

// TODO
func (ts *transactionStorage) GetBalanceTransactions(ctx context.Context, balanceID int64, limit, offset uint64, orderBy string) ([]*entity.Transaction, error) {
	// select * from "transaction" as t where t.source_id = <balanceID> or t.destination_id = <balanceID> ORDER BY id;
	// non-optimized: select * from "transaction" as t where t.source_id = <balanceID> or t.destination_id = <balanceID> ORDER BY id LIMIT <limit> OFFSET <offset>;
	orderBys := make([]string, 0)
	orderBys = append(orderBys, orderBy)
	if orderBy != "id" {
		orderBys = append(orderBys, "id")
	}
	sql, args, buildErr := ts.queryBuilder.
		Select("id", "source_id", "destination_id", "amount",
			"ttype", "date_time_created", "date_time_updated", "status").
		From(ts.tableScheme).Where(
		sq.Or{
			sq.Eq{"source_id": balanceID},
			sq.Eq{"destination_id": balanceID},
		}).Limit(limit).Offset(offset).OrderBy(orderBys...).ToSql()
	ts.logger.Info("select sql",
		zap.String("table", ts.tableScheme),
		zap.String("sql", sql),
		zap.String("args", fmt.Sprintf("%v", args)))
	if buildErr != nil {
		// TODO: add wrapping
		// buildErr
		return nil, buildErr
	}
	rows, err := ts.pool.Query(ctx, sql, args...)
	if err != nil {
		// TODO: add wrapping
		return nil, err
	}
	defer rows.Close()

	page := make([]*entity.Transaction, 0)
	for rows.Next() {
		transaction := entity.Transaction{}
		if err = rows.Scan(
			&transaction.ID,
			&transaction.SourceID,
			&transaction.DestinationID,
			&transaction.Amount,
			&transaction.TType,
			&transaction.DateTimeCreated,
			&transaction.DateTimeUpdated,
			&transaction.Status,
		); err != nil {
			// TODO: add wrapping
			return nil, err
		}

		page = append(page, &transaction)
	}
	return page, nil
}

func (ts *transactionStorage) UpdateStatusByID(ctx context.Context, id uint64, status entity.TransactionStatus) error {
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

func (ts *transactionStorage) Create(ctx context.Context, transaction entity.Transaction) (*entity.Transaction, error) {
	sql, args, buildErr := ts.queryBuilder.
		Insert(ts.tableScheme).Columns(
		"amount", "source_id", "destination_id",
		"ttype", "date_time_created", "datetime_updated", "status").
		Values(
			transaction.Amount, transaction.SourceID, transaction.DestinationID,
			transaction.TType, transaction.DateTimeCreated, transaction.DateTimeUpdated, transaction.Status).
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
		return &transaction, err
	}
	return &transaction, nil
}

func (ts *transactionStorage) GetByID(ctx context.Context, id uint64) (*entity.Transaction, error) {
	sql, args, buildErr := ts.queryBuilder.
		Select("id", "source_id", "destination_id", "amount",
			"ttype", "date_time_created", "date_time_updated", "status").
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
		&obj.TType,
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
