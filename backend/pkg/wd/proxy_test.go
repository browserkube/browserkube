package wd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/browserkube/browserkube/pkg/wd/wdproto"
)

func Test_chain(t *testing.T) {
	type f = func(string) string
	type mFunc = func(f) f

	funcs := []mFunc{
		func(next f) f {
			return func(s string) string {
				s = "first_" + s
				return next(s)
			}
		},
		func(next f) f {
			return func(s string) string {
				s = "second_" + s
				return next(s)
			}
		},
	}
	chainF := chain[f](funcs, func(s string) string {
		return s
	})
	res := chainF("hello")
	fmt.Println(res)
}

func Test_adjustCapabilities(t *testing.T) {
	capsStr := `{"capabilities":{"firstMatch":[{"browserName":"chrome","browserVersion":"108.0","goog:chromeOptions":{}}]}}`
	var rq wdproto.NewSessionRQ
	err := json.Unmarshal([]byte(capsStr), &rq)
	require.NoError(t, err)

	err = adjustCapabilities(&rq)
	require.NoError(t, err)

	require.Equal(t, "chrome", rq.Capabilities.BrowserName)
	require.Equal(t, "108.0", rq.Capabilities.BrowserVersion)
}

func TestProxyManager_cleanupOriginHeaders(t *testing.T) {
	type args struct {
		out *http.Request
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "no headers",
			args: args{out: getDummyRequest(t, map[string]string{})},
		},
		{
			name: "origin header",
			args: args{out: getDummyRequest(t, map[string]string{"Origin": "any"})},
		},
		{
			name: "access control header",
			args: args{out: getDummyRequest(t, map[string]string{"Access-Control-Allow-Origin": "any"})},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ProxyManager{}
			p.cleanupOriginHeaders(tt.args.out)
			require.NotContains(t, tt.args.out.Header, "Origin")
			for h := range tt.args.out.Header {
				require.False(t, strings.HasPrefix(h, "Access-Control"))
			}
		})
	}
}

func getDummyRequest(t *testing.T, headers map[string]string) *http.Request {
	rq, err := http.NewRequest("GET", "", nil)
	require.NoError(t, err)

	for h, v := range headers {
		rq.Header.Set(h, v)
	}

	return rq
}
