package usecase

import (
	"context"
	"github.com/stretchr/testify/suite"
	"payment_processing_system/internal/domain"
	"payment_processing_system/internal/domain/entity"
	"payment_processing_system/internal/domain/usecase/mock"
	"testing"
)

type ApplierUseCaseTestSuite struct {
	suite.Suite
	bs       *mock.BalanceGetChangeService
	producer *mock.ConfirmTransactionProducer
	useCase  *ApplierUseCase
}

func (suite *ApplierUseCaseTestSuite) SetupTest() {
	suite.bs = &mock.BalanceGetChangeService{}
	suite.producer = &mock.ConfirmTransactionProducer{}
	suite.useCase = NewApplierUseCase(suite.bs, suite.producer)
}

func (suite *ApplierUseCaseTestSuite) TestApplyTransaction_UnknownTransactionTypeErr() {
	var expectedErr = domain.UnknownTransactionTypeErr
	ctx := context.Background()
	transaction := entity.Transaction{}
	err := suite.useCase.ApplyTransaction(ctx, transaction)
	suite.ErrorIs(err, expectedErr)
}

func (suite *ApplierUseCaseTestSuite) TestApplyTransaction_Success() {
	var errNil error
	sourceID, destinationID := int64(1), int64(2)
	transaction1 := entity.Transaction{
		ID:            0,
		DestinationID: &destinationID,
		TType:         entity.TypeOuterIncreasing,
	}
	transaction2 := entity.Transaction{
		ID:       0,
		SourceID: &sourceID,
		TType:    entity.TypeOuterDecreasing,
	}
	transaction3 := entity.Transaction{
		ID:       0,
		SourceID: &sourceID,
		TType:    entity.TypePayment,
	}
	testCases := []struct {
		ctx         context.Context
		transaction *entity.Transaction
		mockChain   []mockChainInfo
		expectedErr error
	}{
		{
			ctx:         context.Background(),
			transaction: &transaction1,
			mockChain: []mockChainInfo{
				{
					mockService: "bs",
					methodName:  "ChangeAmount",
					args: []interface{}{
						context.Background(),
						*transaction1.DestinationID,
						transaction1.Amount,
					},
					returnArgs: []interface{}{
						errNil,
					},
				},
				{
					mockService: "producer",
					methodName:  "CompleteTransaction",
					args: []interface{}{
						transaction1.ID,
					},
					returnArgs: []interface{}{
						errNil,
					},
				},
			},
			expectedErr: errNil,
		},
		{
			ctx:         context.Background(),
			transaction: &transaction2,
			mockChain: []mockChainInfo{
				{
					mockService: "bs",
					methodName:  "ChangeAmount",
					args: []interface{}{
						context.Background(),
						*transaction2.SourceID,
						transaction2.Amount.Neg(),
					},
					returnArgs: []interface{}{
						errNil,
					},
				},
				{
					mockService: "producer",
					methodName:  "CompleteTransaction",
					args: []interface{}{
						transaction2.ID,
					},
					returnArgs: []interface{}{
						errNil,
					},
				},
			},
			expectedErr: errNil,
		},
		{
			ctx:         context.Background(),
			transaction: &transaction3,
			mockChain: []mockChainInfo{
				{
					mockService: "bs",
					methodName:  "ChangeAmount",
					args: []interface{}{
						context.Background(),
						*transaction3.SourceID,
						transaction3.Amount.Neg(),
					},
					returnArgs: []interface{}{
						errNil,
					},
				},
				{
					mockService: "producer",
					methodName:  "CompleteTransaction",
					args: []interface{}{
						transaction3.ID,
					},
					returnArgs: []interface{}{
						errNil,
					},
				},
			},
			expectedErr: errNil,
		},
	}
	for _, testCase := range testCases {
		bs := &mock.BalanceGetChangeService{}
		producer := &mock.ConfirmTransactionProducer{}
		useCase := NewApplierUseCase(bs, producer)
		for _, mockAction := range testCase.mockChain {
			if mockAction.mockService == "bs" {
				bs.On(mockAction.methodName, mockAction.args...).Return(mockAction.returnArgs...).Once()
				continue
			}
			if mockAction.mockService == "producer" {
				producer.On(mockAction.methodName, mockAction.args...).Return(mockAction.returnArgs...).Once()
				continue
			}
			suite.Fail("mockAction in mockChain is not bs or ts", mockAction.mockService, mockAction)
		}

		err := useCase.ApplyTransaction(testCase.ctx, *testCase.transaction)
		suite.Equal(err, testCase.expectedErr)
	}
}

func TestApplierUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(ApplierUseCaseTestSuite))
}
