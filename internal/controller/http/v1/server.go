package v1

import "github.com/labstack/echo/v4"

type Server struct {
	balance *balanceHandler
}

// (GET /balances)
func (s *Server) FindBalances(ctx echo.Context, params FindBalancesParams) error {
	return nil
}

// (POST /balances/{idFrom}/transfer/{idTo})
func (s *Server) TransferByIds(ctx echo.Context, idFrom string, idTo string) error {
	return nil
}

// (GET /balances/{id})
func (s *Server) GetBalanceById(ctx echo.Context, id string, params GetBalanceByIdParams) error {
	return s.balance.GetBalanceByID(ctx, id, params)
}

// (POST /balances/{id})
func (s *Server) AccrueOrWriteOffBalance(ctx echo.Context, id string) error {
	return nil
}

// (GET /balances/{id}/transcations)
func (s *Server) GetBindedTransactions(ctx echo.Context, id string, params GetBindedTransactionsParams) error {
	return nil
}

// (POST /reservation/balances/{id})
func (s *Server) ReserveOnSeparateAccount(ctx echo.Context, id string) error {
	return nil
}

// (POST /reservation/balances/{id}/cancel)
func (s *Server) CancelReservation(ctx echo.Context, id string) error {
	return nil
}

// (POST /reservation/balances/{id}/confirm)
func (s *Server) ConfirmReservation(ctx echo.Context, id string) error {
	return nil
}
