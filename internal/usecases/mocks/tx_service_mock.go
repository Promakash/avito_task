// Code generated by mockery v2.52.1. DO NOT EDIT.

package mocks

import (
	domain "avito_shop/internal/domain"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// Transaction is an autogenerated mock type for the Transaction type
type Transaction struct {
	mock.Mock
}

// BuyItemByName provides a mock function with given fields: ctx, uid, name
func (_m *Transaction) BuyItemByName(ctx context.Context, uid int, name string) error {
	ret := _m.Called(ctx, uid, name)

	if len(ret) == 0 {
		panic("no return value specified for BuyItemByName")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, string) error); ok {
		r0 = rf(ctx, uid, name)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SendCoinByName provides a mock function with given fields: ctx, tx, to
func (_m *Transaction) SendCoinByName(ctx context.Context, tx domain.Transaction, to string) error {
	ret := _m.Called(ctx, tx, to)

	if len(ret) == 0 {
		panic("no return value specified for SendCoinByName")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.Transaction, string) error); ok {
		r0 = rf(ctx, tx, to)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewTransaction creates a new instance of Transaction. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTransaction(t interface {
	mock.TestingT
	Cleanup(func())
}) *Transaction {
	mock := &Transaction{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
