package pgx

import (
	"context"
	"fmt"
	"payment_processing_system/internal/domain"
	"payment_processing_system/internal/domain/entity"
	"payment_processing_system/pkg/db/relational/pgx"
	"payment_processing_system/pkg/logger"

	"github.com/shopspring/decimal"

	"go.uber.org/zap"

	sq "github.com/Masterminds/squirrel"
)

type balanceStorage struct {
	tableScheme  string
	queryBuilder sq.StatementBuilderType
	pool         pgx.PgxIface
	logger       *logger.Logger
}

func NewBalanceStorage(pool pgx.PgxIface, logger *logger.Logger) *balanceStorage {
	tableScheme := "public.balance"
	return &balanceStorage{
		tableScheme:  tableScheme,
		queryBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		pool:         pool,
		logger:       logger,
	}
}

// TODO
func (bs *balanceStorage) IncreaseAmount(ctx context.Context, id int64, amount decimal.Decimal) error {
	onConflict := "ON CONFLICT DO UPDATE SET amount = amount + ?"
	sql, args, buildErr := bs.queryBuilder.
		Insert(bs.tableScheme).Columns("id", "amount").
		Values(id, amount).Suffix(onConflict, amount).ToSql()
	bs.logger.Info("upsert sql",
		zap.String("table", bs.tableScheme),
		zap.String("sql", sql),
		zap.String("args", fmt.Sprintf("%v", args)))
	if buildErr != nil {
		return buildErr
	}
	if exec, execErr := bs.pool.Exec(ctx, sql, args...); execErr != nil {
		return execErr
	} else if exec.RowsAffected() == 0 {
		return domain.BalanceWasNotIncreasedErr
	}
	return nil
}

func (bs *balanceStorage) DecreaseAmount(ctx context.Context, id int64, amount decimal.Decimal) error {
	sql, args, buildErr := bs.queryBuilder.
		Update(bs.tableScheme).
		Set("amount", fmt.Sprintf("amount + %s", amount.String())).
		Where(sq.Eq{"id": id}).ToSql()
	bs.logger.Info("update sql",
		zap.String("table", bs.tableScheme),
		zap.String("sql", sql),
		zap.String("args", fmt.Sprintf("%v", args)))
	if buildErr != nil {
		return buildErr
	}
	if exec, execErr := bs.pool.Exec(ctx, sql, args...); execErr != nil {
		return execErr
	} else if exec.RowsAffected() == 0 || !exec.Insert() {
		return domain.BalanceWasNotDecreasedErr
	}
	return nil
}

func (bs *balanceStorage) GetByID(ctx context.Context, id int64) (*entity.Balance, error) {
	sql, args, buildErr := bs.queryBuilder.
		Select("id", "amount").
		From(bs.tableScheme).Where(sq.Eq{"id": id}).ToSql()
	bs.logger.Info("select sql",
		zap.String("table", bs.tableScheme),
		zap.String("sql", sql),
		zap.String("args", fmt.Sprintf("%v", args)))
	if buildErr != nil {
		return nil, buildErr
	}
	var obj entity.Balance
	err := bs.pool.QueryRow(ctx, sql, args...).Scan(
		&obj.ID,
		&obj.Amount)
	if err != nil {
		return nil, err
	}
	return &obj, nil
}
