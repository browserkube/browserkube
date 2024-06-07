package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/browserkube/browserkube/browser-updater/internal/mocks"
	"github.com/browserkube/browserkube/browser-updater/internal/registry"
	v1 "github.com/browserkube/browserkube/operator/api/v1"
)

func TestUpdater(t *testing.T) {
	tests := []struct {
		name    string
		arg1    *v1.BrowserSetList
		wantErr bool
	}{
		{
			name:    "empty",
			arg1:    &v1.BrowserSetList{},
			wantErr: true,
		},
		{
			name: "happy",
			arg1: &v1.BrowserSetList{
				Items: []v1.BrowserSet{
					{
						Spec: v1.BrowserSetSpec{
							WebDriver: map[string]v1.BrowsersConfig{
								"chrome": {
									Versions: map[string]v1.BrowserConfig{
										"108.0": {
											Provider: "k8s",
											Image:    "selenium/standalone-chrome:108.0",
											Port:     "4444",
										},
										"109.0": {
											Provider: "k8s",
											Image:    "selenium/standalone-chrome:109.0",
											Port:     "4444",
										},
									},
								},
								"firefox": {
									Versions: map[string]v1.BrowserConfig{
										"108.0": {
											Provider: "k8s",
											Image:    "selenium/standalone-firefox:108.0",
											Port:     "4444",
										},
									},
								},
							},
							Playwright: map[string]v1.BrowsersConfig{
								"chrome": {
									Versions: map[string]v1.BrowserConfig{
										"108.0": {
											Provider: "k8s",
											Image:    "quay.io/test/elfs/elrond:111.1",
											Port:     "4444",
										},
									},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "duplicate",
			arg1: &v1.BrowserSetList{
				Items: []v1.BrowserSet{
					{
						Spec: v1.BrowserSetSpec{
							WebDriver: map[string]v1.BrowsersConfig{
								"chrome": {
									Versions: map[string]v1.BrowserConfig{
										"108.0": {
											Provider: "k8s",
											Image:    "selenium/standalone-chrome:108.0",
											Port:     "4444",
										},
										"109.0": {
											Provider: "k8s",
											Image:    "selenium/standalone-chrome:108.0",
											Port:     "4444",
										},
									},
								},
								"firefox": {
									Versions: map[string]v1.BrowserConfig{
										"108.0": {
											Provider: "k8s",
											Image:    "selenium/standalone-firefox:108.0",
											Port:     "4444",
										},
									},
								},
							},
							Playwright: map[string]v1.BrowsersConfig{
								"chrome": {
									Versions: map[string]v1.BrowserConfig{
										"108.0": {
											Provider: "k8s",
											Image:    "quay.io/test/elfs/elrond:111.1",
											Port:     "4444",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	mockBrowsersInterface := mocks.NewBrowsersInterface(t)
	mockBrowserSetsInterface := mocks.NewBrowsersSetsInterface(t)
	mockRegistryClient := mocks.NewRegistryClient(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRegistryClient.On("CheckRegistry", mock.AnythingOfType("string"), "selenium/standalone-chrome").Return(&registry.RegistryImageListResp{
				Name: "selenium/standalone-chrome",
				Tags: []string{"108.0", "109.0, 110.0"},
			}, nil).Maybe()
			mockRegistryClient.On("CheckRegistry", mock.AnythingOfType("string"), "selenium/standalone-firefox").Return(&registry.RegistryImageListResp{
				Name: "selenium/standalone-firefox",
				Tags: []string{"108.0", "109.0"},
			}, nil).Maybe()
			mockRegistryClient.On("CheckRegistry", mock.AnythingOfType("string"), "quay.io/test/elfs/elrond").Return(&registry.RegistryImageListResp{
				Name: "quay.io/test/elfs/elrond",
				Tags: []string{"new", "newer", "newest"},
			}, nil).Maybe()
			updater := &BrowserImageUpdater{
				[]string{},
				make(chan struct{}),
				mockBrowserSetsInterface,
				mockBrowsersInterface,
				mockRegistryClient,
			}
			go func() {
				for range updater.BrowserCache {
				}
			}()
			updated, err := updater.UpdateBrowserSet(tt.arg1)
			if tt.wantErr {
				assert.NotNil(t, err)
				return
			}
			assert.Nil(t, err)
			assert.Equal(t, 3, len(updated.Items[0].Spec.WebDriver["chrome"].Versions))
			assert.Equal(t, 2, len(updated.Items[0].Spec.WebDriver["firefox"].Versions))
		})
	}
}
