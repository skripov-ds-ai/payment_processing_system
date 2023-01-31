package v1

import (
	"github.com/labstack/echo/v4"
	"payment_processing_system/pkg/logger"
)

type Server struct {
	balance *balanceHandler
}

func NewServer(balanceUseCase ManagerUseCase, converter Converter, logger *logger.Logger) *Server {
	return &Server{balance: NewBalanceHandler(balanceUseCase, converter, logger)}
}

// (GET /balances)
func (s *Server) FindBalances(ctx echo.Context, params FindBalancesParams) error {
	return nil
}

// (POST /balances/{idFrom}/transfer/{idTo})
func (s *Server) TransferByIds(ctx echo.Context, idFrom, idTo int64) error {
	return nil
}

// (GET /balances/{id})
func (s *Server) GetBalanceById(ctx echo.Context, id int64, params GetBalanceByIdParams) error {
	return s.balance.GetBalanceByID(ctx, id, params)
}

// (POST /balances/{id})
func (s *Server) AccrueOrWriteOffBalance(ctx echo.Context, id int64) error {
	return nil
}

// (GET /balances/{id}/transcations)
func (s *Server) GetBindedTransactions(ctx echo.Context, id int64, params GetBindedTransactionsParams) error {
	return nil
}

// (POST /reservation/balances/{id})
func (s *Server) ReserveOnSeparateAccount(ctx echo.Context, id int64) error {
	return nil
}

// (POST /reservation/balances/{id}/cancel)
func (s *Server) CancelReservation(ctx echo.Context, id int64) error {
	return nil
}

// (POST /reservation/balances/{id}/confirm)
func (s *Server) ConfirmReservation(ctx echo.Context, id int64) error {
	return nil
}
