package observability

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	HTTPRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"service", "method", "path", "status"},
	)

	HTTPRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"service", "method", "path"},
	)

	NATSEventsPublished = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "nats_events_published_total",
			Help: "Total number of NATS events published",
		},
		[]string{"service", "event_type"},
	)

	NATSEventsConsumed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "nats_events_consumed_total",
			Help: "Total number of NATS events consumed",
		},
		[]string{"service", "event_type", "status"},
	)

	GRPCRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_requests_total",
			Help: "Total number of gRPC requests",
		},
		[]string{"service", "method", "status"},
	)

	GRPCRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_request_duration_seconds",
			Help:    "gRPC request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"service", "method"},
	)

	OutboxPendingGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "outbox_pending_events",
			Help: "Number of unpublished events in the outbox",
		},
		[]string{"service"},
	)

	CircuitBreakerState = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "circuit_breaker_state",
			Help: "Circuit breaker state (0=closed, 1=half-open, 2=open)",
		},
		[]string{"service"},
	)

	CircuitBreakerTrips = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "circuit_breaker_trips_total",
			Help: "Total number of circuit breaker trips to open state",
		},
		[]string{"service"},
	)
)

func InitMetrics() {
	prometheus.MustRegister(
		HTTPRequestsTotal,
		HTTPRequestDuration,
		GRPCRequestsTotal,
		GRPCRequestDuration,
		NATSEventsPublished,
		NATSEventsConsumed,
		OutboxPendingGauge,
		CircuitBreakerState,
		CircuitBreakerTrips,
	)
}

func MetricsHandler() http.Handler {
	return promhttp.Handler()
}
