package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"payment_processing_system/internal/adapters/client/kafka"
	pgxinternal "payment_processing_system/internal/adapters/client/relational/pgx"
	"payment_processing_system/internal/adapters/converter"
	v1 "payment_processing_system/internal/controller/http/v1"
	"payment_processing_system/internal/domain/service"
	"payment_processing_system/internal/domain/usecase"
	"payment_processing_system/pkg/db"
	"payment_processing_system/pkg/db/relational"
	"payment_processing_system/pkg/db/relational/pgx"
	"payment_processing_system/pkg/logger"
	"time"
)

func main() {
	port := 8000
	ctx := context.Background()
	cfg := relational.NewSQLConnectConfig("gorm", "gorm", "localhost", "5432", "public")
	retryCfg := db.NewRetryConfig(5, 3*time.Second)
	var err error
	var pool pgx.PgxIface
	pool, err = pgx.NewClient(ctx, cfg, retryCfg)
	if err != nil {
		panic(err)
	}

	l, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	log := logger.NewLogger(l)
	defer func() {
		_ = log.Sync()
	}()

	conv := converter.NewExchangeRatesAPI("", "", time.Minute)

	bStorage := pgxinternal.NewBalanceStorage(pool, log)
	tStorage := pgxinternal.NewTransactionStorage(pool, log)

	bService := service.NewBalanceService(bStorage)
	tService := service.NewTransactionService(tStorage)

	producer := kafka.NewApplyTransactionProducer()

	managerUseCase := usecase.NewManagerUseCase(bService, tService, producer)
	e := echo.New()
	s := v1.NewServer(managerUseCase, conv, log)
	v1.RegisterHandlers(e, s)
	e.Logger.Fatal(e.Start(fmt.Sprintf("0.0.0.0:%d", port)))
}
