package api

import (
	"testing"

	"github.com/stretchr/testify/assert"

	v1 "github.com/browserkube/browserkube/operator/api/v1"
)

func Test_handler_sortBrowsers(t *testing.T) {
	tests := []struct {
		name     string
		browsers []Browser
		want     []Browser
	}{
		{
			name: "positive",
			browsers: []Browser{
				{
					Type:    v1.TypeWebDriver,
					Version: "121.0",
					Name:    "Chrome",
				},
				{
					Type:    v1.TypeWebDriver,
					Version: "120.0",
					Name:    "Chrome",
				},
				{
					Type:    v1.TypeWebDriver,
					Version: "122.0",
					Name:    "Firefox",
				},
				{
					Type:    v1.TypePlaywright,
					Version: "122.0",
					Name:    "Chrome",
				},
				{
					Type:    v1.TypePlaywright,
					Version: "121.0",
					Name:    "Chrome",
				},
			},
			want: []Browser{
				{
					Type:    v1.TypeWebDriver,
					Version: "121.0",
					Name:    "Chrome",
				},
				{
					Type:    v1.TypeWebDriver,
					Version: "120.0",
					Name:    "Chrome",
				},
				{
					Type:    v1.TypeWebDriver,
					Version: "122.0",
					Name:    "Firefox",
				},
				{
					Type:    v1.TypePlaywright,
					Version: "122.0",
					Name:    "Chrome",
				},
				{
					Type:    v1.TypePlaywright,
					Version: "121.0",
					Name:    "Chrome",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &handler{}
			h.sortBrowsers(tt.browsers)
			assert.Equal(t, tt.want, tt.browsers)
		})
	}
}
