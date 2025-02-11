package config

import (
	"os"
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
	RedisHost        string
	RedisPort        string
	RedisPassword    string
	RedisDB          string
}

var AppConfig Config

func Load() {
	AppConfig = Config{
		AppPort:          os.Getenv("APP_PORT"),
		JWTSecret:        os.Getenv("JWT_SECRET"),
		JWTAccessExpiry:  os.Getenv("JWT_ACCESS_EXPIRY"),
		JWTRefreshExpiry: os.Getenv("JWT_REFRESH_EXPIRY"),
		MySQLHost:        os.Getenv("MYSQL_HOST"),
		MySQLPort:        os.Getenv("MYSQL_PORT"),
		MySQLDatabase:    os.Getenv("MYSQL_DATABASE"),
		MySQLUsername:    os.Getenv("MYSQL_USERNAME"),
		MySQLPassword:    os.Getenv("MYSQL_PASSWORD"),
		RedisHost:        os.Getenv("REDIS_HOST"),
		RedisPort:        os.Getenv("REDIS_PORT"),
		RedisPassword:    os.Getenv("REDIS_PASSWORD"),
		RedisDB:          os.Getenv("REDIS_DB"),
	}
}
