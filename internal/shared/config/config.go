package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	AppPort            string
	JWTSecret          string
	JWTAccessExpiry    string
	JWTRefreshExpiry   string
	RedisURL           string
	PGHost             string
	PGPort             string
	PGDatabase         string
	PGUsername         string
	PGPassword         string
	NATSURL            string
	MigrationPath      string
	AuctionServiceURL  string
	AuctionGRPCAddress string
	GRPCPort           string
	AuthServiceURL     string
	BidServiceURL      string
	PaymentServiceURL  string
	PGSSLMode          string
	PGMaxOpenConns     int
	PGMaxIdleConns     int
	CORSAllowOrigins   string
	OTELEndpoint       string
	MigrationMode      string
	ShutdownTimeout    time.Duration
	RateLimitRPS       float64
	RateLimitBurst     int
}

var AppConfig Config

func Load() {
	AppConfig = Config{
		AppPort:            getEnv("APP_PORT", "8080"),
		JWTSecret:          getEnv("JWT_SECRET", ""),
		JWTAccessExpiry:    getEnv("JWT_ACCESS_EXPIRY", "15m"),
		JWTRefreshExpiry:   getEnv("JWT_REFRESH_EXPIRY", "168h"),
		RedisURL:           getEnv("REDIS_URL", "localhost:6379"),
		PGHost:             getEnv("PG_HOST", "localhost"),
		PGPort:             getEnv("PG_PORT", "5432"),
		PGDatabase:         getEnv("PG_DATABASE", "app_db"),
		PGUsername:         getEnv("PG_USERNAME", "postgres"),
		PGPassword:         getEnv("PG_PASSWORD", ""),
		NATSURL:            getEnv("NATS_URL", "nats://localhost:4222"),
		MigrationPath:      getEnv("MIGRATION_PATH", "migrations"),
		AuctionServiceURL:  getEnv("AUCTION_SERVICE_URL", "http://localhost:8082"),
		AuctionGRPCAddress: getEnv("AUCTION_GRPC_ADDRESS", "localhost:9090"),
		GRPCPort:           getEnv("GRPC_PORT", ""),
		AuthServiceURL:     getEnv("AUTH_SERVICE_URL", "http://localhost:8081"),
		BidServiceURL:      getEnv("BID_SERVICE_URL", "http://localhost:8083"),
		PaymentServiceURL:  getEnv("PAYMENT_SERVICE_URL", "http://localhost:8084"),
		PGSSLMode:          getEnv("PG_SSL_MODE", "disable"),
		PGMaxOpenConns:     parseInt(getEnv("PG_MAX_OPEN_CONNS", "25")),
		PGMaxIdleConns:     parseInt(getEnv("PG_MAX_IDLE_CONNS", "10")),
		CORSAllowOrigins:   getEnv("CORS_ALLOW_ORIGINS", "*"),
		OTELEndpoint:       getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://tempo:4318"),
		MigrationMode:      getEnv("MIGRATION_MODE", "auto"),
		ShutdownTimeout:    parseDuration(getEnv("SHUTDOWN_TIMEOUT", "10s")),
		RateLimitRPS:       parseFloat(getEnv("RATE_LIMIT_RPS", "100")),
		RateLimitBurst:     parseInt(getEnv("RATE_LIMIT_BURST", "200")),
	}
}

func requireEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		log.Fatalf("required environment variable %s is not set", key)
	}
	return value
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func parseFloat(s string) float64 {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 100
	}
	return v
}

func parseInt(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return v
}

func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 10 * time.Second
	}
	return d
}
