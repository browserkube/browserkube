package opentelemetry

import (
	"strconv"
	"testing"
)

func Test_replaceSessionId(t *testing.T) {
	tests := []struct {
		urlPath string
		want    string
	}{
		{
			"/events",
			"/events",
		}, {
			"/wd/hub/session",
			"/wd/hub/session",
		}, {
			"/wd/hub/session/",
			"/wd/hub/session/",
		}, {
			"/vnc/06b77a7d-6b7a-417c-bf9a-62c346037715",
			"/vnc/" + uuidPlaceholder,
		}, {
			"/clipboard/06b77a7d-6b7a-417c-bf9a-62c346037715",
			"/clipboard/" + uuidPlaceholder,
		}, {
			"/wd/hub/session/06b77a7d-6b7a-417c-bf9a-62c346037715/browserkube/downloads",
			"/wd/hub/session/" + uuidPlaceholder + "/browserkube/downloads",
		}, {
			"/wd/hub/session/06b77a7d-6b7a-417c-bf9a-62c346037715/element/e8ffb0ea-b55d-4f84-bf3a-a5d8777acd06/screenshot",
			"/wd/hub/session/" + uuidPlaceholder + "/element/" + uuidPlaceholder + "/screenshot",
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if got := replaceSessionID(tt.urlPath); got != tt.want {
				t.Errorf("replaceSessionID() = %v, want %v", got, tt.want)
			}
		})
	}
}
