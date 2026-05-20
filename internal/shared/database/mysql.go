package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/in-jun/go-structure-example/internal/shared/config"
)

func NewMySQL() (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true",
		config.AppConfig.MySQLUsername,
		config.AppConfig.MySQLPassword,
		config.AppConfig.MySQLHost,
		config.AppConfig.MySQLPort,
		config.AppConfig.MySQLDatabase,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if config.AppConfig.MySQLMaxOpenConns > 0 {
		db.SetMaxOpenConns(config.AppConfig.MySQLMaxOpenConns)
	}
	if config.AppConfig.MySQLMaxIdleConns > 0 {
		db.SetMaxIdleConns(config.AppConfig.MySQLMaxIdleConns)
	}
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(3 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	return db, nil
}

func runMigrations(db *sql.DB) error {
	// Advisory lock prevents multiple instances from migrating simultaneously.
	var acquired int
	if err := db.QueryRow("SELECT GET_LOCK('migration_lock', 30)").Scan(&acquired); err != nil {
		return fmt.Errorf("failed to acquire migration lock: %w", err)
	}
	if acquired != 1 {
		return fmt.Errorf("could not acquire migration lock (timeout)")
	}
	defer db.Exec("SELECT RELEASE_LOCK('migration_lock')")

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return fmt.Errorf("failed to create driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+config.AppConfig.MigrationPath,
		"mysql",
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
