// Code generated by mockery v2.20.0. DO NOT EDIT.

package mocks

import (
	http "net/http"

	mock "github.com/stretchr/testify/mock"
)

// RestClient is an autogenerated mock type for the RestClient type
type RestClient struct {
	mock.Mock
}

// Get provides a mock function with given fields: url
func (_m *RestClient) Get(url string) (*http.Response, error) {
	ret := _m.Called(url)

	var r0 *http.Response
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*http.Response, error)); ok {
		return rf(url)
	}
	if rf, ok := ret.Get(0).(func(string) *http.Response); ok {
		r0 = rf(url)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*http.Response)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(url)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewRestClient interface {
	mock.TestingT
	Cleanup(func())
}

// NewRestClient creates a new instance of RestClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRestClient(t mockConstructorTestingTNewRestClient) *RestClient {
	mock := &RestClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
