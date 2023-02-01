package usecase

import (
	"context"
	"errors"
	"payment_processing_system/internal/domain"
	"payment_processing_system/internal/domain/entity"
	"payment_processing_system/internal/domain/usecase/mock"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
)

type ManagerUseCaseTestSuite struct {
	suite.Suite
	bs            *mock.BalanceGetService
	ts            *mock.TransactionGetCreateService
	producer      *mock.ApplyTransactionProducer
	useCase       *ManagerUseCase
	idFrom        int64
	idTo          int64
	transactionID uint64
}

func (suite *ManagerUseCaseTestSuite) SetupTest() {
	suite.bs = &mock.BalanceGetService{}
	suite.ts = &mock.TransactionGetCreateService{}
	suite.producer = &mock.ApplyTransactionProducer{}
	suite.useCase = NewManagerUseCase(suite.bs, suite.ts, suite.producer)
	suite.idFrom = 1
	suite.idTo = 2
	suite.transactionID = 42
}

func (suite *ManagerUseCaseTestSuite) TestPayForService_Error() {
	var idx = int64(1)
	testCases := []struct {
		ctx         context.Context
		id          *int64
		amount      decimal.Decimal
		expectedErr error
	}{
		{
			ctx:         context.Background(),
			id:          nil,
			amount:      decimal.NewFromInt(1),
			expectedErr: domain.TransactionNilSourceErr,
		},
		{
			ctx:         context.Background(),
			id:          &idx,
			amount:      decimal.Zero,
			expectedErr: domain.ChangeBalanceByZeroAmountErr,
		},
		{
			ctx:         context.Background(),
			id:          &idx,
			amount:      decimal.NewFromInt(-1),
			expectedErr: domain.NegativeAmountTransactionErr,
		},
	}
	for _, testCase := range testCases {
		bs := &mock.BalanceGetService{}
		ts := &mock.TransactionGetCreateService{}
		producer := &mock.ApplyTransactionProducer{}
		useCase := NewManagerUseCase(bs, ts, producer)
		_, err := useCase.PayForService(testCase.ctx, testCase.id, testCase.amount)
		suite.ErrorIs(err, testCase.expectedErr)
	}
}

func (suite *ManagerUseCaseTestSuite) TestGetByID() {
	testCases := []struct {
		ctx             context.Context
		id              int64
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
		bs := &mock.BalanceGetService{}
		ts := &mock.TransactionGetCreateService{}
		producer := &mock.ApplyTransactionProducer{}
		useCase := NewManagerUseCase(bs, ts, producer)
		bs.On("GetByID", testCase.ctx, testCase.id).
			Return(testCase.expectedBalance, testCase.expectedErr).Once()
		balance, err := useCase.GetBalanceByID(testCase.ctx, testCase.id)
		if err != nil {
			suite.EqualError(err, testCase.expectedErr.Error())
		}
		suite.Equal(testCase.expectedBalance, balance)
	}
}

