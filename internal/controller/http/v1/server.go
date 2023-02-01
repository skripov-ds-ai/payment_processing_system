package v1

import (
	"payment_processing_system/pkg/logger"

	"github.com/labstack/echo/v4"
)

type Server struct {
	manager *managerHandler
}

func NewServer(balanceUseCase ManagerUseCase, converter Converter, logger *logger.Logger) *Server {
	return &Server{manager: NewBalanceHandler(balanceUseCase, converter, logger)}
}

// (GET /balances)
func (s *Server) FindBalances(ctx echo.Context, params FindBalancesParams) error {
	return nil
}

// (POST /balances/{idFrom}/transfer/{idTo})
func (s *Server) TransferByIds(ctx echo.Context, idFrom, idTo int64) error {
	return s.manager.TransferByIds(ctx, idFrom, idTo)
}

// (GET /balances/{id})
func (s *Server) GetBalanceById(ctx echo.Context, id int64, params GetBalanceByIdParams) error {
	return s.manager.GetBalanceByID(ctx, id, params)
}

// (POST /balances/{id})
func (s *Server) AccrueOrWriteOffBalance(ctx echo.Context, id int64) error {
	return s.manager.AccrueOrWriteOffBalance(ctx, id)
}

// (GET /balances/{id}/transcations)
func (s *Server) GetBindedTransactions(ctx echo.Context, id int64, params GetBindedTransactionsParams) error {
	return s.manager.GetBindedTransactions(ctx, id, params)
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
