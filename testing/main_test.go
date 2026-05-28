package auth_test

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	globalAdminURL string
	globalConnStr  string
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	pgContainer, err := postgres.Run(
		ctx,
		"postgres:18-alpine",
		postgres.WithDatabase("testauthdb"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(15*time.Second),
		),
	)
	if err != nil {
		os.Exit(1)
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		os.Exit(1)
	}
	globalConnStr = connStr
	globalAdminURL = strings.Replace(connStr, "/testauthdb", "/postgres", 1)

	code := m.Run()

	if err := pgContainer.Terminate(ctx); err != nil {
		os.Exit(1)
	}
	os.Exit(code)
}
