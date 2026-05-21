package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad_Defaults(t *testing.T) {
	os.Clearenv()
	t.Setenv("JWT_SECRET", "test-secret")
	Load()

	if AppConfig.AppPort != "8080" {
		t.Errorf("expected default AppPort '8080', got %q", AppConfig.AppPort)
	}
	if AppConfig.PGPort != "5432" {
		t.Errorf("expected default PGPort '5432', got %q", AppConfig.PGPort)
	}
	if AppConfig.ShutdownTimeout != 10*time.Second {
		t.Errorf("expected default ShutdownTimeout 10s, got %v", AppConfig.ShutdownTimeout)
	}
	if AppConfig.RateLimitBurst != 200 {
		t.Errorf("expected default RateLimitBurst 100, got %d", AppConfig.RateLimitBurst)
	}
	if AppConfig.PGMaxOpenConns != 25 {
		t.Errorf("expected default PGMaxOpenConns 25, got %d", AppConfig.PGMaxOpenConns)
	}
	if AppConfig.CORSAllowOrigins != "*" {
		t.Errorf("expected default CORSAllowOrigins '*', got %q", AppConfig.CORSAllowOrigins)
	}
	if AppConfig.JWTAccessExpiry != "15m" {
		t.Errorf("expected default JWTAccessExpiry '15m', got %q", AppConfig.JWTAccessExpiry)
	}
	if AppConfig.JWTRefreshExpiry != "168h" {
		t.Errorf("expected default JWTRefreshExpiry '168h', got %q", AppConfig.JWTRefreshExpiry)
	}
	if AppConfig.NATSURL != "nats://localhost:4222" {
		t.Errorf("expected default NATSURL 'nats://localhost:4222', got %q", AppConfig.NATSURL)
	}
}

func TestLoad_CustomEnv(t *testing.T) {
	os.Clearenv()
	t.Setenv("JWT_SECRET", "custom-secret")
	t.Setenv("APP_PORT", "9090")
	t.Setenv("PG_HOST", "myhost")
	t.Setenv("SHUTDOWN_TIMEOUT", "30s")
	t.Setenv("RATE_LIMIT_BURST", "200")
	t.Setenv("CORS_ALLOW_ORIGINS", "https://example.com")

	Load()

	if AppConfig.AppPort != "9090" {
		t.Errorf("expected AppPort '9090', got %q", AppConfig.AppPort)
	}
	if AppConfig.PGHost != "myhost" {
		t.Errorf("expected PGHost 'myhost', got %q", AppConfig.PGHost)
	}
	if AppConfig.ShutdownTimeout != 30*time.Second {
		t.Errorf("expected ShutdownTimeout 30s, got %v", AppConfig.ShutdownTimeout)
	}
	if AppConfig.RateLimitBurst != 200 {
		t.Errorf("expected RateLimitBurst 200, got %d", AppConfig.RateLimitBurst)
	}
	if AppConfig.CORSAllowOrigins != "https://example.com" {
		t.Errorf("expected CORSAllowOrigins 'https://example.com', got %q", AppConfig.CORSAllowOrigins)
	}
}

func TestParseDuration_Invalid(t *testing.T) {
	d := parseDuration("invalid")
	if d != 10*time.Second {
		t.Errorf("expected fallback 10s, got %v", d)
	}
}

func TestParseInt_Invalid(t *testing.T) {
	v := parseInt("abc")
	if v != 0 {
		t.Errorf("expected fallback 0, got %d", v)
	}
}
