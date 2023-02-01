package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"payment_processing_system/internal/adapters/client/kafka"
	pgxinternal "payment_processing_system/internal/adapters/client/relational/pgx"
	"payment_processing_system/internal/adapters/converter"
	v1 "payment_processing_system/internal/controller/http/v1"
	"payment_processing_system/internal/domain/service"
	"payment_processing_system/internal/domain/usecase"
	"payment_processing_system/pkg/db"
	"payment_processing_system/pkg/db/relational"
	"payment_processing_system/pkg/db/relational/pgx"
	"payment_processing_system/pkg/kafka/pubsub"
	"payment_processing_system/pkg/logger"
	"syscall"
	"time"

	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func main() {
	port := 8000
	ctx := context.Background()
	cfg := relational.NewSQLConnectConfig("gorm", "gorm", "postgres", "5432", "gorm")
	retryCfg := db.NewRetryConfig(5, 3*time.Second)
	var err error
	var pool pgx.PgxIface
	pool, err = pgx.NewClient(ctx, cfg, retryCfg)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

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

	// TODO
	address := make([]string, 0)
	kafkaPubSub := pubsub.NewProducer(address, log)
	producer := kafka.NewApplyTransactionProducer("apply", kafkaPubSub)

	managerUseCase := usecase.NewManagerUseCase(bService, tService, producer)

	swagger, err := v1.GetSwagger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading swagger spec\n: %s", err)
		return
	}
	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swagger.Servers = nil

	e := echo.New()
	// Log all requests
	e.Use(echomiddleware.Logger())
	// Use our validation middleware to check all requests against the
	// OpenAPI schema.
	e.Use(middleware.OapiRequestValidator(swagger))

	s := v1.NewServer(managerUseCase, conv, log)
	v1.RegisterHandlers(e, s)

	// start server in goroutine
	go func() {
		if err := e.Start(fmt.Sprintf(":%d", port)); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server...")
		}
	}()

	// graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)
	<-sigs

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
