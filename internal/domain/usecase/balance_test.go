package usecase

import (
	"context"
	"errors"
	"payment_processing_system/internal/domain"
	"payment_processing_system/internal/domain/entity"
	"payment_processing_system/internal/domain/usecase/mock"
	"testing"

	"github.com/stretchr/testify/suite"
)

type mockChainInfo struct {
	mockService string
	methodName  string
	args        []interface{}
	returnArgs  []interface{}
}

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

func (suite *BalanceUseCaseTestSuite) TestChangeAmount_Error() {
	var idNil *string
	var nilTransaction *entity.Transaction
	var validTransaction = entity.Transaction{
		ID: suite.transactionID,
	}
	exampleError := errors.New("example error")
	cancelError := errors.New("cancel error")
	testCases := []struct {
		ctx            context.Context
		id             *string
		amount         float32
		expectedErrors []error
		mockChain      []mockChainInfo
	}{
		{
			ctx:            context.Background(),
			id:             nil,
			amount:         1,
			expectedErrors: []error{domain.TransactionNilSourceOrDestinationErr},
		},
		{
			ctx:            context.Background(),
			id:             &suite.idFrom,
			amount:         0,
			expectedErrors: []error{domain.ChangeBalanceByZeroAmountErr},
		},
		{
			ctx:            context.Background(),
			id:             &suite.idTo,
			amount:         1,
			expectedErrors: []error{},
			mockChain: []mockChainInfo{
				{
					mockService: "ts",
					methodName:  "CreateDefaultTransaction",
					args: []interface{}{
						context.Background(),
						idNil,
						&suite.idTo,
						float32(1),
						entity.TypeOuterIncreasing,
					},
					returnArgs: []interface{}{
						nilTransaction,
						exampleError,
					},
				},
			},
		},
		{
			ctx:            context.Background(),
			id:             &suite.idTo,
			amount:         1,
			expectedErrors: []error{},
			mockChain: []mockChainInfo{
				{
					mockService: "ts",
					methodName:  "CreateDefaultTransaction",
					args: []interface{}{
						context.Background(),
						idNil,
						&suite.idTo,
						float32(1),
						entity.TypeOuterIncreasing,
					},
					returnArgs: []interface{}{
						&validTransaction,
						exampleError,
					},
				},
				{
					mockService: "ts",
					methodName:  "CancelByID",
					args: []interface{}{
						context.Background(),
						suite.transactionID,
					},
					returnArgs: []interface{}{
						cancelError,
					},
				},
			},
		},
	}
	for _, testCase := range testCases {
		bs := &mock.BalanceService{}
		ts := &mock.TransactionService{}
		useCase := &BalanceUseCase{bs: bs, ts: ts}
		for _, mockAction := range testCase.mockChain {
			if mockAction.mockService == "ts" {
				ts.On(mockAction.methodName, mockAction.args...).Return(mockAction.returnArgs...).Once()
				continue
			}
			if mockAction.mockService == "bs" {
				bs.On(mockAction.methodName, mockAction.args...).Return(mockAction.returnArgs...).Once()
				continue
			}
			suite.Fail("mockAction in mockChain is not necessary", mockAction.mockService, mockAction)
		}

		err := useCase.ChangeAmount(testCase.ctx, testCase.id, testCase.amount)
		for _, expectedErr := range testCase.expectedErrors {
			suite.ErrorIs(err, expectedErr)
		}
	}
}

func (suite *BalanceUseCaseTestSuite) TestChangeAmount_Success() {
	testCases := []struct {
		ctx             context.Context
		idFrom          *string
		idTo            *string
		amount          float32
		transaction     entity.Transaction
		transactionType entity.TransactionType
		expectedErr     error
	}{
		{
			ctx:             context.Background(),
			idFrom:          &suite.idTo,
			amount:          -1.3,
			transaction:     entity.Transaction{ID: suite.transactionID},
			transactionType: entity.TypeOuterDecreasing,
			expectedErr:     nil,
		},
		{
			ctx:             context.Background(),
			idTo:            &suite.idTo,
			amount:          2.3,
			transaction:     entity.Transaction{ID: suite.transactionID},
			transactionType: entity.TypeOuterIncreasing,
			expectedErr:     nil,
		},
	}
	for _, testCase := range testCases {
		bs := &mock.BalanceService{}
		ts := &mock.TransactionService{}
		useCase := &BalanceUseCase{bs: bs, ts: ts}

		bsChangeID := testCase.idFrom
		if bsChangeID == nil {
			bsChangeID = testCase.idTo
		}

		tsAmount := testCase.amount
		if tsAmount < 0 {
			tsAmount *= -1
		}

		ts.On("CreateDefaultTransaction", testCase.ctx, testCase.idFrom,
			testCase.idTo, tsAmount, testCase.transactionType).
			Return(&testCase.transaction, testCase.expectedErr).Once()
		ts.On("ProcessingByID",
			testCase.ctx, testCase.transaction.ID).Return(testCase.expectedErr).Once()
		bs.On("ChangeAmount",
			testCase.ctx, *bsChangeID, testCase.amount).Return(testCase.expectedErr).Once()
		ts.On("CompletedByID",
			testCase.ctx, testCase.transaction.ID).Return(testCase.expectedErr).Once()
		err := useCase.ChangeAmount(testCase.ctx, bsChangeID, testCase.amount)
		suite.Equal(testCase.expectedErr, err)
	}
}

func (suite *BalanceUseCaseTestSuite) TestTransfer_TransferNoError() {
	ctx := context.Background()
	var amount float32 = 1.3
	var expectedErr error
	var entityBalancePtr *entity.Balance
	var transaction = entity.Transaction{ID: suite.transactionID}
	suite.bs.On("GetByID", ctx, suite.idFrom).
		Return(entityBalancePtr, expectedErr).Once()
	suite.ts.On("CreateDefaultTransaction", ctx, &suite.idFrom,
		&suite.idTo, amount, entity.TypeTransfer).
		Return(&transaction, expectedErr).Once()
	suite.ts.On("ProcessingByID",
		ctx, transaction.ID).Return(expectedErr).Once()
	suite.bs.On("ChangeAmount",
		ctx, suite.idTo, amount).Return(expectedErr).Once()
	suite.bs.On("ChangeAmount",
		ctx, suite.idFrom, -amount).Return(expectedErr).Once()
	suite.ts.On("CompletedByID",
		ctx, transaction.ID).Return(expectedErr).Once()
	err := suite.useCase.Transfer(ctx, &suite.idFrom, &suite.idTo, amount)
	suite.Equal(expectedErr, err)
}

func TestBalanceUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(BalanceUseCaseTestSuite))
}
