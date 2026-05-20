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
	if AppConfig.MySQLPort != "3306" {
		t.Errorf("expected default MySQLPort '3306', got %q", AppConfig.MySQLPort)
	}
	if AppConfig.ShutdownTimeout != 10*time.Second {
		t.Errorf("expected default ShutdownTimeout 10s, got %v", AppConfig.ShutdownTimeout)
	}
	if AppConfig.RateLimitBurst != 100 {
		t.Errorf("expected default RateLimitBurst 100, got %d", AppConfig.RateLimitBurst)
	}
	if AppConfig.MySQLMaxOpenConns != 25 {
		t.Errorf("expected default MySQLMaxOpenConns 25, got %d", AppConfig.MySQLMaxOpenConns)
	}
	if AppConfig.CORSAllowOrigins != "*" {
		t.Errorf("expected default CORSAllowOrigins '*', got %q", AppConfig.CORSAllowOrigins)
	}
}

func TestLoad_CustomEnv(t *testing.T) {
	os.Clearenv()
	t.Setenv("JWT_SECRET", "custom-secret")
	t.Setenv("APP_PORT", "9090")
	t.Setenv("MYSQL_HOST", "myhost")
	t.Setenv("SHUTDOWN_TIMEOUT", "30s")
	t.Setenv("RATE_LIMIT_BURST", "200")
	t.Setenv("CORS_ALLOW_ORIGINS", "https://example.com")

	Load()

	if AppConfig.AppPort != "9090" {
		t.Errorf("expected AppPort '9090', got %q", AppConfig.AppPort)
	}
	if AppConfig.MySQLHost != "myhost" {
		t.Errorf("expected MySQLHost 'myhost', got %q", AppConfig.MySQLHost)
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
