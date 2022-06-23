// Code generated by mockery v2.13.1. DO NOT EDIT.

package mocks

import (
	model "users-service/model"

	mock "github.com/stretchr/testify/mock"
)

// SignUCer is an autogenerated mock type for the SignUCer type
type SignUCer struct {
	mock.Mock
}

// Refresh provides a mock function with given fields: _a0
func (_m *SignUCer) Refresh(_a0 *string) (string, error) {
	ret := _m.Called(_a0)

	var r0 string
	if rf, ok := ret.Get(0).(func(*string) string); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SignIn provides a mock function with given fields: _a0
func (_m *SignUCer) SignIn(_a0 *model.LogIn) (string, string, error) {
	ret := _m.Called(_a0)

	var r0 string
	if rf, ok := ret.Get(0).(func(*model.LogIn) string); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 string
	if rf, ok := ret.Get(1).(func(*model.LogIn) string); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Get(1).(string)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(*model.LogIn) error); ok {
		r2 = rf(_a0)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// SignOut provides a mock function with given fields: refresh
func (_m *SignUCer) SignOut(refresh *string) error {
	ret := _m.Called(refresh)

	var r0 error
	if rf, ok := ret.Get(0).(func(*string) error); ok {
		r0 = rf(refresh)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SignUp provides a mock function with given fields: _a0
func (_m *SignUCer) SignUp(_a0 *model.LogIn) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.LogIn) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewSignUCer interface {
	mock.TestingT
	Cleanup(func())
}

// NewSignUCer creates a new instance of SignUCer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewSignUCer(t mockConstructorTestingTNewSignUCer) *SignUCer {
	mock := &SignUCer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
