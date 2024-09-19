package opentelemetry

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/caarlos0/env/v11"
	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.uber.org/zap"
)

type config struct {
	TelemetryEnabled      bool   `env:"TELEMETRY_PROVIDER_ENABLED"`
	TelemetryProviderType string `env:"TELEMETRY_PROVIDER_TYPE"`
	TelemetryHost         string `env:"BROWSERKUBE_TEMPO_SERVICE_HOST"`
	OTLPPort              string `env:"BROWSERKUBE_TEMPO_SERVICE_PORT_TEMPO_OTLP_HTTP"`
	ZipkinPort            string `env:"BROWSERKUBE_TEMPO_SERVICE_PORT_TEMPO_ZIPKIN"`
}

const (
	zipkinProvider        = "zipkin"
	otlptracehttpProvider = "otlptracehttp"
)

func InitProvider(serviceName string) (*sdktrace.TracerProvider, error) {
	cfg := config{}

	if err := cfg.readEnv(); err != nil {
		return nil, fmt.Errorf("unable to init config : %w", err)
	}

	if !cfg.TelemetryEnabled {
		return nil, nil //nolint:nilnil
	}

	// Create exporter
	exporter, err := cfg.initExporter()
	if err != nil {
		return nil, fmt.Errorf("unable to init trace exporter: %w", err)
	}

	// Create resource
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to init trace resource: %w", err)
	}

	// Create trace provider with the exporter
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		&propagation.TraceContext{},
		&propagation.Baggage{},
	))

	return tp, nil
}

func HTTPMiddleware(provider *sdktrace.TracerProvider) func(http.Handler) http.Handler {
	opts := []otelhttp.Option{
		otelhttp.WithTracerProvider(provider),
		otelhttp.WithSpanNameFormatter(func(opName string, r *http.Request) string {
			if rc := chi.RouteContext(r.Context()); rc != nil {
				return rc.RoutePattern()
			}
			return "unknown"
		}),
	}
	return otelhttp.NewMiddleware("", opts...)
}

func (c *config) initExporter() (sdktrace.SpanExporter, error) {
	switch c.TelemetryProviderType {
	case zipkinProvider:
		// Zipkin requires http schema before
		url := fmt.Sprintf("http://%s:%s/api/v2/spans", c.TelemetryHost, c.ZipkinPort) //nolint:nosprintfhostport
		exporterZipkin, err := zipkin.New(url)
		if err != nil {
			return nil, fmt.Errorf("creating zipkin trace exporter: %w", err)
		}

		return exporterZipkin, nil
	case otlptracehttpProvider:
		url := fmt.Sprintf("%s:%s", c.TelemetryHost, c.OTLPPort)
		exporterOtlptracehttp, err := otlptracehttp.New(
			context.Background(),
			otlptracehttp.WithInsecure(),
			otlptracehttp.WithEndpoint(url),
		)
		if err != nil {
			return nil, fmt.Errorf("creating otlptracehttp trace exporter: %w", err)
		}

		return exporterOtlptracehttp, nil
	default:
		return nil, errors.New("unexpected TelemetryProviderType")
	}
}

func (c *config) readEnv() error {
	opts := env.Options{
		OnSet: func(tag string, value interface{}, isDefault bool) {
			if tag != "TELEMETRY_PROVIDER_TYPE" {
				return
			}
			if v, ok := value.(string); ok {
				if !compareStrs(v) {
					zap.S().Infof("Value is not %s of allowed values %v", tag, value)
					return
				}
			}
		},
	}
	if err := env.ParseWithOptions(c, opts); err != nil {
		return fmt.Errorf("unable to parse config: %w", err)
	}

	return nil
}

func compareStrs(str string) bool {
	switch str {
	case zipkinProvider, otlptracehttpProvider:
		return true
	default:
		return false
	}
}
