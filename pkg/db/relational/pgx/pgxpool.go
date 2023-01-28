package pgx

import (
	"context"
	"fmt"
	"log"
	retry "payment_processing_system/pkg/db"
	"payment_processing_system/pkg/db/relational"

	"github.com/jackc/pgx/v5/pgxpool"
)

// https://pkg.go.dev/github.com/jackc/pgx/v5/pgxpool#hdr-Creating_a_Pool

// type Client interface {
//	Begin(context.Context) (pgx.Tx, error)
//	BeginFunc(ctx context.Context, f func(pgx.Tx) error) error
//	BeginTxFunc(ctx context.Context, txOptions pgx.TxOptions, f func(pgx.Tx) error) error
//	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
//	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
//	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
//}

// NewClient for pgxpool
func NewClient(ctx context.Context, cfg *relational.SQLConnectConfig, retryCfg *retry.RetryConfig) (pool *pgxpool.Pool, err error) {
	dsn := fmt.Sprintf("postgresql://%s", cfg.DSN())
	err = retry.Retry(func() error {
		ctx, cancel := context.WithTimeout(ctx, retryCfg.PingDelay)
		defer cancel()

		pgxCfg, err := pgxpool.ParseConfig(dsn)
		if err != nil {
			log.Printf("Unable to parse config: %v\n", err)
			return err
		}

		pool, err = pgxpool.NewWithConfig(ctx, pgxCfg)
		if err != nil {
			log.Println("Failed to connect to postgres... Going to do the next attempt")
			return err
		}
		return nil
	}, retryCfg)
	return pool, err
}
