// Package v1 provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.2 DO NOT EDIT.
package v1

// Defines values for FindBalancesParamsSort.
const (
	FindBalancesParamsSortDate FindBalancesParamsSort = "date"
	FindBalancesParamsSortId   FindBalancesParamsSort = "id"
	FindBalancesParamsSortSum  FindBalancesParamsSort = "sum"
)

// Defines values for GetBindedTransactionsParamsSort.
const (
	GetBindedTransactionsParamsSortDate GetBindedTransactionsParamsSort = "date"
	GetBindedTransactionsParamsSortId   GetBindedTransactionsParamsSort = "id"
	GetBindedTransactionsParamsSortSum  GetBindedTransactionsParamsSort = "sum"
)

// Balance defines model for Balance.
type Balance struct {
	Amount int64  `json:"amount"`
	Id     string `json:"id"`
}

// Error defines model for Error.
type Error struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

// NewBalance defines model for NewBalance.
type NewBalance struct {
	Amount int64  `json:"amount"`
	Id     string `json:"id"`
}

// FindBalancesParams defines parameters for FindBalances.
type FindBalancesParams struct {
	// Sort key to sort by - id, date, sum
	Sort *FindBalancesParamsSort `form:"sort,omitempty" json:"sort,omitempty"`

	// Limit maximum number of results to return
	Limit *int64 `form:"limit,omitempty" json:"limit,omitempty"`

	// Page page of collection
	Page *int64 `form:"page,omitempty" json:"page,omitempty"`
}

// FindBalancesParamsSort defines parameters for FindBalances.
type FindBalancesParamsSort string

// GetBalanceByIdParams defines parameters for GetBalanceByID.
type GetBalanceByIdParams struct {
	// Currency Currency to display balance
	Currency *string `form:"currency,omitempty" json:"currency,omitempty"`
}

// GetBindedTransactionsParams defines parameters for GetBindedTransactions.
type GetBindedTransactionsParams struct {
	// Sort key to sort by - id, date, sum
	Sort *GetBindedTransactionsParamsSort `form:"sort,omitempty" json:"sort,omitempty"`

	// Limit maximum number of results to return
	Limit *int64 `form:"limit,omitempty" json:"limit,omitempty"`

	// Page page of collection
	Page *int64 `form:"page,omitempty" json:"page,omitempty"`
}

// GetBindedTransactionsParamsSort defines parameters for GetBindedTransactions.
type GetBindedTransactionsParamsSort string

// AccrueOrWriteOffBalanceJSONRequestBody defines body for AccrueOrWriteOffBalance for application/json ContentType.
type AccrueOrWriteOffBalanceJSONRequestBody = NewBalance
