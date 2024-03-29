// Code generated by mockery v2.13.1. DO NOT EDIT.

package mocks

import (
	model "users-service/model"

	mock "github.com/stretchr/testify/mock"
)

// EMPLStorager is an autogenerated mock type for the EMPLStorager type
type EMPLStorager struct {
	mock.Mock
}

// Fire provides a mock function with given fields: _a0
func (_m *EMPLStorager) Fire(_a0 uint) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(uint) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: from, target
func (_m *EMPLStorager) Get(from *model.UserRole, target uint) (model.UserJobs, error) {
	ret := _m.Called(from, target)

	var r0 model.UserJobs
	if rf, ok := ret.Get(0).(func(*model.UserRole, uint) model.UserJobs); ok {
		r0 = rf(from, target)
	} else {
		r0 = ret.Get(0).(model.UserJobs)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.UserRole, uint) error); ok {
		r1 = rf(from, target)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Hire provides a mock function with given fields: _a0
func (_m *EMPLStorager) Hire(_a0 *model.UserRole) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.UserRole) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Search provides a mock function with given fields: _a0
func (_m *EMPLStorager) Search(_a0 *model.SearchEMPL) ([]model.User, error) {
	ret := _m.Called(_a0)

	var r0 []model.User
	if rf, ok := ret.Get(0).(func(*model.SearchEMPL) []model.User); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.SearchEMPL) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SearchWaiters provides a mock function with given fields: _a0, _a1
func (_m *EMPLStorager) SearchWaiters(_a0 uint, _a1 *model.Search) ([]model.User, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []model.User
	if rf, ok := ret.Get(0).(func(uint, *model.Search) []model.User); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint, *model.Search) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Self provides a mock function with given fields: _a0
func (_m *EMPLStorager) Self(_a0 uint) (model.UserJobs, error) {
	ret := _m.Called(_a0)

	var r0 model.UserJobs
	if rf, ok := ret.Get(0).(func(uint) model.UserJobs); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(model.UserJobs)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewEMPLStorager interface {
	mock.TestingT
	Cleanup(func())
}

// NewEMPLStorager creates a new instance of EMPLStorager. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewEMPLStorager(t mockConstructorTestingTNewEMPLStorager) *EMPLStorager {
	mock := &EMPLStorager{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
