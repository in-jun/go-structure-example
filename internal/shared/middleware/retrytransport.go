package middleware

import (
	"crypto/rand"
	"log/slog"
	"math/big"
	"net/http"
	"time"
)

type RetryTransport struct {
	Base       http.RoundTripper
	MaxRetries int
	BaseDelay  time.Duration
}

func NewRetryTransport(base http.RoundTripper, maxRetries int) *RetryTransport {
	if base == nil {
		base = http.DefaultTransport
	}
	return &RetryTransport{
		Base:       base,
		MaxRetries: maxRetries,
		BaseDelay:  100 * time.Millisecond,
	}
}

func isIdempotent(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodPut, http.MethodDelete:
		return true
	}
	return false
}

func isRetryable(statusCode int) bool {
	switch statusCode {
	case http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return true
	}
	return false
}

func (t *RetryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	canRetry := isIdempotent(req.Method) || req.Header.Get("Idempotency-Key") != ""

	resp, err := t.Base.RoundTrip(req)
	if !canRetry {
		return resp, err
	}

	delay := t.BaseDelay
	for attempt := 1; attempt <= t.MaxRetries; attempt++ {
		if err == nil && !isRetryable(resp.StatusCode) {
			return resp, nil
		}

		if resp != nil {
			if err := resp.Body.Close(); err != nil {
				slog.Warn("failed to close upstream response body", "error", err)
			}
		}

		n, _ := rand.Int(rand.Reader, big.NewInt(int64(delay/2)))
		sleepTime := delay + time.Duration(n.Int64())

		slog.Warn("upstream request failed, retrying",
			"attempt", attempt,
			"max_retries", t.MaxRetries,
			"delay", sleepTime,
			"error", err,
		)

		timer := time.NewTimer(sleepTime)
		select {
		case <-req.Context().Done():
			timer.Stop()
			return nil, req.Context().Err()
		case <-timer.C:
		}

		resp, err = t.Base.RoundTrip(req)
		delay *= 2
	}

	return resp, err
}
