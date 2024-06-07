package opentelemetry

import (
	"net/http"
	"regexp"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/zap"
)

const uuidPlaceholder = "{uuid}"

// regex matching UUID RFC4122.
var rfc4122regex = regexp.MustCompile(`[\da-zA-Z]{8}-[\da-zA-Z]{4}-[\da-zA-Z]{4}-[\da-zA-Z]{4}-[\da-zA-Z]{12}`)

// NewMetricsMiddleware creates new HTTP middleware that populates the context with metrics
func NewMetricsMiddleware(moduleName string) func(http.Handler) http.Handler {
	meter := otel.Meter("")
	log := zap.S()
	requestCount, err := meter.Int64Counter(
		"http.request_count",
		metric.WithDescription("The total number of HTTP requests."),
		metric.WithUnit("request"),
	)
	if err != nil {
		log.Fatalln(err)
	}

	requestLatency, err := meter.Float64Histogram(
		"http.request_latency_ms",
		metric.WithDescription("The latency distribution of HTTP requests."),
		metric.WithUnit("ms"),
	)
	if err != nil {
		log.Fatalln(err)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			attrs := []attribute.KeyValue{
				attribute.String("module", moduleName),
				attribute.String("method", r.Method),
				attribute.String("path", replaceSessionID(r.URL.Path)),
			}

			// Measure the request count and start a timer for latency
			requestCount.Add(ctx, 1, metric.WithAttributes(attrs...))
			startTime := time.Now()

			next.ServeHTTP(w, r)

			// Record the request latency
			requestLatency.Record(ctx, float64(time.Since(startTime).Milliseconds()), metric.WithAttributes(attrs...))
		})
	}
}

func replaceSessionID(urlPath string) string {
	if len(urlPath) < 36 {
		// RFC4122 is 36 characters long
		return urlPath
	}

	return rfc4122regex.ReplaceAllString(urlPath, uuidPlaceholder)
}
