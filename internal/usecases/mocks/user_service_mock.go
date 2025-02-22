// Code generated by mockery v2.52.1. DO NOT EDIT.

package mocks

import (
	domain "avito_shop/internal/domain"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// User is an autogenerated mock type for the User type
type User struct {
	mock.Mock
}

// GetByName provides a mock function with given fields: ctx, name
func (_m *User) GetByName(ctx context.Context, name string) (domain.User, error) {
	ret := _m.Called(ctx, name)

	if len(ret) == 0 {
		panic("no return value specified for GetByName")
	}

	var r0 domain.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (domain.User, error)); ok {
		return rf(ctx, name)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) domain.User); ok {
		r0 = rf(ctx, name)
	} else {
		r0 = ret.Get(0).(domain.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetInfoByID provides a mock function with given fields: ctx, id
func (_m *User) GetInfoByID(ctx context.Context, id int) (domain.UserInfo, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetInfoByID")
	}

	var r0 domain.UserInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) (domain.UserInfo, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) domain.UserInfo); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(domain.UserInfo)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Put provides a mock function with given fields: ctx, user
func (_m *User) Put(ctx context.Context, user domain.User) (int, error) {
	ret := _m.Called(ctx, user)

	if len(ret) == 0 {
		panic("no return value specified for Put")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.User) (int, error)); ok {
		return rf(ctx, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, domain.User) int); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, domain.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewUser creates a new instance of User. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUser(t interface {
	mock.TestingT
	Cleanup(func())
}) *User {
	mock := &User{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
