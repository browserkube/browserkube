// Code generated by mockery v2.24.0. DO NOT EDIT.

package mocks

import (
	context "context"

	driver "gocloud.dev/blob/driver"
	gcerr "gocloud.dev/gcerrors"

	mock "github.com/stretchr/testify/mock"
)

// Bucket is an autogenerated mock type for the Bucket type
type Bucket struct {
	mock.Mock
}

// As provides a mock function with given fields: i
func (_m *Bucket) As(i interface{}) bool {
	ret := _m.Called(i)

	var r0 bool
	if rf, ok := ret.Get(0).(func(interface{}) bool); ok {
		r0 = rf(i)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Attributes provides a mock function with given fields: ctx, key
func (_m *Bucket) Attributes(ctx context.Context, key string) (*driver.Attributes, error) {
	ret := _m.Called(ctx, key)

	var r0 *driver.Attributes
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*driver.Attributes, error)); ok {
		return rf(ctx, key)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *driver.Attributes); ok {
		r0 = rf(ctx, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*driver.Attributes)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Close provides a mock function with given fields:
func (_m *Bucket) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Copy provides a mock function with given fields: ctx, dstKey, srcKey, opts
func (_m *Bucket) Copy(ctx context.Context, dstKey, srcKey string, opts *driver.CopyOptions) error {
	ret := _m.Called(ctx, dstKey, srcKey, opts)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, *driver.CopyOptions) error); ok {
		r0 = rf(ctx, dstKey, srcKey, opts)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: ctx, key
func (_m *Bucket) Delete(ctx context.Context, key string) error {
	ret := _m.Called(ctx, key)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ErrorAs provides a mock function with given fields: _a0, _a1
func (_m *Bucket) ErrorAs(_a0 error, _a1 interface{}) bool {
	ret := _m.Called(_a0, _a1)

	var r0 bool
	if rf, ok := ret.Get(0).(func(error, interface{}) bool); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// ErrorCode provides a mock function with given fields: _a0
func (_m *Bucket) ErrorCode(_a0 error) gcerr.ErrorCode {
	ret := _m.Called(_a0)

	var r0 gcerr.ErrorCode
	if rf, ok := ret.Get(0).(func(error) gcerr.ErrorCode); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(gcerr.ErrorCode)
	}

	return r0
}

// ListPaged provides a mock function with given fields: ctx, opts
func (_m *Bucket) ListPaged(ctx context.Context, opts *driver.ListOptions) (*driver.ListPage, error) {
	ret := _m.Called(ctx, opts)

	var r0 *driver.ListPage
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *driver.ListOptions) (*driver.ListPage, error)); ok {
		return rf(ctx, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *driver.ListOptions) *driver.ListPage); ok {
		r0 = rf(ctx, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*driver.ListPage)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *driver.ListOptions) error); ok {
		r1 = rf(ctx, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewRangeReader provides a mock function with given fields: ctx, key, offset, length, opts
func (_m *Bucket) NewRangeReader(ctx context.Context, key string, offset, length int64, opts *driver.ReaderOptions) (driver.Reader, error) {
	ret := _m.Called(ctx, key, offset, length, opts)

	var r0 driver.Reader
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, int64, int64, *driver.ReaderOptions) (driver.Reader, error)); ok {
		return rf(ctx, key, offset, length, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, int64, int64, *driver.ReaderOptions) driver.Reader); ok {
		r0 = rf(ctx, key, offset, length, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(driver.Reader)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, int64, int64, *driver.ReaderOptions) error); ok {
		r1 = rf(ctx, key, offset, length, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewTypedWriter provides a mock function with given fields: ctx, key, contentType, opts
func (_m *Bucket) NewTypedWriter(ctx context.Context, key, contentType string, opts *driver.WriterOptions) (driver.Writer, error) {
	ret := _m.Called(ctx, key, contentType, opts)

	var r0 driver.Writer
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, *driver.WriterOptions) (driver.Writer, error)); ok {
		return rf(ctx, key, contentType, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, *driver.WriterOptions) driver.Writer); ok {
		r0 = rf(ctx, key, contentType, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(driver.Writer)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, *driver.WriterOptions) error); ok {
		r1 = rf(ctx, key, contentType, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SignedURL provides a mock function with given fields: ctx, key, opts
func (_m *Bucket) SignedURL(ctx context.Context, key string, opts *driver.SignedURLOptions) (string, error) {
	ret := _m.Called(ctx, key, opts)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *driver.SignedURLOptions) (string, error)); ok {
		return rf(ctx, key, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, *driver.SignedURLOptions) string); ok {
		r0 = rf(ctx, key, opts)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, *driver.SignedURLOptions) error); ok {
		r1 = rf(ctx, key, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewBucket interface {
	mock.TestingT
	Cleanup(func())
}

// NewBucket creates a new instance of Bucket. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewBucket(t mockConstructorTestingTNewBucket) *Bucket {
	mock := &Bucket{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
