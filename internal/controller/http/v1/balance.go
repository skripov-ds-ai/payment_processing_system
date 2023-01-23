package v1

import (
	"context"
	"fmt"
	"net/http"
	"payment_processing_system/internal/domain/entity"

	"github.com/labstack/echo/v4"
)

// BalanceService is standard service for balance
type BalanceService interface {
	GetByID(ctx context.Context, id string) (*entity.Balance, error)
	ChangeAmount(ctx context.Context, id string, amount int64) error
}

type balanceHandler struct {
	service BalanceService
}

// GetBalanceByID returns json of balance object or error
// (GET /balances/{id})
func (b *balanceHandler) GetBalanceByID(ctx echo.Context, id string, params GetBalanceByIdParams) error {
	// TODO: add currency processing
	balance, err := b.service.GetByID(ctx.Request().Context(), id)
	if err != nil {
		// TODO: add logging
		e := Error{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("something went wrong during getting balance by id = %s", id),
		}
		// TODO: wrap an error
		_ = ctx.JSON(http.StatusBadRequest, e)
		return err
	}
	if balance == nil {
		e := Error{Code: http.StatusNotFound, Message: "balance not found"}
		// TODO: wrap an error
		_ = ctx.JSON(http.StatusNotFound, e)
		return nil
	}
	// TODO: wrap an error
	_ = ctx.JSON(http.StatusOK, *balance)
	return nil
}
