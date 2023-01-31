package service

import (
	"context"
	"errors"
	"github.com/shopspring/decimal"
	"payment_processing_system/internal/domain"
	"payment_processing_system/internal/domain/entity"
	"payment_processing_system/internal/domain/service/mock"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/tkuchiki/faketime"
)

type TransactionServiceTestSuite struct {
	suite.Suite
	testService   *TransactionService
	testStorage   *mock.TransactionStorage
	ctx           context.Context
	id            int64
	transactionID uint64
}

func (suite *TransactionServiceTestSuite) SetupTest() {
	suite.testStorage = &mock.TransactionStorage{}
	suite.testService = NewTransactionService(suite.testStorage)

	suite.ctx = context.Background()
	suite.id = 1            //"example"
	suite.transactionID = 7 //"example-transaction"
}

func (suite *TransactionServiceTestSuite) TestCancelByID() {
	status := entity.StatusCancelled
	var expectedErr error
	suite.testStorage.On("UpdateStatusByID", suite.ctx, suite.id, status).
		Return(expectedErr).Once()
	err := suite.testService.CancelByID(suite.ctx, suite.id)
	suite.Equal(expectedErr, err)
}

func (suite *TransactionServiceTestSuite) TestProcessingByID() {
	status := entity.StatusProcessing
	var expectedErr error
	suite.testStorage.On("UpdateStatusByID", suite.ctx, suite.id, status).
		Return(expectedErr).Once()
	err := suite.testService.ProcessingByID(suite.ctx, suite.id)
	suite.Equal(expectedErr, err)
}

func (suite *TransactionServiceTestSuite) TestCompletedByID() {
	status := entity.StatusCompleted
	var expectedErr error
	suite.testStorage.On("UpdateStatusByID", suite.ctx, suite.id, status).
		Return(expectedErr).Once()
	err := suite.testService.CompletedByID(suite.ctx, suite.id)
	suite.Equal(expectedErr, err)
}

func (suite *TransactionServiceTestSuite) TestShouldRetryByID() {
	status := entity.StatusShouldRetry
	var expectedErr error
	suite.testStorage.On("UpdateStatusByID", suite.ctx, suite.id, status).
		Return(expectedErr).Once()
	err := suite.testService.ShouldRetryByID(suite.ctx, suite.id)
	suite.Equal(expectedErr, err)
}

func (suite *TransactionServiceTestSuite) TestCannotApplyByID() {
	status := entity.StatusCannotApply
	var expectedErr error
	suite.testStorage.On("UpdateStatusByID", suite.ctx, suite.id, status).
		Return(expectedErr).Once()
	err := suite.testService.CannotApplyByID(suite.ctx, suite.id)
	suite.Equal(expectedErr, err)
}

func (suite *TransactionServiceTestSuite) TestCreateDefaultTransaction() {
	f := faketime.NewFaketime(2021, time.March, 01, 01, 01, 01, 0, time.UTC)
	defer f.Undo()
	f.Do()

	var sourceID int64 = 2      //"a"
	var destinationID int64 = 3 //"b"
	testCases := []struct {
		ctx                   context.Context
		sourceID              *int64
		destinationID         *int64
		amount                decimal.Decimal
		ttype                 entity.TransactionType
		expectedTransaction   *entity.Transaction
		expectedErr           error
		expectMockStorageCall bool
	}{
		{
			ctx:                 context.Background(),
			amount:              decimal.Zero,
			expectedTransaction: nil,
			expectedErr:         domain.ZeroAmountTransactionErr,
		},
		{
			ctx:                 context.Background(),
			amount:              decimal.NewFromInt(-1),
			expectedTransaction: nil,
			expectedErr:         domain.NegativeAmountTransactionErr,
		},
		{
			ctx:                 context.Background(),
			amount:              decimal.NewFromInt(1),
			expectedTransaction: nil,
			expectedErr:         domain.TransactionNilSourceAndDestinationErr,
		},
		{
			ctx:                 context.Background(),
			amount:              decimal.NewFromInt(1),
			sourceID:            &sourceID,
			destinationID:       &sourceID,
			expectedTransaction: nil,
			expectedErr:         domain.TransactionSourceDestinationAreEqualErr,
		},
		{
			ctx:                   context.Background(),
			amount:                decimal.NewFromInt(1),
			sourceID:              &sourceID,
			destinationID:         &destinationID,
			ttype:                 entity.TypeTransfer,
			expectedTransaction:   &entity.Transaction{ID: suite.transactionID},
			expectedErr:           nil,
			expectMockStorageCall: true,
		},
		{
			ctx:                   context.Background(),
			amount:                decimal.NewFromInt(1),
			sourceID:              &sourceID,
			destinationID:         &destinationID,
			ttype:                 entity.TypeTransfer,
			expectedTransaction:   nil,
			expectedErr:           errors.New("database error"),
			expectMockStorageCall: true,
		},
	}
	for _, testCase := range testCases {
		storage := &mock.TransactionStorage{}
		service := NewTransactionService(storage)

		if testCase.expectMockStorageCall {
			transaction := entity.Transaction{
				SourceID:        testCase.sourceID,
				DestinationID:   testCase.destinationID,
				Amount:          testCase.amount,
				TType:           testCase.ttype,
				DateTimeCreated: time.Now(),
				DateTimeUpdated: time.Now(),
				Status:          entity.StatusCreated,
			}
			storage.On("Create", testCase.ctx, transaction).
				Return(testCase.expectedTransaction, testCase.expectedErr).
				Once()
		}
		id, err := service.CreateDefaultTransaction(
			testCase.ctx, testCase.sourceID, testCase.destinationID,
			testCase.amount, testCase.ttype)
		if err != nil {
			suite.EqualError(err, testCase.expectedErr.Error())
		}
		suite.Equal(testCase.expectedTransaction, id)
	}
}

func (suite *TransactionServiceTestSuite) TestGetByID() {
	testCases := []struct {
		ctx                 context.Context
		id                  int64
		expectedTransaction *entity.Transaction
		expectedErr         error
	}{
		{
			ctx: context.Background(),
			id:  4, //"example",
			expectedTransaction: &entity.Transaction{
				ID: 4, //"example",
			},
			expectedErr: nil,
		},
		{
			ctx:                 context.Background(),
			id:                  8, //"another-example",
			expectedTransaction: nil,
			expectedErr:         errors.New("example error"),
		},
	}
	for _, testCase := range testCases {
		storage := &mock.TransactionStorage{}
		service := NewTransactionService(storage)

		storage.On("GetByID", testCase.ctx, testCase.id).
			Return(testCase.expectedTransaction, testCase.expectedErr).
			Once()

		transaction, err := service.GetByID(testCase.ctx, testCase.id)
		if testCase.expectedErr != nil {
			suite.EqualError(err, testCase.expectedErr.Error())
			continue
		}
		if testCase.expectedTransaction != nil {
			suite.Equal(testCase.expectedTransaction.ID, transaction.ID)
			suite.Equal(*testCase.expectedTransaction, *transaction)
		}
	}
}

func TestTransactionServiceTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionServiceTestSuite))
}
