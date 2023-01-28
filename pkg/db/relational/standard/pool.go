package standard

import (
	"context"
	"database/sql"
	"log"
	retry "payment_processing_system/pkg/db"
	"payment_processing_system/pkg/db/relational"
	"time"
	// importing MySQL
	// _ "github.com/go-sql-driver/mysql"
)

// PoolConfig to configure database/sql pool
type PoolConfig struct {
	MaxOpenConns    int
	MaxIdleConns    int
	MaxLifetime     time.Duration
	ConnMaxIdleTime time.Duration
}

// NewClient for database/sql
func NewClient(ctx context.Context, driverName string, poolCfg *PoolConfig,
	cfg *relational.SQLConnectConfig, retryCfg *retry.RetryConfig) (pool *sql.DB, err error) {
	pool, err = sql.Open(driverName, cfg.DSN())
	err = retry.Retry(func() error {
		pool, err = sql.Open(driverName, cfg.DSN())
		if err != nil {
			log.Printf("Failed to connect to sql(%s) database... Going to do the next attempt\n", driverName)
			return err
		}
		pool.SetMaxOpenConns(poolCfg.MaxOpenConns)
		pool.SetMaxIdleConns(poolCfg.MaxIdleConns)
		pool.SetConnMaxLifetime(poolCfg.MaxLifetime)
		pool.SetConnMaxIdleTime(poolCfg.ConnMaxIdleTime)

		ctx, cancel := context.WithTimeout(ctx, retryCfg.PingDelay)
		defer cancel()

		if err = pool.PingContext(ctx); err != nil {
			log.Printf("Failed to ping sql(%s) database... Going to do the next attempt\n", driverName)
			return err
		}
		return nil
	}, retryCfg)
	return pool, err
}
