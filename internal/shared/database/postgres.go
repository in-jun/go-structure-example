package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/in-jun/go-structure-example/internal/shared/config"
)

var ErrMigrateOnly = errors.New("migrate-only: migrations complete")

func NewPostgres() (*sql.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config.AppConfig.PGUsername,
		config.AppConfig.PGPassword,
		config.AppConfig.PGHost,
		config.AppConfig.PGPort,
		config.AppConfig.PGDatabase,
		config.AppConfig.PGSSLMode,
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if config.AppConfig.PGMaxOpenConns > 0 {
		db.SetMaxOpenConns(config.AppConfig.PGMaxOpenConns)
	}
	if config.AppConfig.PGMaxIdleConns > 0 {
		db.SetMaxIdleConns(config.AppConfig.PGMaxIdleConns)
	}
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(3 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	switch config.AppConfig.MigrationMode {
	case "skip":
		slog.Info("Migration mode: skip — skipping migrations")
	case "migrate-only":
		if err := runMigrations(db); err != nil {
			return nil, fmt.Errorf("migration failed: %w", err)
		}
		slog.Info("Migrations complete (migrate-only mode)")
		return db, ErrMigrateOnly
	default:
		if err := runMigrations(db); err != nil {
			return nil, fmt.Errorf("migration failed: %w", err)
		}
	}

	return db, nil
}

func runMigrations(db *sql.DB) error {
	if _, err := db.Exec("SELECT pg_advisory_lock(1)"); err != nil {
		return fmt.Errorf("failed to acquire migration lock: %w", err)
	}
	defer func() {
		if _, err := db.Exec("SELECT pg_advisory_unlock(1)"); err != nil {
			slog.Warn("failed to release migration advisory lock", "error", err)
		}
	}()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+config.AppConfig.MigrationPath,
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	version, dirty, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return fmt.Errorf("failed to get migration version: %w", err)
	}

	if dirty {
		return fmt.Errorf("dirty migration at version %d — manual intervention required", version)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	slog.Info("migrations applied", "version", version)
	return nil
}
