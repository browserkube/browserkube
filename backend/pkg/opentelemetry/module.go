package opentelemetry

import (
	"context"
	"fmt"
	"time"

	"github.com/caarlos0/env/v9"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/time/rate"

	browserkubev1 "github.com/browserkube/browserkube/operator/api/v1"
	wdsession "github.com/browserkube/browserkube/pkg/session"
	"github.com/browserkube/browserkube/pkg/sessionresult"
	"github.com/browserkube/browserkube/storage"
)

const usageObservePeriod = 5 * time.Minute

var Module = fx.Options(
	fx.Provide(
		provideConfig,
		fx.Annotate(
			provideMetricsProxyPlugin,
			fx.ResultTags(`group:"wd-extensions"`),
		),
	),
	fx.Invoke(
		registerActiveSessionsObserver,
		// registerSessionResultsObserver,
		// registerStorageUsageObserver,
		initMetricsProvider,
	),
)

const srvName = "browserkube-backend"

type telemetryConfig struct {
	Prometheus struct {
		Enabled bool `env:"TELEMETRY_PROMETHEUS_ENABLED" envDefault:"false"`
	}

	Environment string `env:"ENVIRONMENT"       envDefault:"dev"`
	Version     string `env:"TELEMETRY_VERSION" envDefault:"1.0.0"`
}

func provideConfig() (*telemetryConfig, error) {
	var cfg telemetryConfig
	return &cfg, errors.WithStack(env.Parse(&cfg))
}

func initMetricsProvider(cfg *telemetryConfig, mux chi.Router) error {
	if !cfg.Prometheus.Enabled {
		return nil
	}

	exporter, err := prometheus.New()
	if err != nil {
		return errors.Wrap(err, "failed to initialize prometheus exporter")
	}

	otel.SetMeterProvider(sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(exporter),
		sdkmetric.WithResource(buildResource(cfg)),
	))

	zap.S().Info("Starting metrics endpoint")
	mux.Handle("/metrics", promhttp.Handler())

	return nil
}

func registerActiveSessionsObserver(sessionRepo wdsession.Repository) error {
	meter := otel.Meter("")
	activeSessions, err := meter.Int64ObservableGauge(
		"browserkube.active_sessions",
		metric.WithDescription("Number of active sessions"),
		metric.WithUnit("session"),
	)
	if err != nil {
		return errors.Wrap(err, "active_sessions metric registration failed")
	}

	runningBrowsers, err := meter.Int64ObservableGauge(
		"browserkube.browsers_running",
		metric.WithDescription("Number of running browsers"),
		metric.WithUnit("browser"),
	)
	if err != nil {
		return errors.Wrap(err, "browsers_running metric registration failed")
	}

	pendingBrowsers, err := meter.Int64ObservableGauge(
		"browserkube.browsers_pending",
		metric.WithDescription("Number of pending browsers"),
		metric.WithUnit("browser"),
	)
	if err != nil {
		return errors.Wrap(err, "browsers_pending metric registration failed")
	}

	terminatedBrowsers, err := meter.Int64ObservableGauge(
		"browserkube.browsers_terminated",
		metric.WithDescription("Number of terminated browsers"),
		metric.WithUnit("browser"),
	)
	if err != nil {
		return errors.Wrap(err, "browsers_terminated metric registration failed")
	}

	_, err = meter.RegisterCallback(func(ctx context.Context, o metric.Observer) error {
		sessions, fErr := sessionRepo.FindAll()
		if fErr != nil {
			return errors.Wrap(err, "failed to query sessions repo")
		}

		var runningC, pendingC, terminatedC int64
		for _, sess := range sessions {
			switch sess.Browser.Status.Phase {
			case browserkubev1.PhaseRunning:
				runningC++
			case browserkubev1.PhasePending:
				pendingC++
			case browserkubev1.PhaseTerminated:
				terminatedC++
			}
		}

		o.ObserveInt64(activeSessions, int64(len(sessions)))
		o.ObserveInt64(runningBrowsers, runningC)
		o.ObserveInt64(pendingBrowsers, pendingC)
		o.ObserveInt64(terminatedBrowsers, terminatedC)
		return nil
	}, activeSessions, runningBrowsers, pendingBrowsers, terminatedBrowsers)
	if err != nil {
		return errors.Wrap(err, "callback registration for active sessions failed")
	}

	return nil
}

func registerSessionResultsObserver(resultsRepo sessionresult.Repository) error { //nolint:deadcode
	meter := otel.Meter("")
	activeSessions, err := meter.Int64ObservableGauge(
		"browserkube.session_results",
		metric.WithDescription("Number of session results"),
		metric.WithUnit("results"),
	)
	if err != nil {
		return errors.Wrap(err, "session_results metric registration failed")
	}

	_, err = meter.RegisterCallback(func(ctx context.Context, o metric.Observer) error {
		resultsPage, fErr := resultsRepo.FindAll(ctx, 0, "")
		if fErr != nil {
			return errors.Wrap(err, "failed to query resultsPage repo")
		}
		o.ObserveInt64(activeSessions, int64(len(resultsPage.Items)))
		return nil
	}, activeSessions)
	if err != nil {
		return errors.Wrap(err, "callback registration for session results failed")
	}

	return nil
}

func registerStorageUsageObserver(st storage.BlobSessionStorage) error { //nolint:deadcode
	meter := otel.Meter("")
	storageUsed, err := meter.Int64ObservableGauge(
		"browserkube.storage_usage",
		metric.WithDescription("Size of occupied storage"),
		metric.WithUnit("byte"),
	)
	if err != nil {
		return errors.Wrap(err, "storage_usage metric registration failed")
	}

	rateLimit := rate.Sometimes{Interval: usageObservePeriod}
	_, err = meter.RegisterCallback(func(ctx context.Context, o metric.Observer) error {
		rateLimit.Do(func() {
			used, uErr := st.SizeUsed()
			if uErr != nil {
				zap.S().Error(uErr)
				return
			}
			o.ObserveInt64(storageUsed, used)
		})
		return nil
	}, storageUsed)
	if err != nil {
		return errors.Wrap(err, "callback registration for storage usage failed")
	}

	return nil
}

func buildResource(cfg *telemetryConfig) *sdkresource.Resource {
	var svcName string
	if cfg.Environment != "" {
		svcName = fmt.Sprintf("%s/%s", cfg.Environment, srvName)
	} else {
		svcName = srvName
	}
	return sdkresource.NewSchemaless(
		semconv.ServiceNameKey.String(svcName),
		semconv.ServiceVersionKey.String(cfg.Version),
		semconv.DeploymentEnvironmentKey.String(cfg.Environment),
	)
}
