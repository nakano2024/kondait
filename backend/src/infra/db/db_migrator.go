package db

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"kondait-backend/infra/config"
)

type DbMigrator struct{}

func NewDbMigrator() *DbMigrator {
	return &DbMigrator{}
}

func (dbMigrator *DbMigrator) Migrate(cfg config.Config) error {
	// migrations path
	migrationsPath := os.Getenv("DB_MIGRATIONS_PATH")
	if migrationsPath == "" {
		return fmt.Errorf("DB_MIGRATIONS_PATH is not set.")
	}
	absMigrationsPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		return fmt.Errorf("filepath.Abs: %w", err)
	}
	sourceURL := "file://" + absMigrationsPath

	// sslmode default
	sslMode := cfg.DBSSLMode
	if sslMode == "" {
		sslMode = "disable"
	}

	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.DBUser, cfg.DBPassword),
		Host:   fmt.Sprintf("%s:%s", cfg.DBHost, cfg.DBPort),
		Path:   cfg.DBName,
	}
	q := u.Query()
	q.Set("sslmode", sslMode)
	u.RawQuery = q.Encode()
	dbURL := u.String()

	m, err := migrate.New(sourceURL, dbURL)
	if err != nil {
		return fmt.Errorf("migrate.New: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migrate.Up: %w", err)
	}

	return nil
}
