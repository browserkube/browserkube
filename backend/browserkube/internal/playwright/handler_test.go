package playwright

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/browserkube/browserkube/browserkube/internal/playwright/mocks"
	"github.com/browserkube/browserkube/browserkube/internal/provision"
	browserkubev1 "github.com/browserkube/browserkube/operator/api/v1"
)

func Test_playwrightProxy_start(t *testing.T) {
	type args struct {
		browser          string
		enableVNC        string
		enableVideo      string
		screenResolution string
		videoFileName    string
	}
	tests := []struct {
		name                   string
		args                   args
		wantErr                bool
		prepareMockProvisioner func(mockProvisioner *mocks.Provisioner)
		prepareMockCore        func(mockCore *mocks.Core)
	}{
		{
			name: "playwrightProxy_start: success",
			args: args{
				browser:          "chrome",
				enableVNC:        "true",
				enableVideo:      "true",
				screenResolution: "1920x1080x24",
			},
			wantErr: false,
			prepareMockProvisioner: func(mockProvisioner *mocks.Provisioner) {
				mockProvisioner.On("Provision", mock.Anything, mock.Anything, mock.Anything).Return(b, nil).Maybe()
				mockProvisioner.On("Delete", mock.Anything, mock.Anything).Return(nil).Maybe()
				mockProvisioner.On("Logs", mock.Anything, mock.Anything, mock.Anything).Return(io.NopCloser(strings.NewReader("some logs")), nil).Maybe()
			},
			prepareMockCore: func(mockCore *mocks.Core) {
				core := NewNopCore()
				mockCore.On("With", mock.Anything).Return(core).Maybe()
				mockCore.On("Enabled", mock.Anything).Return(false).Maybe()
			},
		},
		{
			name: "playwrightProxy_start: return error, error expected",
			args: args{
				browser: "chrome",
			},
			wantErr: true,
			prepareMockProvisioner: func(mockProvisioner *mocks.Provisioner) {
				mockProvisioner.On("Provision", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("error")).Maybe()
			},
			prepareMockCore: func(mockCore *mocks.Core) {
				core := NewNopCore()
				mockCore.On("With", mock.Anything).Return(core).Maybe()
			},
		},
		{
			name: "playwrightProxy_start: browser not provided, error expected",
			args: args{
				browser: "",
			},
			wantErr: true,
			prepareMockCore: func(mockCore *mocks.Core) {
				core := NewNopCore()
				mockCore.On("With", mock.Anything).Return(core).Maybe()
			},
		},
		{
			name: "playwrightProxy_start: unable to delete browser, error not expected",
			args: args{
				browser: "chrome",
			},
			wantErr: false,
			prepareMockProvisioner: func(mockProvisioner *mocks.Provisioner) {
				mockProvisioner.On("Provision", mock.Anything, mock.Anything, mock.Anything).Return(b, nil).Maybe()
				mockProvisioner.On("Delete", mock.Anything, mock.Anything).Return(errors.New("error")).Maybe()
				mockProvisioner.On("Logs", mock.Anything, mock.Anything, mock.Anything).Return(io.NopCloser(strings.NewReader("some logs")), nil).Maybe()
			},
			prepareMockCore: func(mockCore *mocks.Core) {
				core := NewNopCore()
				mockCore.On("With", mock.Anything).Return(core).Maybe()
				mockCore.On("Enabled", mock.Anything).Return(false).Maybe()
			},
		},
		{
			name: "playwrightProxy_start: fail to decode query, error expected",
			args: args{
				browser:       "chrome",
				videoFileName: "videoFileName",
			},
			wantErr: true,
			prepareMockProvisioner: func(mockProvisioner *mocks.Provisioner) {
				mockProvisioner.On("Provision", mock.Anything, mock.Anything, mock.Anything).Return(b, nil).Maybe()
				mockProvisioner.On("Delete", mock.Anything, mock.Anything).Return(nil).Maybe()
				mockProvisioner.On("Logs", mock.Anything, mock.Anything, mock.Anything).Return(io.NopCloser(strings.NewReader("some logs")), nil).Maybe()
			},
			prepareMockCore: func(mockCore *mocks.Core) {
				core := NewNopCore()
				mockCore.On("With", mock.Anything).Return(core).Maybe()
				mockCore.On("Enabled", mock.Anything).Return(false).Maybe()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProvisioner := mocks.NewProvisioner(t)
			if tt.prepareMockProvisioner != nil {
				tt.prepareMockProvisioner(mockProvisioner)
			}

			mockCore := mocks.NewCore(t)
			if tt.prepareMockCore != nil {
				tt.prepareMockCore(mockCore)
			}

			url := fmt.Sprintf(
				"/playwright/%s/?enableVNC=%s&enableVideo=%s&screenResolution=%s",
				tt.args.browser,
				tt.args.enableVNC,
				tt.args.enableVideo,
				tt.args.screenResolution,
			)

			if tt.args.videoFileName != "" {
				url = fmt.Sprintf(
					"/playwright/%s/?videoFileName=%s",
					tt.args.browser,
					tt.args.videoFileName,
				)
			}

			req, err := http.NewRequest("GET", url, nil)
			require.NoError(t, err)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("browser", tt.args.browser)
			rctx.URLParams.Add("enableVNC", tt.args.enableVNC)
			rctx.URLParams.Add("enableVideo", tt.args.enableVideo)
			rctx.URLParams.Add("screenResolution", tt.args.screenResolution)

			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()

			mockLogger := zap.New(mockCore).Sugar()

			pp := newMockPlaywrightProxy(mockLogger, mockProvisioner)

			err = pp.start(rr, req)
			if !tt.wantErr {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
			require.NotNil(t, mockLogger)
		})
	}
}

var b = &browserkubev1.Browser{
	Status: browserkubev1.BrowserStatus{
		PodName: "chrome",
		Host:    "host",
		PortConfig: browserkubev1.PortConfig{
			Browser: "123",
		},
	},
}

type nopCore struct{}

func NewNopCore() zapcore.Core                                                        { return nopCore{} }
func (nopCore) Enabled(zapcore.Level) bool                                            { return false }
func (n nopCore) With([]zapcore.Field) zapcore.Core                                   { return n }
func (nopCore) Check(_ zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry { return ce }
func (nopCore) Write(zapcore.Entry, []zapcore.Field) error                            { return nil }
func (nopCore) Sync() error                                                           { return nil }
func newMockPlaywrightProxy(logger *zap.SugaredLogger, manager provision.Provisioner) *playwrightProxy {
	return &playwrightProxy{
		logger:        logger,
		manager:       manager,
		sessionRecord: false,
	}
}
