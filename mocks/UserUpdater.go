// Code generated by mockery v2.13.1. DO NOT EDIT.

package mocks

import (
	model "users-service/model"

	mock "github.com/stretchr/testify/mock"
)

// UserUpdater is an autogenerated mock type for the UserUpdater type
type UserUpdater struct {
	mock.Mock
}

// Update provides a mock function with given fields: _a0, _a1
func (_m *UserUpdater) Update(_a0 uint, _a1 *model.User) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(uint, *model.User) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewUserUpdater interface {
	mock.TestingT
	Cleanup(func())
}

// NewUserUpdater creates a new instance of UserUpdater. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUserUpdater(t mockConstructorTestingTNewUserUpdater) *UserUpdater {
	mock := &UserUpdater{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
