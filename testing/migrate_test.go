package auth_test

import (
	"context"
	"testing"

	"dpm/services/auth"
)

func TestCreateDB(t *testing.T) {
	ctx := context.Background()

	cfg := &auth.Config{
		AdminURL: globalAdminURL,
		DBName:   "testcreatedb",
	}

	t.Run("creates database when it does not exist", func(t *testing.T) {
		err := auth.CreateDB(ctx, cfg)
		if err != nil {
			t.Fatalf("CreateDB() returned error: %v", err)
		}
	})

	t.Run("no error when database already exists", func(t *testing.T) {
		err := auth.CreateDB(ctx, cfg)
		if err != nil {
			t.Fatalf("CreateDB() should not error when DB exists: %v", err)
		}
	})
}

func TestMigrateUp(t *testing.T) {
	t.Run("applies migrations successfully", func(t *testing.T) {
		err := auth.MigrateUp(globalConnStr)
		if err != nil {
			t.Fatalf("MigrateUp() returned error: %v", err)
		}
	})

	t.Run("no error when no new migrations", func(t *testing.T) {
		err := auth.MigrateUp(globalConnStr)
		if err != nil {
			t.Fatalf("MigrateUp() should not error on re-run: %v", err)
		}
	})
}
