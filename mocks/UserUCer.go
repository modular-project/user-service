// Code generated by mockery v2.13.1. DO NOT EDIT.

package mocks

import (
	model "users-service/model"

	mock "github.com/stretchr/testify/mock"
)

// UserUCer is an autogenerated mock type for the UserUCer type
type UserUCer struct {
	mock.Mock
}

// ChangePassword provides a mock function with given fields: _a0, _a1
func (_m *UserUCer) ChangePassword(_a0 uint, _a1 *string) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(uint, *string) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Data provides a mock function with given fields: _a0
func (_m *UserUCer) Data(_a0 uint) (model.User, error) {
	ret := _m.Called(_a0)

	var r0 model.User
	if rf, ok := ret.Get(0).(func(uint) model.User); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(model.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GenerateCode provides a mock function with given fields: _a0
func (_m *UserUCer) GenerateCode(_a0 uint) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(uint) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateData provides a mock function with given fields: _a0
func (_m *UserUCer) UpdateData(_a0 *model.User) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.User) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Verify provides a mock function with given fields: _a0, _a1
func (_m *UserUCer) Verify(_a0 uint, _a1 string) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(uint, string) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewUserUCer interface {
	mock.TestingT
	Cleanup(func())
}

// NewUserUCer creates a new instance of UserUCer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUserUCer(t mockConstructorTestingTNewUserUCer) *UserUCer {
	mock := &UserUCer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
