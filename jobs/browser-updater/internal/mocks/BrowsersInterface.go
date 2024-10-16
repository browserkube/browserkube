// Code generated by mockery v2.24.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "github.com/browserkube/browserkube/operator/api/v1"

	watch "k8s.io/apimachinery/pkg/watch"
)

// BrowsersInterface is an autogenerated mock type for the BrowsersInterface type
type BrowsersInterface struct {
	mock.Mock
}

// Create provides a mock function with given fields: _a0, _a1
func (_m *BrowsersInterface) Create(_a0 context.Context, _a1 *v1.Browser) (*v1.Browser, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *v1.Browser
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *v1.Browser) (*v1.Browser, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *v1.Browser) *v1.Browser); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.Browser)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *v1.Browser) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, name, opts
func (_m *BrowsersInterface) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	ret := _m.Called(ctx, name, opts)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, metav1.DeleteOptions) error); ok {
		r0 = rf(ctx, name, opts)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: ctx, name, options
func (_m *BrowsersInterface) Get(ctx context.Context, name string, options metav1.GetOptions) (*v1.Browser, error) {
	ret := _m.Called(ctx, name, options)

	var r0 *v1.Browser
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, metav1.GetOptions) (*v1.Browser, error)); ok {
		return rf(ctx, name, options)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, metav1.GetOptions) *v1.Browser); ok {
		r0 = rf(ctx, name, options)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.Browser)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, metav1.GetOptions) error); ok {
		r1 = rf(ctx, name, options)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: ctx, opts
func (_m *BrowsersInterface) List(ctx context.Context, opts metav1.ListOptions) (*v1.BrowserList, error) {
	ret := _m.Called(ctx, opts)

	var r0 *v1.BrowserList
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, metav1.ListOptions) (*v1.BrowserList, error)); ok {
		return rf(ctx, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, metav1.ListOptions) *v1.BrowserList); ok {
		r0 = rf(ctx, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.BrowserList)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, metav1.ListOptions) error); ok {
		r1 = rf(ctx, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Watch provides a mock function with given fields: ctx, pts
func (_m *BrowsersInterface) Watch(ctx context.Context, pts metav1.ListOptions) (watch.Interface, error) {
	ret := _m.Called(ctx, pts)

	var r0 watch.Interface
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, metav1.ListOptions) (watch.Interface, error)); ok {
		return rf(ctx, pts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, metav1.ListOptions) watch.Interface); ok {
		r0 = rf(ctx, pts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(watch.Interface)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, metav1.ListOptions) error); ok {
		r1 = rf(ctx, pts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WatchByName provides a mock function with given fields: ctx, name
func (_m *BrowsersInterface) WatchByName(ctx context.Context, name string) (watch.Interface, error) {
	ret := _m.Called(ctx, name)

	var r0 watch.Interface
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (watch.Interface, error)); ok {
		return rf(ctx, name)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) watch.Interface); ok {
		r0 = rf(ctx, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(watch.Interface)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewBrowsersInterface interface {
	mock.TestingT
	Cleanup(func())
}

// NewBrowsersInterface creates a new instance of BrowsersInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewBrowsersInterface(t mockConstructorTestingTNewBrowsersInterface) *BrowsersInterface {
	mock := &BrowsersInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
