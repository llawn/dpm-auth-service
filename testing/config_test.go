package auth_test

import (
	"os"
	"path/filepath"
	"testing"

	"dpm/services/auth"

	"github.com/google/go-cmp/cmp"
)

func setEnv(t *testing.T, env map[string]string) {
	t.Helper()
	for k, v := range env {
		t.Setenv(k, v)
	}
}

func unsetEnv(t *testing.T, keys ...string) {
	t.Helper()
	for _, k := range keys {
		_ = os.Unsetenv(k)
	}
}

func TestLoad_Defaults(t *testing.T) {
	unsetEnv(t, "DB_NAME", "DB_USER", "DB_HOST", "DB_PORT", "DB_PASSWORD")

	setEnv(t, map[string]string{
		"DB_NAME": "testdb",
		"DB_USER": "postgres",
	})

	cfg, err := auth.Load("")
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	want := &auth.Config{
		DBName:     "testdb",
		DBUser:     "postgres",
		DBHost:     "localhost",
		DBPort:     "5432",
		DBPassword: "",
		DBURL:      "postgres://postgres@localhost:5432/testdb?sslmode=disable",
		AdminURL:   "postgres://postgres@localhost:5432/postgres?sslmode=disable",
	}
	if diff := cmp.Diff(want, cfg); diff != "" {
		t.Errorf("Load() mismatch (-want +got):\n%s", diff)
	}
}

func TestLoad_TrimsWhitespace(t *testing.T) {
	unsetEnv(t, "DB_NAME", "DB_USER", "DB_HOST", "DB_PORT", "DB_PASSWORD")

	setEnv(t, map[string]string{
		"DB_NAME":     "  testdb  ",
		"DB_USER":     "  admin  ",
		"DB_HOST":     "  192.168.1.1  ",
		"DB_PORT":     "  6432  ",
		"DB_PASSWORD": "  s3cret  ",
	})

	cfg, err := auth.Load("")
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	want := &auth.Config{
		DBName:     "testdb",
		DBUser:     "admin",
		DBHost:     "192.168.1.1",
		DBPort:     "6432",
		DBPassword: "s3cret",
		DBURL:      "postgres://admin:s3cret@192.168.1.1:6432/testdb?sslmode=disable",
		AdminURL:   "postgres://admin:s3cret@192.168.1.1:6432/postgres?sslmode=disable",
	}
	if diff := cmp.Diff(want, cfg); diff != "" {
		t.Errorf("Load() mismatch (-want +got):\n%s", diff)
	}
}

func TestLoad_FromDotenvFile(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", "..", ".env.example"))
	if err != nil {
		t.Fatalf("failed to read .env.example: %v", err)
	}

	dir := t.TempDir()
	p := filepath.Join(dir, ".env.test")

	err = os.WriteFile(p, data, 0644)
	if err != nil {
		t.Fatalf("failed to write .env file: %v", err)
	}

	unsetEnv(t, "DB_NAME", "DB_USER", "DB_HOST", "DB_PORT", "DB_PASSWORD")

	cfg, err := auth.Load(p)
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	want := &auth.Config{
		DBName:     "dpm",
		DBUser:     "postgres",
		DBHost:     "localhost",
		DBPort:     "5432",
		DBPassword: "secret",
		DBURL:      "postgres://postgres:secret@localhost:5432/dpm?sslmode=disable",
		AdminURL:   "postgres://postgres:secret@localhost:5432/postgres?sslmode=disable",
	}
	if diff := cmp.Diff(want, cfg); diff != "" {
		t.Errorf("Load() mismatch (-want +got):\n%s", diff)
	}
}
