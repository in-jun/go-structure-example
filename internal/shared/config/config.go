package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	AppPort          string
	JWTSecret        string
	JWTAccessExpiry  string
	JWTRefreshExpiry string
	MySQLHost        string
	MySQLPort        string
	MySQLDatabase    string
	MySQLUsername    string
	MySQLPassword    string
	MySQLMaxOpenConns int
	MySQLMaxIdleConns int
	RedisHost        string
	RedisPort        string
	RedisPassword    string
	CORSAllowOrigins string
	ShutdownTimeout  time.Duration
	RateLimitBurst   int
}

var AppConfig Config

func Load() {
	AppConfig = Config{
		AppPort:           getEnv("APP_PORT", "8080"),
		JWTSecret:         requireEnv("JWT_SECRET"),
		JWTAccessExpiry:   getEnv("JWT_ACCESS_EXPIRY", "15m"),
		JWTRefreshExpiry:  getEnv("JWT_REFRESH_EXPIRY", "168h"),
		MySQLHost:         getEnv("MYSQL_HOST", "localhost"),
		MySQLPort:         getEnv("MYSQL_PORT", "3306"),
		MySQLDatabase:     getEnv("MYSQL_DATABASE", "app_db"),
		MySQLUsername:     getEnv("MYSQL_USERNAME", "root"),
		MySQLPassword:     getEnv("MYSQL_PASSWORD", ""),
		MySQLMaxOpenConns: parseInt(getEnv("MYSQL_MAX_OPEN_CONNS", "25")),
		MySQLMaxIdleConns: parseInt(getEnv("MYSQL_MAX_IDLE_CONNS", "10")),
		RedisHost:         getEnv("REDIS_HOST", "localhost"),
		RedisPort:         getEnv("REDIS_PORT", "6379"),
		RedisPassword:     getEnv("REDIS_PASSWORD", ""),
		CORSAllowOrigins:  getEnv("CORS_ALLOW_ORIGINS", "*"),
		ShutdownTimeout:   parseDuration(getEnv("SHUTDOWN_TIMEOUT", "10s")),
		RateLimitBurst:    parseInt(getEnv("RATE_LIMIT_BURST", "100")),
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
