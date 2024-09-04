// Code generated by mockery v2.43.1. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/ARUMANDESU/go-revise/internal/domain"
	mock "github.com/stretchr/testify/mock"
)

// UserProvider is an autogenerated mock type for the UserProvider type
type UserProvider struct {
	mock.Mock
}

// GetUser provides a mock function with given fields: ctx, id
func (_m *UserProvider) GetUser(ctx context.Context, id string) (domain.User, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetUser")
	}

	var r0 domain.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (domain.User, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) domain.User); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(domain.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserByChatID provides a mock function with given fields: ctx, chatID
func (_m *UserProvider) GetUserByChatID(ctx context.Context, chatID int64) (domain.User, error) {
	ret := _m.Called(ctx, chatID)

	if len(ret) == 0 {
		panic("no return value specified for GetUserByChatID")
	}

	var r0 domain.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (domain.User, error)); ok {
		return rf(ctx, chatID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) domain.User); ok {
		r0 = rf(ctx, chatID)
	} else {
		r0 = ret.Get(0).(domain.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, chatID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewUserProvider creates a new instance of UserProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserProvider(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserProvider {
	mock := &UserProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}