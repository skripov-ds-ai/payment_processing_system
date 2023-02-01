package v1

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"payment_processing_system/internal/domain/entity"
	"payment_processing_system/pkg/logger"

	"github.com/shopspring/decimal"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// Converter of currency
type Converter interface {
	ConvertFromRUBToCurrency(amount decimal.Decimal, currency string) (decimal.Decimal, error)
}

// BalanceService is standard useCase for manager
type ManagerUseCase interface {
	GetBalanceByID(ctx context.Context, id int64) (*entity.Balance, error)
	GetBalanceTransactions(ctx context.Context, balanceID int64, limit, offset uint64, orderBy string) ([]*entity.Transaction, error)
	ChangeAmount(ctx context.Context, id *int64, amount decimal.Decimal) (transaction *entity.Transaction, err error)
	// ChangeAmount(ctx context.Context, id *int64, amount decimal.Decimal) error
	// PayForService(ctx context.Context, id *int64, amount decimal.Decimal) error
	// Transfer(ctx context.Context, idFrom, idTo *int64, amount decimal.Decimal) error
}

type managerHandler struct {
	useCase   ManagerUseCase
	converter Converter
	logger    *logger.Logger
}

func NewBalanceHandler(useCase ManagerUseCase, converter Converter, logger *logger.Logger) *managerHandler {
	return &managerHandler{
		useCase:   useCase,
		converter: converter,
		logger:    logger,
	}
}

// TODO
// (POST /balances/{id})
func (m *managerHandler) AccrueOrWriteOffBalance(ctx echo.Context, id int64) error {
	var balanceChangeBody NewBalance
	err := ctx.Bind(&balanceChangeBody)
	if err != nil {
		e := Error{
			Message: fmt.Sprintf("something went wrong during read body ; id = %d ; %v", id, err),
		}
		_ = ctx.JSON(http.StatusBadRequest, e)
	}
	amount, err := decimal.NewFromString(balanceChangeBody.Amount)
	if err != nil {
		e := Error{
			Message: fmt.Sprintf("something went wrong during amount converting ; id = %d ; %v", id, err),
		}
		_ = ctx.JSON(http.StatusBadRequest, e)
		return err
	}
	transaction, err := m.useCase.ChangeAmount(ctx.Request().Context(), &id, amount)
	if err != nil {
		e := Error{
			Message: fmt.Sprintf("something went wrong during balance changing ; id = %d ; %v", id, err),
		}
		_ = ctx.JSON(http.StatusInternalServerError, e)
		return err
	}
	if transaction == nil {
		e := Error{
			Message: fmt.Sprintf("transaction was not created during balance changing ; id = %d ; %v", id),
		}
		_ = ctx.JSON(http.StatusInternalServerError, e)
		// TODO
		return errors.New("transaction was not created")
	}
	_ = ctx.JSON(http.StatusOK, *transaction)
	return nil
}

// (GET /balances/{id}/transcations)
func (m *managerHandler) GetBindedTransactions(ctx echo.Context, id int64, params GetBindedTransactionsParams) error {
	// TODO: add validation by validator
	var limit, offset uint64
	limit = 10
	var orderBy = "id"
	if params.Limit != nil && *params.Limit > 0 {
		limit = *params.Limit
	}
	if params.Page != nil {
		offset = *params.Page
	}
	if params.Sort != nil {
		orderBy = string(*params.Sort)
	}
	transactions, err := m.useCase.GetBalanceTransactions(ctx.Request().Context(), id, limit, offset, orderBy)
	if err != nil {
		m.logger.Error("error during getting manager transactions", zap.Int64("id", id), zap.Error(err))
		e := Error{
			Message: fmt.Sprintf("something went wrong during getting manager transactions by id = %d ; %v", id, err),
		}
		err1 := ctx.JSON(http.StatusBadRequest, e)
		if err1 != nil {
			m.logger.Error("error during sending error json", zap.Error(err1))
		}
		return err
	}
	_ = ctx.JSON(http.StatusOK, transactions)
	return nil
}

// GetBalanceByID returns json of manager object or error
// (GET /balances/{id})
func (m *managerHandler) GetBalanceByID(ctx echo.Context, id int64, params GetBalanceByIdParams) error {
	balance, err := m.useCase.GetBalanceByID(ctx.Request().Context(), id)
	if err != nil {
		m.logger.Error("error during getting manager", zap.Int64("id", id), zap.Error(err))
		e := Error{
			Message: fmt.Sprintf("something went wrong during getting manager by id = %d ; %v", id, err),
		}
		err1 := ctx.JSON(http.StatusBadRequest, e)
		if err1 != nil {
			m.logger.Error("error during sending error json", zap.Error(err1))
		}
		return err
	}
	if balance == nil {
		e := Error{Message: "manager not found"}
		err1 := ctx.JSON(http.StatusNotFound, e)
		if err1 != nil {
			m.logger.Error("error during sending error json", zap.Error(err1))
		}
		return nil
	}
	// convert
	if params.Currency != nil && *params.Currency != "RUB" {
		newAmount, err1 := m.converter.ConvertFromRUBToCurrency(balance.Amount, *params.Currency)
		if err1 != nil {
			m.logger.Error("error during manager convert",
				zap.Int64("id", id),
				zap.String("amount", balance.Amount.String()),
				zap.String("currency", *params.Currency),
				zap.Error(err1))
			e := Error{
				Message: fmt.Sprintf("something went wrong during convertation to %s ; %v", *params.Currency, err1),
			}
			err2 := ctx.JSON(http.StatusNotFound, e)
			if err2 != nil {
				m.logger.Error("error during sending error json", zap.Error(err2))
			}
			return err
		}
		balance.Amount = newAmount
	}
	_ = ctx.JSON(http.StatusOK, *balance)
	return nil
}
