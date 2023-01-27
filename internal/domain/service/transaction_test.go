package service

import (
	"context"
	"errors"
	"github.com/tkuchiki/faketime"
	"payment_processing_system/internal/domain"
	"payment_processing_system/internal/domain/entity"
	"payment_processing_system/internal/domain/service/mock"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type TransactionServiceTestSuite struct {
	suite.Suite
	testService *TransactionService
	testStorage *mock.TransactionStorage
}

func (suite *TransactionServiceTestSuite) SetupTest() {
	suite.testStorage = &mock.TransactionStorage{}
	suite.testService = NewTransactionService(suite.testStorage)
}

func (suite *TransactionServiceTestSuite) TestCreateDefaultTransaction() {
	f := faketime.NewFaketime(2021, time.March, 01, 01, 01, 01, 0, time.UTC)
	defer f.Undo()
	f.Do()

	sourceID := "a"
	destinationID := "b"
	testCases := []struct {
		ctx                   context.Context
		sourceID              *string
		destinationID         *string
		amount                float32
		ttype                 entity.TransactionType
		expectedID            string
		expectedErr           error
		expectMockStorageCall bool
	}{
		{
			ctx:         context.Background(),
			amount:      0,
			expectedID:  "",
			expectedErr: domain.ZeroAmountTransactionErr,
		},
		{
			ctx:         context.Background(),
			amount:      -1,
			expectedID:  "",
			expectedErr: domain.NegativeAmountTransactionErr,
		},
		{
			ctx:         context.Background(),
			amount:      1,
			expectedID:  "",
			expectedErr: domain.TransactionNilSourceAndDestinationErr,
		},
		{
			ctx:           context.Background(),
			amount:        1,
			sourceID:      &sourceID,
			destinationID: &sourceID,
			expectedID:    "",
			expectedErr:   domain.TransactionSourceDestinationAreEqualErr,
		},
		{
			ctx:                   context.Background(),
			amount:                1,
			sourceID:              &sourceID,
			destinationID:         &destinationID,
			ttype:                 entity.TypeTransfer,
			expectedID:            "42",
			expectedErr:           nil,
			expectMockStorageCall: true,
		},
		{
			ctx:                   context.Background(),
			amount:                1,
			sourceID:              &sourceID,
			destinationID:         &destinationID,
			ttype:                 entity.TypeTransfer,
			expectedID:            "",
			expectedErr:           errors.New("database error"),
			expectMockStorageCall: true,
		},
	}
	for _, testCase := range testCases {
		storage := &mock.TransactionStorage{}
		service := NewTransactionService(storage)

		if testCase.expectMockStorageCall {
			transaction := entity.Transaction{
				SourceID:      testCase.sourceID,
				DestinationID: testCase.destinationID,
				Amount:        testCase.amount,
				Type:          testCase.ttype,
				DateTime:      time.Now(),
				Status:        entity.StatusCreated,
			}
			storage.On("Create", testCase.ctx, transaction).
				Return(testCase.expectedID, testCase.expectedErr).
				Once()
		}
		id, err := service.CreateDefaultTransaction(
			testCase.ctx, testCase.sourceID, testCase.destinationID,
			testCase.amount, testCase.ttype)
		if err != nil {
			suite.EqualError(err, testCase.expectedErr.Error())
		}
		suite.Equal(testCase.expectedID, id)
	}
}

func (suite *TransactionServiceTestSuite) TestGetByID() {
	testCases := []struct {
		ctx                 context.Context
		id                  string
		expectedTransaction *entity.Transaction
		expectedErr         error
	}{
		{
			ctx: context.Background(),
			id:  "example",
			expectedTransaction: &entity.Transaction{
				ID: "example",
			},
			expectedErr: nil,
		},
		{
			ctx:                 context.Background(),
			id:                  "another-example",
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
