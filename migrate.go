package auth

import (
	"context"
	"fmt"

	"dpm/services/auth/migrations"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateDB(ctx context.Context, cfg *Config) error {
	adminURL := cfg.AdminURL

	pool, err := pgxpool.New(ctx, adminURL)
	if err != nil {
		return fmt.Errorf("connect to postgres: %w", err)
	}
	defer pool.Close()

	var exists bool
	err = pool.QueryRow(
		ctx,
		"SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)",
		cfg.DBName,
	).Scan(&exists)
	if err != nil {
		return fmt.Errorf("check database existence: %w", err)
	}

	if exists {
		return nil
	}

	_, err = pool.Exec(ctx, "CREATE DATABASE "+cfg.DBName)
	if err != nil {
		return fmt.Errorf("create database: %w", err)
	}

	return nil
}

func MigrateUp(connStr string) error {
	d, err := iofs.New(migrations.FS, ".")
	if err != nil {
		return fmt.Errorf("migrate iofs: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, connStr)
	if err != nil {
		return fmt.Errorf("migrate instance: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migrate up: %w", err)
	}

	return nil
}
