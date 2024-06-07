package opentelemetry

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInitProvider(t *testing.T) {
	type args struct {
		serviceName string
		config      *config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "init provider: successful with zipkin exporter",
			args: args{
				serviceName: "api",
				config: &config{
					TelemetryEnabled:      true,
					TelemetryProviderType: "zipkin",
					TelemetryHost:         "localhost",
					ZipkinPort:            "9411",
				},
			},
			wantErr: false,
		},
		{
			name: "init provider: successful with otlptracehttp exporter",
			args: args{
				serviceName: "api",
				config: &config{
					TelemetryEnabled:      true,
					TelemetryProviderType: "otlptracehttp",
					TelemetryHost:         "localhost",
					OTLPPort:              "4318",
				},
			},
			wantErr: false,
		},
		{
			name: "init provider: config not found, error expected",
			args: args{
				serviceName: "api",
				config:      &config{TelemetryEnabled: false},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.config != nil {
				t.Setenv("TELEMETRY_PROVIDER_ENABLED", "true")
				t.Setenv("TELEMETRY_PROVIDER_TYPE", tt.args.config.TelemetryProviderType)
				t.Setenv("BROWSERKUBE_TEMPO_SERVICE_HOST", tt.args.config.TelemetryHost)
				t.Setenv("BROWSERKUBE_TEMPO_SERVICE_PORT_TEMPO_OTLP_HTTP", tt.args.config.OTLPPort)
				t.Setenv("BROWSERKUBE_TEMPO_SERVICE_PORT_TEMPO_ZIPKIN", tt.args.config.ZipkinPort)
			} else {
				t.Setenv("TELEMETRY_PROVIDER_ENABLED", "false")
			}

			got, err := InitProvider(tt.args.serviceName)
			if !tt.wantErr {
				require.NoError(t, err)
				if tt.args.config.TelemetryEnabled {
					require.NotEmpty(t, got)
				}

			} else {
				require.Error(t, err)
			}
		})
	}
}
