package storage

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"testing"
	"time"

	"github.com/browserkube/browserkube/storage/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gocloud.dev/blob"
	"gocloud.dev/blob/driver"
	gcerr "gocloud.dev/gcerrors"
)

//go:generate mockery --name Bucket --replace-type gocloud.dev/internal/gcerr=gocloud.dev/gcerrors --dir $GOPATH/pkg/mod/gocloud.dev@v0.36.0/blob/driver --output mocks
//go:generate mockery --name BucketURLOpener --dir $GOPATH/pkg/mod/gocloud.dev@v0.36.0/blob --output mocks
func Test_blobStorage_GetSessionRecord(t *testing.T) {
	type args struct {
		sessionID string
		filename  string
		scheme    string
		urlstr    string
	}
	tests := []struct {
		name    string
		args    args
		prepare func(f *mocks.Bucket)
		wantErr bool
	}{
		{
			name: "GetSessionRecord: success",
			args: args{
				sessionID: "1",
				filename:  "filename",
				scheme:    "test1",
				urlstr:    "test1://",
			},
			prepare: func(f *mocks.Bucket) {
				f.On("NewRangeReader", context.Background(), mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&r, nil).Maybe()
			},
			wantErr: false,
		},
		{
			name: "GetSessionRecord: returns error, error expected",
			args: args{
				sessionID: "1",
				filename:  "filename",
				scheme:    "test2",
				urlstr:    "test2://",
			},
			prepare: func(f *mocks.Bucket) {
				f.On("ErrorCode", mock.Anything).Return(gcerr.Unknown).Maybe()
				f.On("NewRangeReader", context.Background(), mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("error")).Maybe()
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blobStorage, fakeDriver := mockNew(t, tt.args.scheme, tt.args.urlstr)

			if tt.prepare != nil {
				tt.prepare(fakeDriver)
			}

			sessionRecord, err := blobStorage.GetFile(context.Background(), tt.args.sessionID, tt.args.filename)
			if !tt.wantErr {
				require.NoError(t, err)
				require.NotEmpty(t, sessionRecord)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func Test_blobStorage_SaveSessionRecord(t *testing.T) {
	type args struct {
		sessionID string
		prefix    string
		sr        *BlobFile
		scheme    string
		urlstr    string
	}
	tests := []struct {
		name    string
		args    args
		prepare func(f *mocks.Bucket)
		wantErr bool
	}{
		{
			name: "SaveSessionRecord: success",
			args: args{
				sessionID: "1",
				sr:        sessionRecord,
				scheme:    "test3",
				urlstr:    "test3://",
			},
			prepare: func(f *mocks.Bucket) {
				f.On("NewTypedWriter", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&w, nil).Maybe()
			},
			wantErr: false,
		},
		{
			name: "SaveSessionRecord: returns error, error expected",
			args: args{
				sessionID: "1",
				sr:        sessionRecord,
				scheme:    "test4",
				urlstr:    "test4://",
			},
			prepare: func(f *mocks.Bucket) {
				f.On("ErrorCode", mock.Anything).Return(gcerr.Unknown).Maybe()
				f.On("NewTypedWriter", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("error")).Maybe()
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blobStorage, fakeDriver := mockNew(t, tt.args.scheme, tt.args.urlstr)

			if tt.prepare != nil {
				tt.prepare(fakeDriver)
			}

			err := blobStorage.SaveFile(context.Background(), tt.args.sessionID, tt.args.prefix, tt.args.sr)
			if !tt.wantErr {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func Test_blobStorage_DeleteSessionRecord(t *testing.T) {
	type args struct {
		sessionID string
		filename  string
		scheme    string
		urlstr    string
	}
	tests := []struct {
		name    string
		args    args
		prepare func(f *mocks.Bucket)
		wantErr bool
	}{
		{
			name: "DeleteSessionRecord: success",
			args: args{
				sessionID: "1",
				filename:  "filename",
				scheme:    "test5",
				urlstr:    "test5://",
			},
			prepare: func(f *mocks.Bucket) {
				f.On("Delete", mock.Anything, mock.Anything).Return(nil).Maybe()
			},
			wantErr: false,
		},
		{
			name: "DeleteSessionRecord: returns error, error expected",
			args: args{
				sessionID: "1",
				filename:  "filename",
				scheme:    "test6",
				urlstr:    "test6://",
			},
			prepare: func(f *mocks.Bucket) {
				f.On("ErrorCode", mock.Anything).Return(gcerr.Unknown).Maybe()
				f.On("Delete", mock.Anything, mock.Anything).Return(errors.New("error")).Maybe()
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blobStorage, fakeDriver := mockNew(t, tt.args.scheme, tt.args.urlstr)

			if tt.prepare != nil {
				tt.prepare(fakeDriver)
			}

			err := blobStorage.DeleteFile(context.Background(), tt.args.sessionID, tt.args.filename)
			if !tt.wantErr {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func Test_blobStorage_Close(t *testing.T) {
	type args struct {
		scheme string
		urlstr string
	}
	tests := []struct {
		name    string
		args    args
		prepare func(f *mocks.Bucket)
		wantErr bool
	}{
		{
			name: "Close: success",
			args: args{
				scheme: "test7",
				urlstr: "test7://",
			},
			prepare: func(f *mocks.Bucket) {
				f.On("Close").Return(nil).Maybe()
			},
			wantErr: false,
		},
		{
			name: "Close: returns error, error expected",
			args: args{
				scheme: "test8",
				urlstr: "test8://",
			},
			prepare: func(f *mocks.Bucket) {
				f.On("ErrorCode", mock.Anything).Return(gcerr.Unknown).Maybe()
				f.On("Close").Return(errors.New("error")).Maybe()
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs, fakeDriver := mockNew(t, tt.args.scheme, tt.args.urlstr)

			if tt.prepare != nil {
				tt.prepare(fakeDriver)
			}

			err := bs.Close()
			if !tt.wantErr {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func Test_blobStorage_ListSessionRecords(t *testing.T) {
	type args struct {
		sessionID string
		prefix    string
		scheme    string
		urlstr    string
	}
	tests := []struct {
		name    string
		args    args
		prepare func(f *mocks.Bucket)
		wantErr bool
	}{
		{
			name: "ListSessionRecords: success",
			args: args{
				sessionID: "1",
				scheme:    "test9",
				urlstr:    "test9://",
			},
			prepare: func(f *mocks.Bucket) {
				f.On("ListPaged", mock.Anything, mock.Anything).Return(&listPage, nil).Maybe()
			},
			wantErr: false,
		},
		{
			name: "ListSessionRecords: returns error, error expected",
			args: args{
				sessionID: "1",
				scheme:    "test10",
				urlstr:    "test10://",
			},
			prepare: func(f *mocks.Bucket) {
				f.On("ErrorCode", mock.Anything).Return(gcerr.Unknown).Maybe()
				f.On("ListPaged", mock.Anything, mock.Anything).Return(nil, errors.New("error")).Maybe()
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blobStorage, fakeDriver := mockNew(t, tt.args.scheme, tt.args.urlstr)

			if tt.prepare != nil {
				tt.prepare(fakeDriver)
			}

			strings, err := blobStorage.ListFileNames(context.Background(), tt.args.sessionID, tt.args.prefix)
			if !tt.wantErr {
				require.NoError(t, err)
				require.NotEmpty(t, strings)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func mockNew(t *testing.T, scheme, urlstr string) (Storage, *mocks.Bucket) {
	fakeDriver := mocks.NewBucket(t)

	// Attributes(ctx context.Context, key string) (*Attributes, error)
	fakeDriver.On("Attributes", mock.AnythingOfType("*context.valueCtx"), "test-bucket").Return(&driver.Attributes{}, nil).Maybe()
	bucket := blob.NewBucket(fakeDriver)
	buo := mocks.NewBucketURLOpener(t)

	// OpenBucketURL(ctx context.Context, u *url.URL) (*Bucket, error)
	buo.On("OpenBucketURL", context.Background(), mock.Anything).Return(bucket, nil)
	blob.DefaultURLMux().RegisterBucket(scheme, buo)

	bucket, err := blob.OpenBucket(context.Background(), urlstr)
	require.NoError(t, err)
	require.NotNil(t, bucket)

	exists, err := bucket.Exists(context.Background(), "test-bucket")
	require.NoError(t, err)
	require.True(t, exists)

	return &blobStorage{bucket: bucket}, fakeDriver
}

var readCloser = ioutil.NopCloser(bytes.NewBufferString(""))

var sessionRecord = &BlobFile{
	FileName:    "fileName",
	ContentType: "text/plain; charset=utf-8",
	Content:     readCloser,
}

type reader struct {
	r     io.Reader
	attrs driver.ReaderAttributes
}

func (r *reader) Read(p []byte) (int, error) {
	return r.r.Read(p)
}

func (r *reader) Close() error {
	return nil
}

func (r *reader) Attributes() *driver.ReaderAttributes {
	return &r.attrs
}

func (r *reader) As(i interface{}) bool { return false }

var r = reader{
	r: bytes.NewReader([]byte("byte")),
	attrs: driver.ReaderAttributes{
		ContentType: "Content-type: application/json",
		ModTime:     time.Now(),
		Size:        1,
	},
}

type writer struct{}

func (w *writer) Write(p []byte) (n int, err error) { return 0, nil }

func (w *writer) Close() error { return nil }

var w = writer{}

var listPage = driver.ListPage{
	Objects: []*driver.ListObject{
		{
			Key:     "key",
			ModTime: time.Now(),
			Size:    1,
			MD5:     []byte{},
			IsDir:   true,
		},
	},
}
