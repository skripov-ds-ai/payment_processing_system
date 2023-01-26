package usecase

import (
	"context"
	"payment_processing_system/internal/domain/entity"
	"payment_processing_system/internal/domain/usecase/mock"
	"testing"

	"github.com/stretchr/testify/suite"
)

type BalanceUseCaseTestSuite struct {
	suite.Suite
	bs            *mock.BalanceService
	ts            *mock.TransactionService
	useCase       *BalanceUseCase
	idFrom        string
	idTo          string
	transactionID string
}

func (suite *BalanceUseCaseTestSuite) SetupTest() {
	suite.bs = &mock.BalanceService{}
	suite.ts = &mock.TransactionService{}
	suite.useCase = &BalanceUseCase{bs: suite.bs, ts: suite.ts}
	suite.idFrom = "example-1"
	suite.idTo = "example-2"
	suite.transactionID = "transaction-1"
}

func (suite *BalanceUseCaseTestSuite) TestChangeAmount_OuterIncreasingNoError() {
	ctx := context.Background()
	var amount float32 = 1.3
	var expectedErr error
	var idFromPtr *string
	suite.ts.On("CreateDefaultTransaction", ctx, idFromPtr,
		&suite.idTo, amount, entity.TypeOuterIncreasing).
		Return(suite.transactionID, expectedErr).Once()
	suite.ts.On("ProcessingByID",
		ctx, suite.transactionID).Return(expectedErr).Once()
	suite.bs.On("ChangeAmount",
		ctx, suite.idTo, amount).Return(expectedErr).Once()
	suite.ts.On("CompletedByID",
		ctx, suite.transactionID).Return(expectedErr).Once()
	err := suite.useCase.ChangeAmount(ctx, &suite.idTo, amount)
	suite.Equal(expectedErr, err)
}
func (suite *BalanceUseCaseTestSuite) TestChangeAmount_OuterDecreasingNoError() {
	ctx := context.Background()
	var amount float32 = -10.3
	var expectedErr error
	var idToPtr *string
	suite.ts.On("CreateDefaultTransaction", ctx, &suite.idFrom,
		idToPtr, -amount, entity.TypeOuterDecreasing).
		Return(suite.transactionID, expectedErr).Once()
	suite.ts.On("ProcessingByID",
		ctx, suite.transactionID).Return(expectedErr).Once()
	suite.bs.On("ChangeAmount",
		ctx, suite.idFrom, amount).Return(expectedErr).Once()
	suite.ts.On("CompletedByID",
		ctx, suite.transactionID).Return(expectedErr).Once()
	err := suite.useCase.ChangeAmount(ctx, &suite.idFrom, amount)
	suite.Equal(expectedErr, err)
}

func (suite *BalanceUseCaseTestSuite) TestTransfer_TransferNoError() {
	ctx := context.Background()
	var amount float32 = 1.3
	var expectedErr error
	var entityBalancePtr *entity.Balance
	suite.bs.On("GetByID", ctx, suite.idFrom).
		Return(entityBalancePtr, expectedErr).Once()
	suite.ts.On("CreateDefaultTransaction", ctx, &suite.idFrom,
		&suite.idTo, amount, entity.TypeTransfer).
		Return(suite.transactionID, expectedErr).Once()
	suite.ts.On("ProcessingByID",
		ctx, suite.transactionID).Return(expectedErr).Once()
	suite.bs.On("ChangeAmount",
		ctx, suite.idTo, amount).Return(expectedErr).Once()
	suite.bs.On("ChangeAmount",
		ctx, suite.idFrom, -amount).Return(expectedErr).Once()
	suite.ts.On("CompletedByID",
		ctx, suite.transactionID).Return(expectedErr).Once()
	err := suite.useCase.Transfer(ctx, &suite.idFrom, &suite.idTo, amount)
	suite.Equal(expectedErr, err)
}

func TestBalanceUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(BalanceUseCaseTestSuite))
}
