package service

import (
	"context"
	"payment_processing_system/internal/domain/entity"
	"payment_processing_system/internal/domain/service/mock"
	"testing"

	"github.com/stretchr/testify/suite"
)

type TransactionServiceTestSuite struct {
	suite.Suite
	service *TransactionService
	storage *mock.TransactionStorage
	id      string
}

func (suite *TransactionServiceTestSuite) SetupTest() {
	suite.storage = &mock.TransactionStorage{}
	suite.service = NewTransactionService(suite.storage)
	suite.id = "example"
}

func (suite *TransactionServiceTestSuite) TestGetByID_EqualIDNoError() {
	ctx := context.Background()
	var expectedErr error
	var expectedTransaction = entity.Transaction{ID: suite.id}
	suite.storage.On("GetByID", ctx, suite.id).
		Return(&expectedTransaction, expectedErr).
		Once()
	actual, err := suite.service.GetByID(ctx, suite.id)
	suite.Equal(expectedErr, err)
	suite.Equal(expectedTransaction.ID, actual.ID)
	suite.Equal(expectedTransaction, *actual)
}

func TestTransactionServiceTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionServiceTestSuite))
}
