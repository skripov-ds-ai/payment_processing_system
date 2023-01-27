package v1

import (
	"context"
	"fmt"
	"net/http"
	"payment_processing_system/internal/domain/entity"
	"payment_processing_system/pkg/logger"

	"go.uber.org/zap"

	"github.com/labstack/echo/v4"
)

// Converter of currency
type Converter interface {
	ConvertFromRUBToCurrency(amount float32, currency string) (float32, error)
}

// BalanceService is standard service for balance
type BalanceService interface {
	GetByID(ctx context.Context, id string) (*entity.Balance, error)
	ChangeAmount(ctx context.Context, id string, amount float32) error
}

type balanceHandler struct {
	service   BalanceService
	converter Converter
	logger    *logger.Logger
}

func NewBalanceHandler(service BalanceService, converter Converter, logger *logger.Logger) *balanceHandler {
	return &balanceHandler{
		service:   service,
		converter: converter,
		logger:    logger,
	}
}

// GetBalanceByID returns json of balance object or error
// (GET /balances/{id})
func (b *balanceHandler) GetBalanceByID(ctx echo.Context, id string, params GetBalanceByIdParams) error {
	balance, err := b.service.GetByID(ctx.Request().Context(), id)
	if err != nil {
		b.logger.Error("error during getting balance", zap.String("id", id), zap.Error(err))
		e := Error{
			Message: fmt.Sprintf("something went wrong during getting balance by id = %s", id),
		}
		err1 := ctx.JSON(http.StatusBadRequest, e)
		if err1 != nil {
			b.logger.Error("error during sending error json", zap.Error(err1))
		}
		return err
	}
	if balance == nil {
		e := Error{Message: "balance not found"}
		err1 := ctx.JSON(http.StatusNotFound, e)
		if err != nil {
			b.logger.Error("error during sending error json", zap.Error(err1))
		}
		return nil
	}
	// convert
	if params.Currency != nil && *params.Currency != "RUB" {
		newAmount, err1 := b.converter.ConvertFromRUBToCurrency(balance.Amount, *params.Currency)
		if err1 != nil {
			b.logger.Error("error during balance convert",
				zap.Float32("amount", balance.Amount),
				zap.String("currency", *params.Currency),
				zap.Error(err1))
			e := Error{
				Message: fmt.Sprintf("something went wrong during convertation to %s", *params.Currency),
			}
			err2 := ctx.JSON(http.StatusNotFound, e)
			if err2 != nil {
				b.logger.Error("error during sending error json", zap.Error(err2))
			}
			return err
		}
		balance.Amount = newAmount
	}
	_ = ctx.JSON(http.StatusOK, *balance)
	return nil
}
