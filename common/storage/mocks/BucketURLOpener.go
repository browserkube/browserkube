// Code generated by mockery v2.24.0. DO NOT EDIT.

package mocks

import (
	context "context"

	blob "gocloud.dev/blob"

	mock "github.com/stretchr/testify/mock"

	url "net/url"
)

// BucketURLOpener is an autogenerated mock type for the BucketURLOpener type
type BucketURLOpener struct {
	mock.Mock
}

// OpenBucketURL provides a mock function with given fields: ctx, u
func (_m *BucketURLOpener) OpenBucketURL(ctx context.Context, u *url.URL) (*blob.Bucket, error) {
	ret := _m.Called(ctx, u)

	var r0 *blob.Bucket
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *url.URL) (*blob.Bucket, error)); ok {
		return rf(ctx, u)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *url.URL) *blob.Bucket); ok {
		r0 = rf(ctx, u)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*blob.Bucket)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *url.URL) error); ok {
		r1 = rf(ctx, u)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewBucketURLOpener interface {
	mock.TestingT
	Cleanup(func())
}

// NewBucketURLOpener creates a new instance of BucketURLOpener. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewBucketURLOpener(t mockConstructorTestingTNewBucketURLOpener) *BucketURLOpener {
	mock := &BucketURLOpener{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
