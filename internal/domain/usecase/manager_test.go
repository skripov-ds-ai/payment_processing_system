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

// TODO: rewrite

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
	useCase       *ManagerUseCase
	idFrom        string
	idTo          string
	transactionID string
}

func (suite *BalanceUseCaseTestSuite) SetupTest() {
	suite.bs = &mock.BalanceService{}
	suite.ts = &mock.TransactionService{}
	suite.useCase = NewBalanceUseCase(suite.bs, suite.ts)
	suite.idFrom = "example-1"
	suite.idTo = "example-2"
	suite.transactionID = "transaction-1"
}

func (suite *BalanceUseCaseTestSuite) TestGetByID() {
	testCases := []struct {
		ctx             context.Context
		id              string
		expectedBalance *entity.Balance
		expectedErr     error
	}{
		{
			ctx:             context.Background(),
			id:              suite.idTo,
			expectedBalance: &entity.Balance{ID: suite.idTo},
			expectedErr:     nil,
		},
		{
			ctx:             context.Background(),
			id:              suite.idFrom,
			expectedBalance: nil,
			expectedErr:     errors.New("example get error"),
		},
	}
	for _, testCase := range testCases {
		bs := &mock.BalanceService{}
		ts := &mock.TransactionService{}
		useCase := &ManagerUseCase{bs: bs, ts: ts}
		bs.On("GetByID", testCase.ctx, testCase.id).
			Return(testCase.expectedBalance, testCase.expectedErr).Once()
		balance, err := useCase.GetByID(testCase.ctx, testCase.id)
		if err != nil {
			suite.EqualError(err, testCase.expectedErr.Error())
		}
		suite.Equal(testCase.expectedBalance, balance)
	}
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
			expectedErrors: []error{exampleError},
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
			expectedErrors: []error{cancelError, exampleError},
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
		useCase := NewBalanceUseCase(bs, ts)
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
		useCase := NewBalanceUseCase(bs, ts)

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

func (suite *BalanceUseCaseTestSuite) TestTransfer_Error() {
	exampleError := errors.New("example error")
	testCases := []struct {
		ctx            context.Context
		idFrom         *string
		idTo           *string
		amount         float32
		expectedErrors []error
		mockChain      []mockChainInfo
	}{
		{
			ctx:            context.Background(),
			idFrom:         nil,
			amount:         1,
			expectedErrors: []error{domain.TransactionNilSourceErr},
		},
		{
			ctx:            context.Background(),
			idFrom:         &suite.idFrom,
			idTo:           nil,
			amount:         1,
			expectedErrors: []error{domain.TransactionNilDestinationErr},
		},
		{
			ctx:            context.Background(),
			idFrom:         &suite.idFrom,
			idTo:           &suite.idTo,
			amount:         0,
			expectedErrors: []error{domain.ChangeBalanceByZeroAmountErr},
		},
		{
			ctx:            context.Background(),
			idFrom:         &suite.idFrom,
			idTo:           &suite.idTo,
			amount:         -42,
			expectedErrors: []error{domain.NegativeAmountTransactionErr},
		},
		{
			ctx:            context.Background(),
			idFrom:         &suite.idFrom,
			idTo:           &suite.idFrom,
			amount:         100,
			expectedErrors: []error{domain.TransactionSourceDestinationAreEqualErr},
		},
		{
			ctx:            context.Background(),
			idFrom:         &suite.idFrom,
			idTo:           &suite.idTo,
			amount:         100,
			expectedErrors: []error{exampleError},
			mockChain: []mockChainInfo{
				{
					mockService: "bs",
					methodName:  "GetByID",
					args: []interface{}{
						context.Background(),
						suite.idFrom,
					},
					returnArgs: []interface{}{
						&entity.Balance{},
						exampleError,
					},
				},
			},
		},
	}
	for _, testCase := range testCases {
		bs := &mock.BalanceService{}
		ts := &mock.TransactionService{}
		useCase := NewBalanceUseCase(bs, ts)
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

		err := useCase.Transfer(testCase.ctx, testCase.idFrom, testCase.idTo, testCase.amount)
		for _, expectedErr := range testCase.expectedErrors {
			suite.ErrorIs(err, expectedErr)
		}
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