func (suite *ManagerUseCaseTestSuite) TestChangeAmount_Error() {
	var idNil *int64
	var nilTransaction *entity.Transaction
	var validTransaction = entity.Transaction{
		ID: suite.transactionID,
	}
	exampleError := errors.New("example error")
	cancelError := errors.New("cancel error")
	testCases := []struct {
		ctx            context.Context
		id             *int64
		amount         decimal.Decimal
		expectedErrors []error
		mockChain      []mockChainInfo
	}{
		{
			ctx:            context.Background(),
			id:             nil,
			amount:         decimal.NewFromInt(1),
			expectedErrors: []error{domain.TransactionNilSourceOrDestinationErr},
		},
		{
			ctx:            context.Background(),
			id:             &suite.idFrom,
			amount:         decimal.Zero,
			expectedErrors: []error{domain.ChangeBalanceByZeroAmountErr},
		},
		{
			ctx:            context.Background(),
			id:             &suite.idTo,
			amount:         decimal.NewFromInt(1),
			expectedErrors: []error{exampleError},
			mockChain: []mockChainInfo{
				{
					mockService: "ts",
					methodName:  "CreateDefaultTransaction",
					args: []interface{}{
						context.Background(),
						idNil,
						&suite.idTo,
						decimal.NewFromInt(1),
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
			amount:         decimal.NewFromInt(1),
			expectedErrors: []error{cancelError, exampleError},
			mockChain: []mockChainInfo{
				{
					mockService: "ts",
					methodName:  "CreateDefaultTransaction",
					args: []interface{}{
						context.Background(),
						idNil,
						&suite.idTo,
						decimal.NewFromInt(1),
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
		bs := &mock.BalanceGetService{}
		ts := &mock.TransactionGetCreateService{}
		producer := &mock.ApplyTransactionProducer{}
		useCase := NewManagerUseCase(bs, ts, producer)
		for _, mockAction := range testCase.mockChain {
			if mockAction.mockService == "ts" {
				ts.On(mockAction.methodName, mockAction.args...).Return(mockAction.returnArgs...).Once()
				continue
			}
			if mockAction.mockService == "bs" {
				bs.On(mockAction.methodName, mockAction.args...).Return(mockAction.returnArgs...).Once()
				continue
			}
			suite.Fail("mockAction in mockChain is not bs or ts", mockAction.mockService, mockAction)
		}

		_, err := useCase.ChangeAmount(testCase.ctx, testCase.id, testCase.amount)
		for _, expectedErr := range testCase.expectedErrors {
			suite.ErrorIs(err, expectedErr)
		}
	}
}

func (suite *ManagerUseCaseTestSuite) TestChangeAmount_Success() {
	testCases := []struct {
		ctx             context.Context
		idFrom          *int64
		idTo            *int64
		amount          decimal.Decimal
		transaction     entity.Transaction
		transactionType entity.TransactionType
		expectedErr     error
	}{
		{
			ctx:             context.Background(),
			idFrom:          &suite.idTo,
			amount:          decimal.NewFromInt(-1),
			transaction:     entity.Transaction{ID: suite.transactionID},
			transactionType: entity.TypeOuterDecreasing,
			expectedErr:     nil,
		},
		{
			ctx:             context.Background(),
			idTo:            &suite.idTo,
			amount:          decimal.NewFromInt(2),
			transaction:     entity.Transaction{ID: suite.transactionID},
			transactionType: entity.TypeOuterIncreasing,
			expectedErr:     nil,
		},
	}
	for _, testCase := range testCases {
		bs := &mock.BalanceGetService{}
		ts := &mock.TransactionGetCreateService{}
		producer := &mock.ApplyTransactionProducer{}
		useCase := NewManagerUseCase(bs, ts, producer)

		bsChangeID := testCase.idFrom
		if bsChangeID == nil {
			bsChangeID = testCase.idTo
		}

		tsAmount := testCase.amount.Copy()
		if tsAmount.IsNegative() {
			tsAmount = tsAmount.Neg()
		}

		ts.On("CreateDefaultTransaction", testCase.ctx, testCase.idFrom,
			testCase.idTo, tsAmount, testCase.transactionType).
			Return(&testCase.transaction, testCase.expectedErr).Once()
		producer.On("ApplyTransaction", testCase.ctx, testCase.transaction).
			Return(testCase.expectedErr).Once()
		transaction, err := useCase.ChangeAmount(testCase.ctx, bsChangeID, testCase.amount)
		suite.Equal(testCase.expectedErr, err)
		suite.Equal(testCase.transaction, *transaction)
	}
}

func (suite *ManagerUseCaseTestSuite) TestTransfer_Error() {
	exampleError := errors.New("example error")
	testCases := []struct {
		ctx            context.Context
		idFrom         *int64
		idTo           *int64
		amount         decimal.Decimal
		expectedErrors []error
		mockChain      []mockChainInfo
	}{
		{
			ctx:            context.Background(),
			idFrom:         nil,
			amount:         decimal.NewFromInt(1),
			expectedErrors: []error{domain.TransactionNilSourceErr},
		},
		{
			ctx:            context.Background(),
			idFrom:         &suite.idFrom,
			idTo:           nil,
			amount:         decimal.NewFromInt(1),
			expectedErrors: []error{domain.TransactionNilDestinationErr},
		},
		{
			ctx:            context.Background(),
			idFrom:         &suite.idFrom,
			idTo:           &suite.idTo,
			amount:         decimal.Zero,
			expectedErrors: []error{domain.ChangeBalanceByZeroAmountErr},
		},
		{
			ctx:            context.Background(),
			idFrom:         &suite.idFrom,
			idTo:           &suite.idTo,
			amount:         decimal.NewFromInt(-42),
			expectedErrors: []error{domain.NegativeAmountTransactionErr},
		},
		{
			ctx:            context.Background(),
			idFrom:         &suite.idFrom,
			idTo:           &suite.idFrom,
			amount:         decimal.NewFromInt(100),
			expectedErrors: []error{domain.TransactionSourceDestinationAreEqualErr},
		},
		{
			ctx:            context.Background(),
			idFrom:         &suite.idFrom,
			idTo:           &suite.idTo,
			amount:         decimal.NewFromInt(100),
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
		bs := &mock.BalanceGetService{}
		ts := &mock.TransactionGetCreateService{}
		producer := &mock.ApplyTransactionProducer{}
		useCase := NewManagerUseCase(bs, ts, producer)
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

		_, err := useCase.Transfer(testCase.ctx, testCase.idFrom, testCase.idTo, testCase.amount)
		for _, expectedErr := range testCase.expectedErrors {
			suite.ErrorIs(err, expectedErr)
		}
	}
}

func (suite *ManagerUseCaseTestSuite) TestTransfer_TransferNoError() {
	ctx := context.Background()
	amount := decimal.NewFromInt(7)
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
		ctx, suite.idFrom, amount.Neg()).Return(expectedErr).Once()
	suite.ts.On("CompletedByID",
		ctx, transaction.ID).Return(expectedErr).Once()
	suite.producer.On("ApplyTransaction", ctx, transaction).
		Return(expectedErr).Once()
	_, err := suite.useCase.Transfer(ctx, &suite.idFrom, &suite.idTo, amount)
	suite.Equal(expectedErr, err)
}

func TestBalanceUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(ManagerUseCaseTestSuite))
}
