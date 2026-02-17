package opentelemetry

import (
	"context"
	"fmt"
	"net/http"

	"github.com/caarlos0/env/v11"
	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type config struct {
	TelemetryEnabled bool   `env:"TELEMETRY_PROVIDER_ENABLED"`
	TelemetryHost    string `env:"BROWSERKUBE_TEMPO_SERVICE_HOST"`
	OTLPPort         string `env:"BROWSERKUBE_TEMPO_SERVICE_PORT_TEMPO_OTLP_HTTP"`
	ZipkinPort       string `env:"BROWSERKUBE_TEMPO_SERVICE_PORT_TEMPO_ZIPKIN"`
}

const (
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
}

func (c *config) readEnv() error {
	opts := env.Options{}
	if err := env.ParseWithOptions(c, opts); err != nil {
		return fmt.Errorf("unable to parse config: %w", err)
	}

	return nil
}
