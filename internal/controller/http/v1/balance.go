package v1

import (
	"context"
	"fmt"
	"net/http"
	"payment_processing_system/internal/domain/entity"
	"payment_processing_system/pkg/logger"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// Converter of currency
type Converter interface {
	ConvertFromRUBToCurrency(amount int64, currency string) (int64, error)
}

// BalanceService is standard useCase for balance
type BalanceUseCase interface {
	GetByID(ctx context.Context, id int64) (*entity.Balance, error)
	ChangeAmount(ctx context.Context, id *int64, amount int64) error
	PayForService(ctx context.Context, id *int64, amount int64) error
	Transfer(ctx context.Context, idFrom, idTo *int64, amount int64) error
}

type balanceHandler struct {
	useCase   BalanceUseCase
	converter Converter
	logger    *logger.Logger
}

func NewBalanceHandler(useCase BalanceUseCase, converter Converter, logger *logger.Logger) *balanceHandler {
	return &balanceHandler{
		useCase:   useCase,
		converter: converter,
		logger:    logger,
	}
}

// GetBalanceByID returns json of balance object or error
// (GET /balances/{id})
func (b *balanceHandler) GetBalanceByID(ctx echo.Context, id int64, params GetBalanceByIdParams) error {
	balance, err := b.useCase.GetByID(ctx.Request().Context(), id)
	if err != nil {
		b.logger.Error("error during getting balance", zap.Int64("id", id), zap.Error(err))
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
		if err1 != nil {
			b.logger.Error("error during sending error json", zap.Error(err1))
		}
		return nil
	}
	// convert
	if params.Currency != nil && *params.Currency != "RUB" {
		newAmount, err1 := b.converter.ConvertFromRUBToCurrency(balance.Amount, *params.Currency)
		if err1 != nil {
			b.logger.Error("error during balance convert",
				zap.Int64("id", id),
				zap.Int64("amount", balance.Amount),
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
