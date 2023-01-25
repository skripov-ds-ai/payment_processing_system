package service

import (
	"context"
	"payment_processing_system/internal/domain/entity"
	"payment_processing_system/internal/domain/service/mock"
	"testing"

	"github.com/stretchr/testify/suite"
)

type BalanceServiceTestSuite struct {
	suite.Suite
	service *BalanceService
	storage *mock.BalanceStorage
}

func (suite *BalanceServiceTestSuite) SetupTest() {
	suite.storage = &mock.BalanceStorage{}
	suite.service = NewBalanceService(suite.storage)
}

func (suite *BalanceServiceTestSuite) TestGetByID_EqualIDNoError() {
	ctx := context.Background()
	id := "example"
	var expectedErr error
	var expectedBalance = entity.Balance{ID: id}
	suite.storage.On("GetByID", ctx, id).
		Return(&expectedBalance, expectedErr).
		Once()
	actual, err := suite.service.GetByID(ctx, id)
	suite.Equal(expectedErr, err)
	suite.Equal(expectedBalance.ID, actual.ID)
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(BalanceServiceTestSuite))
}
