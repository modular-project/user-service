// Code generated by mockery v2.13.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// IsForbiddener is an autogenerated mock type for the IsForbiddener type
type IsForbiddener struct {
	mock.Mock
}

// IsForbidden provides a mock function with given fields:
func (_m *IsForbiddener) IsForbidden() {
	_m.Called()
}

type mockConstructorTestingTNewIsForbiddener interface {
	mock.TestingT
	Cleanup(func())
}

// NewIsForbiddener creates a new instance of IsForbiddener. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewIsForbiddener(t mockConstructorTestingTNewIsForbiddener) *IsForbiddener {
	mock := &IsForbiddener{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
