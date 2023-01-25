package pgx

import (
	"context"
	"fmt"
	"payment_processing_system/internal/domain/entity"
	"payment_processing_system/pkg/logger"

	"go.uber.org/zap"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type balanceStorage struct {
	tableScheme  string
	queryBuilder sq.StatementBuilderType
	pool         *pgxpool.Pool
	logger       *logger.Logger
}

func NewBalanceStorage(pool *pgxpool.Pool, logger *logger.Logger) *balanceStorage {
	tableScheme := "public.balance"
	return &balanceStorage{
		tableScheme:  tableScheme,
		queryBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		pool:         pool,
		logger:       logger,
	}
}

func (bs *balanceStorage) IncreaseAmount(ctx context.Context, id string, amount float32) error {
	onConflict := "ON CONFLICT DO UPDATE SET amount = amount + ?"
	sql, args, buildErr := bs.queryBuilder.
		Insert(bs.tableScheme).Columns("id", "amount").
		Values(id, amount).Suffix(onConflict, amount).ToSql()
	bs.logger.Info("upsert sql",
		zap.String("table", bs.tableScheme),
		zap.String("sql", sql),
		zap.String("args", fmt.Sprintf("%v", args)))
	if buildErr != nil {
		// TODO: add wrapping
		// buildErr
		return buildErr
	}
	if exec, execErr := bs.pool.Exec(ctx, sql, args...); execErr != nil {
		// TODO: wrap
		return execErr
	} else if exec.RowsAffected() == 0 || !exec.Insert() {
		// TODO: create err
		return fmt.Errorf("balance was not increased")
	}
	return nil
}

func (bs *balanceStorage) DecreaseAmount(ctx context.Context, id string, amount float32) error {
	sql, args, buildErr := bs.queryBuilder.
		Update(bs.tableScheme).
		Set("amount", fmt.Sprintf("amount + %f", amount)).
		Where(sq.Eq{"id": id}).ToSql()
	bs.logger.Info("update sql",
		zap.String("table", bs.tableScheme),
		zap.String("sql", sql),
		zap.String("args", fmt.Sprintf("%v", args)))
	if buildErr != nil {
		// TODO: add wrapping
		// buildErr
		return buildErr
	}
	if exec, execErr := bs.pool.Exec(ctx, sql, args...); execErr != nil {
		// TODO: wrap
		return execErr
	} else if exec.RowsAffected() == 0 || !exec.Insert() {
		// TODO: create err
		return fmt.Errorf("balance was not decreased")
	}
	return nil
}

func (bs *balanceStorage) GetByID(ctx context.Context, id string) (*entity.Balance, error) {
	sql, args, buildErr := bs.queryBuilder.
		Select("id", "amount").
		From(bs.tableScheme).Where(sq.Eq{"id": id}).ToSql()
	bs.logger.Info("select sql",
		zap.String("table", bs.tableScheme),
		zap.String("sql", sql),
		zap.String("args", fmt.Sprintf("%v", args)))
	if buildErr != nil {
		// TODO: add wrapping
		// buildErr
		return nil, buildErr
	}
	var obj entity.Balance
	err := bs.pool.QueryRow(ctx, sql, args...).Scan(
		&obj.ID,
		&obj.Amount)
	if err != nil {
		// TODO: wrap error
		return nil, err
	}
	// bs.pool.Query()
	return &obj, nil
}
