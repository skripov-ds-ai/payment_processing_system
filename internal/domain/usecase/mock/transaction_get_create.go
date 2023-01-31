// Code generated by mockery v2.16.0. DO NOT EDIT.

package mock

import (
	context "context"
	entity "payment_processing_system/internal/domain/entity"

	decimal "github.com/shopspring/decimal"

	mock "github.com/stretchr/testify/mock"
)

// TransactionGetCreateService is an autogenerated mock type for the TransactionGetCreateService type
type TransactionGetCreateService struct {
	mock.Mock
}

// CancelByID provides a mock function with given fields: ctx, id
func (_m *TransactionGetCreateService) CancelByID(ctx context.Context, id uint64) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateDefaultTransaction provides a mock function with given fields: ctx, sourceID, destinationID, amount, ttype
func (_m *TransactionGetCreateService) CreateDefaultTransaction(ctx context.Context, sourceID *int64, destinationID *int64, amount decimal.Decimal, ttype entity.TransactionType) (*entity.Transaction, error) {
	ret := _m.Called(ctx, sourceID, destinationID, amount, ttype)

	var r0 *entity.Transaction
	if rf, ok := ret.Get(0).(func(context.Context, *int64, *int64, decimal.Decimal, entity.TransactionType) *entity.Transaction); ok {
		r0 = rf(ctx, sourceID, destinationID, amount, ttype)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Transaction)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *int64, *int64, decimal.Decimal, entity.TransactionType) error); ok {
		r1 = rf(ctx, sourceID, destinationID, amount, ttype)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByID provides a mock function with given fields: ctx, id
func (_m *TransactionGetCreateService) GetByID(ctx context.Context, id uint64) (*entity.Transaction, error) {
	ret := _m.Called(ctx, id)

	var r0 *entity.Transaction
	if rf, ok := ret.Get(0).(func(context.Context, uint64) *entity.Transaction); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Transaction)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewTransactionGetCreateService interface {
	mock.TestingT
	Cleanup(func())
}

// NewTransactionGetCreateService creates a new instance of TransactionGetCreateService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTransactionGetCreateService(t mockConstructorTestingTNewTransactionGetCreateService) *TransactionGetCreateService {
	mock := &TransactionGetCreateService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}