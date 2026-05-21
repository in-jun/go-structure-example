package middleware

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/in-jun/go-structure-example/internal/shared/observability"
	"github.com/sony/gobreaker/v2"
)

type CircuitBreakerTransport struct {
	base http.RoundTripper
	cb   *gobreaker.TwoStepCircuitBreaker[any]
	name string
}

func NewCircuitBreakerTransport(base http.RoundTripper, name string) *CircuitBreakerTransport {
	settings := gobreaker.Settings{
		Name:        name,
		MaxRequests: 3,
		Interval:    60 * time.Second,
		Timeout:     30 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > 5
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			slog.Warn("circuit breaker state change",
				"service", name,
				"from", from.String(),
				"to", to.String(),
			)
			stateValue := float64(0)
			switch to {
			case gobreaker.StateHalfOpen:
				stateValue = 1
			case gobreaker.StateOpen:
				stateValue = 2
				observability.CircuitBreakerTrips.WithLabelValues(name).Inc()
			}
			observability.CircuitBreakerState.WithLabelValues(name).Set(stateValue)
		},
	}

	return &CircuitBreakerTransport{
		base: base,
		cb:   gobreaker.NewTwoStepCircuitBreaker[any](settings),
		name: name,
	}
}

func (t *CircuitBreakerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	done, err := t.cb.Allow()
	if err != nil {
		slog.Warn("circuit breaker rejected request",
			"service", t.name,
			"path", req.URL.Path,
		)
		return &http.Response{
			StatusCode: http.StatusServiceUnavailable,
			Status:     "503 Service Unavailable",
			Body:       io.NopCloser(strings.NewReader(`{"status":503,"code":"SERVICE_UNAVAILABLE","message":"Service temporarily unavailable"}`)),
			Header:     http.Header{"Content-Type": {"application/json"}},
			Request:    req,
		}, nil
	}

	resp, err := t.base.RoundTrip(req)
	if err != nil {
		done(err)
		return nil, err
	}

	if resp.StatusCode >= 500 {
		done(fmt.Errorf("upstream %s returned %d", t.name, resp.StatusCode))
	} else {
		done(nil)
	}
	return resp, nil
}

func (t *CircuitBreakerTransport) State() gobreaker.State {
	return t.cb.State()
}
