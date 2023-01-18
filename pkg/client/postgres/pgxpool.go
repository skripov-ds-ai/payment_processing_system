package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	retry "payment_processing_system/pkg/client"
	"time"
)

type pgConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

// NewPgConfig creates new pg config instance
func NewPgConfig(username, password, host, port, database string) *pgConfig {
	return &pgConfig{
		Username: username,
		Password: password,
		Host:     host,
		Port:     port,
		Database: database,
	}
}

// https://pkg.go.dev/github.com/jackc/pgx/v5/pgxpool#hdr-Creating_a_Pool

type Client interface {
	Begin(context.Context) (pgx.Tx, error)
	BeginFunc(ctx context.Context, f func(pgx.Tx) error) error
	BeginTxFunc(ctx context.Context, txOptions pgx.TxOptions, f func(pgx.Tx) error) error
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
}

func DoWithAttempts(fn func() error, retryCfg *retry.RetryConfig) error {
	var err error
	maxAttempts := retryCfg.MaxAttempts
	for maxAttempts > 0 {
		if err = fn(); err == nil {
			return nil
		}
		time.Sleep(retryCfg.Delay)
		maxAttempts--
	}
	return err
}

func NewClient(ctx context.Context, retryCfg *retry.RetryConfig, cfg *pgConfig) (pool *pgxpool.Pool, err error) {
	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		cfg.Username, cfg.Password,
		cfg.Host, cfg.Port, cfg.Database,
	)

	err = DoWithAttempts(func() error {
		ctx, cancel := context.WithTimeout(ctx, retryCfg.Delay)
		defer cancel()

		pgxCfg, err := pgxpool.ParseConfig(dsn)
		if err != nil {
			log.Println("Unable to parse config: %v\n", err)
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
