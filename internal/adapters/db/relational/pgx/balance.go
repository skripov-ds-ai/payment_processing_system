package pgx

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/gommon/log"
	"payment_processing_system/internal/domain/entity"
)

type balanceStorage struct {
	tableScheme  string
	queryBuilder sq.StatementBuilderType
	pool         *pgxpool.Pool // TODO: move to interface
}

func NewBalanceStorage(pool *pgxpool.Pool) *balanceStorage {
	tableScheme := "public.balance"
	return &balanceStorage{
		tableScheme:  tableScheme,
		queryBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		pool:         pool,
	}
}

func (bs *balanceStorage) IncreaseAmount(ctx context.Context, id string, amount int64) error {
	column := fmt.Sprintf("%s.%s", bs.tableScheme, "amount")
	onConflict := fmt.Sprintf("ON CONFLICT DO UPDATE SET %s = %s + ?", column, column)
	sql, args, buildErr := bs.queryBuilder.
		Insert(bs.tableScheme).Columns("id", "amount").
		Values(id, amount).Suffix(onConflict, amount).ToSql()
	log.Info(fmt.Sprintf("table = %q ; sql = %q ; args = %q", bs.tableScheme, sql, args))
	if buildErr != nil {
		// TODO: add wrapping
		//buildErr
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

func (bs *balanceStorage) GetByID(ctx context.Context, id string) (*entity.Balance, error) {
	sql, args, buildErr := bs.queryBuilder.
		Select("id").
		Columns("amount").
		From(bs.tableScheme).Where(sq.Eq{"id": id}).ToSql()
	// TODO: add logging
	log.Info(fmt.Sprintf("table = %q ; sql = %q ; args = %q", bs.tableScheme, sql, args))
	if buildErr != nil {
		// TODO: add wrapping
		//buildErr
		return nil, buildErr
	}
	var obj entity.Balance
	err := bs.pool.QueryRow(ctx, sql, args...).Scan(
		&obj.UserID,
		&obj.Amount)
	if err != nil {
		// TODO: wrap error
		return nil, err
	}
	//bs.pool.Query()
	return &obj, nil
}
