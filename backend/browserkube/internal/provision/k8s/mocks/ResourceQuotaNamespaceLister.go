// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	labels "k8s.io/apimachinery/pkg/labels"

	v1 "k8s.io/api/core/v1"
)

// ResourceQuotaNamespaceLister is an autogenerated mock type for the ResourceQuotaNamespaceLister type
type ResourceQuotaNamespaceLister struct {
	mock.Mock
}

// Get provides a mock function with given fields: name
func (_m *ResourceQuotaNamespaceLister) Get(name string) (*v1.ResourceQuota, error) {
	ret := _m.Called(name)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *v1.ResourceQuota
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*v1.ResourceQuota, error)); ok {
		return rf(name)
	}
	if rf, ok := ret.Get(0).(func(string) *v1.ResourceQuota); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.ResourceQuota)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: selector
func (_m *ResourceQuotaNamespaceLister) List(selector labels.Selector) ([]*v1.ResourceQuota, error) {
	ret := _m.Called(selector)

	if len(ret) == 0 {
		panic("no return value specified for List")
	}

	var r0 []*v1.ResourceQuota
	var r1 error
	if rf, ok := ret.Get(0).(func(labels.Selector) ([]*v1.ResourceQuota, error)); ok {
		return rf(selector)
	}
	if rf, ok := ret.Get(0).(func(labels.Selector) []*v1.ResourceQuota); ok {
		r0 = rf(selector)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*v1.ResourceQuota)
		}
	}

	if rf, ok := ret.Get(1).(func(labels.Selector) error); ok {
		r1 = rf(selector)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewResourceQuotaNamespaceLister creates a new instance of ResourceQuotaNamespaceLister. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewResourceQuotaNamespaceLister(t interface {
	mock.TestingT
	Cleanup(func())
},
) *ResourceQuotaNamespaceLister {
	mock := &ResourceQuotaNamespaceLister{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
