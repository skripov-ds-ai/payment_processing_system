package service

import (
	"context"
	"payment_processing_system/internal/domain"
	"payment_processing_system/internal/domain/entity"
	"payment_processing_system/internal/domain/service/mock"
	"testing"

	"github.com/stretchr/testify/suite"
)

type BalanceServiceTestSuite struct {
	suite.Suite
	service *BalanceService
	storage *mock.BalanceStorage
	id      int64
}

func (suite *BalanceServiceTestSuite) SetupTest() {
	suite.storage = &mock.BalanceStorage{}
	suite.service = NewBalanceService(suite.storage)
	suite.id = 42
}

func (suite *BalanceServiceTestSuite) TestGetByID_EqualIDNoError() {
	ctx := context.Background()
	var expectedErr error
	var expectedBalance = entity.Balance{ID: suite.id}
	suite.storage.On("GetByID", ctx, suite.id).
		Return(&expectedBalance, expectedErr).
		Once()
	actual, err := suite.service.GetByID(ctx, suite.id)
	suite.Equal(expectedErr, err)
	suite.Equal(expectedBalance.ID, actual.ID)
	suite.Equal(expectedBalance, *actual)
}

func (suite *BalanceServiceTestSuite) TestChangeAmount_ByZeroErr() {
	ctx := context.Background()
	var amount float32 = 0
	err := suite.service.ChangeAmount(ctx, suite.id, amount)
	suite.ErrorIsf(err, domain.ChangeBalanceByZeroAmountErr, "id = %q ; amount = %f ; %w", suite.id, amount, domain.ChangeBalanceByZeroAmountErr)
}

func (suite *BalanceServiceTestSuite) TestChangeAmount_IncreaseAmount() {
	ctx := context.Background()
	var amount float32 = 1
	var expectedError error
	suite.storage.On("IncreaseAmount", ctx, suite.id, amount).
		Return(expectedError).
		Once()
	err := suite.service.ChangeAmount(ctx, suite.id, amount)
	suite.Equal(expectedError, err)
}

func (suite *BalanceServiceTestSuite) TestChangeAmount_DecreaseAmount() {
	ctx := context.Background()
	var amount float32 = -1
	var expectedError error
	suite.storage.On("DecreaseAmount", ctx, suite.id, -amount).
		Return(expectedError).
		Once()
	err := suite.service.ChangeAmount(ctx, suite.id, amount)
	suite.Equal(expectedError, err)
}

func TestBalanceServiceTestSuite(t *testing.T) {
	suite.Run(t, new(BalanceServiceTestSuite))
}
