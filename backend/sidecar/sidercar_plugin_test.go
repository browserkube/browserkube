package main

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/browserkube/browserkube/pkg/wd/wdproto"
)

func Test_replaceWebsocketUrl(t *testing.T) {
	wsURLFieldKey := "webSocketUrl"
	tests := []struct {
		name                 string
		rq                   *http.Request
		rsPayload            *wdproto.NewSessionRS
		wantSeleniumWsURL    string
		wantBrowserkubeWsURL string
	}{
		{
			name: "webSocketUrl not found",
			rq:   &http.Request{},
			rsPayload: &wdproto.NewSessionRS{
				Value: wdproto.NewSessionRSValue{
					Capabilities: map[string]interface{}{},
				},
			},
			wantSeleniumWsURL:    "",
			wantBrowserkubeWsURL: "",
		},
		{
			name: "normal result",
			rq: &http.Request{
				Header: map[string][]string{
					"X-Forwarded-Host":   {"example.com"},
					"X-Forwarded-Prefix": {"/browserkube"},
				},
				URL:  must(url.Parse("/wd/hub/session")),
				Host: "example.com",
			},
			rsPayload: &wdproto.NewSessionRS{
				Value: wdproto.NewSessionRSValue{
					SessionID: "1234567",
					Capabilities: map[string]interface{}{
						wsURLFieldKey: "ws://localhost:32451/session/qwerty12345",
					},
				},
			},
			wantSeleniumWsURL:    "ws://localhost:32451/session/qwerty12345",
			wantBrowserkubeWsURL: "ws://example.com/browserkube/wd/hub/bidi/1234567",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := replaceWebsocketURL(tt.rq, tt.rsPayload, tt.rsPayload.Value.SessionID, wsURLFieldKey)
			if got != tt.wantSeleniumWsURL {
				t.Errorf("SeleniumWsURL invalid result. Want: %v, Got: %v", tt.wantSeleniumWsURL, got)
			}

			replacedBrowserkubeUrl, _ := tt.rsPayload.Value.Capabilities[wsURLFieldKey].(string)
			if replacedBrowserkubeUrl != tt.wantBrowserkubeWsURL {
				t.Errorf("BrowserkubeWsURL invalid result. Want: %v, Got: %v", tt.wantBrowserkubeWsURL, replacedBrowserkubeUrl)
			}
		})
	}
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}

	return v
}
