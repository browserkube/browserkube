// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	resource "k8s.io/apimachinery/pkg/api/resource"

	session "github.com/browserkube/browserkube/pkg/session"
)

// SessionWatchInterface is an autogenerated mock type for the SessionWatchInterface type
type SessionWatchInterface struct {
	mock.Mock
}

// Exists provides a mock function with given fields: sessionID
func (_m *SessionWatchInterface) Exists(sessionID string) (bool, error) {
	ret := _m.Called(sessionID)

	if len(ret) == 0 {
		panic("no return value specified for Exists")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (bool, error)); ok {
		return rf(sessionID)
	}
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(sessionID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(sessionID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetQuotas provides a mock function with given fields:
func (_m *SessionWatchInterface) GetQuotas() (resource.Quantity, resource.Quantity) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetQuotas")
	}

	var r0 resource.Quantity
	var r1 resource.Quantity
	if rf, ok := ret.Get(0).(func() (resource.Quantity, resource.Quantity)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() resource.Quantity); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(resource.Quantity)
	}

	if rf, ok := ret.Get(1).(func() resource.Quantity); ok {
		r1 = rf()
	} else {
		r1 = ret.Get(1).(resource.Quantity)
	}

	return r0, r1
}

// GetSessions provides a mock function with given fields:
func (_m *SessionWatchInterface) GetSessions() []string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetSessions")
	}

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// LoadByID provides a mock function with given fields: sessionID
func (_m *SessionWatchInterface) LoadByID(sessionID string) (*session.Session, error) {
	ret := _m.Called(sessionID)

	if len(ret) == 0 {
		panic("no return value specified for LoadByID")
	}

	var r0 *session.Session
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*session.Session, error)); ok {
		return rf(sessionID)
	}
	if rf, ok := ret.Get(0).(func(string) *session.Session); ok {
		r0 = rf(sessionID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*session.Session)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(sessionID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Watch provides a mock function with given fields: ctx
func (_m *SessionWatchInterface) Watch(ctx context.Context) <-chan *session.Session {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Watch")
	}

	var r0 <-chan *session.Session
	if rf, ok := ret.Get(0).(func(context.Context) <-chan *session.Session); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan *session.Session)
		}
	}

	return r0
}

// NewSessionWatchInterface creates a new instance of SessionWatchInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSessionWatchInterface(t interface {
	mock.TestingT
	Cleanup(func())
},
) *SessionWatchInterface {
	mock := &SessionWatchInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
