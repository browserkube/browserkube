package opentelemetry

import (
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"

	"github.com/browserkube/browserkube/pkg/session"
	"github.com/browserkube/browserkube/pkg/wd"
)

func provideMetricsProxyPlugin() wd.PluginOpts {
	return wd.PluginOpts{
		Weight: 1,
		Opts: []wd.PluginOpt{
			wd.WithQuitSession(onQuitSession()),
		},
	}
}

func onQuitSession() func(next wd.OnSessionQuit) wd.OnSessionQuit {
	log := zap.S()
	meter := otel.Meter("")
	sessionDuration, err := meter.Int64Histogram(
		"browserkube.session_duration_seconds",
		metric.WithDescription("The duration distribution of Browserkube sessions."),
		metric.WithUnit("second"),
	)
	if err != nil {
		log.Fatalln(err)
	}

	return func(next wd.OnSessionQuit) wd.OnSessionQuit {
		return func(ctx *wd.Context, s *session.Session) error {
			dur := time.Since(s.Browser.CreationTimestamp.Time).Seconds()
			sessionDuration.Record(ctx, int64(dur))
			return next(ctx, s)
		}
	}
}
