package auth_test

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strings"
	"testing"

	"dpm/services/auth"
	"dpm/services/auth/migrations"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
)

func randomDBName() string {
	b := make([]byte, 4)
	_, _ = rand.Read(b)
	return "test_" + hex.EncodeToString(b)
}

func initAuthDB(ctx context.Context, t *testing.T, connStr string) *auth.DB {
	t.Helper()
	authDB, err := auth.NewDB(ctx, connStr)
	if err != nil {
		t.Fatalf("failed to connect to test database: %s", err)
	}
	return authDB
}

func createSchema(t *testing.T, connStr string) *migrate.Migrate {
	t.Helper()
	d, err := iofs.New(migrations.FS, ".")
	if err != nil {
		t.Fatalf("failed to create iofs source: %s", err)
	}
	m, err := migrate.NewWithSourceInstance("iofs", d, connStr)
	if err != nil {
		t.Fatalf("failed to create migrate instance: %s", err)
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		t.Fatalf("failed to apply migrations: %s", err)
	}
	return m
}

func setupTestDB(ctx context.Context, t *testing.T) (*auth.DB, func()) {
	t.Helper()

	dbName := randomDBName()

	pool, err := pgxpool.New(ctx, globalAdminURL)
	if err != nil {
		t.Fatalf("failed to connect as admin: %s", err)
	}
	_, err = pool.Exec(ctx, "CREATE DATABASE "+dbName)
	pool.Close()
	if err != nil {
		t.Fatalf("failed to create database %s: %s", dbName, err)
	}

	connStr := strings.Replace(globalAdminURL, "/postgres?", "/"+dbName+"?", 1)

	m := createSchema(t, connStr)

	authDB := initAuthDB(ctx, t, connStr)

	cleanup := func() {
		if m != nil {
			err := m.Down()
			if err != nil && err != migrate.ErrNoChange {
				t.Fatalf("migrate down failed: %v", err)
			}
			_, _ = m.Close()
		}

		if authDB != nil && authDB.Pool() != nil {
			authDB.Pool().Close()
		}

		cleanPool, err := pgxpool.New(ctx, globalAdminURL)
		if err != nil {
			t.Fatalf("failed to connect as admin for cleanup: %s", err)
		}
		defer cleanPool.Close()

		_, err = cleanPool.Exec(ctx,
			"SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = $1 AND pid <> pg_backend_pid()",
			dbName,
		)
		if err != nil {
			t.Fatalf("failed to terminate connections to %s: %s", dbName, err)
		}

		_, err = cleanPool.Exec(ctx, "DROP DATABASE IF EXISTS "+dbName)
		if err != nil {
			t.Fatalf("failed to drop database %s: %s", dbName, err)
		}
	}

	return authDB, cleanup
}
