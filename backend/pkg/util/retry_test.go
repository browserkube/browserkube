package browserkubeutil

import (
	"errors"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"k8s.io/utils/ptr"
)

func TestRetryWithTimeout(t *testing.T) {
	i := ptr.To[int32](0)
	_, err := RetryWithTimeout(10*time.Second, 0, 500*time.Millisecond, func() (res interface{}, err error) {
		*i++

		if *i < 2 {
			return nil, errors.New("error")
		}
		return i, nil
	})
	require.NoError(t, err)
	require.Equal(t, int32(2), *i)
}

func TestRetryWithTimeoutAttempts(t *testing.T) {
	_, err := RetryWithTimeout(1*time.Second, 0, 500*time.Millisecond, func() (res interface{}, err error) {
		return nil, errors.New("error")
	})
	require.Error(t, err)
	require.Regexp(t, regexp.MustCompile("after \\d attempts, last error: error"), err.Error())
}

func TestRetry(t *testing.T) {
	type args struct {
		attempts int
		timeout  time.Duration
		callback func() (interface{}, error)
	}
	type testCase struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}
	tests := []testCase{
		{
			name: "empty",
			args: args{
				attempts: 1,
				timeout:  1,
				callback: func() (res interface{}, err error) {
					return res, err
				},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "err",
			args: args{
				attempts: 1,
				timeout:  1,
				callback: func() (res interface{}, err error) {
					return res, errors.New("some horrible error")
				},
			},
			want:    nil,
			wantErr: true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Retry(tt.args.attempts, tt.args.timeout, tt.args.callback)
			if (err != nil) != tt.wantErr {
				t.Errorf("Retry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Retry() = %v, want %v", got, tt.want)
			}
		})
	}
}
