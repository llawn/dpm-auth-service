package auth

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DBName      string
	DBUser      string
	DBHost      string
	DBPort      string
	DBPassword  string
	DBURL       string
	AdminURL    string
}

func Load(path string) (*Config, error) {
	_ = godotenv.Load(path)

	cfg := &Config{
		DBName:     strings.TrimSpace(os.Getenv("DB_NAME")),
		DBUser:     strings.TrimSpace(os.Getenv("DB_USER")),
		DBHost:     strings.TrimSpace(os.Getenv("DB_HOST")),
		DBPort:     strings.TrimSpace(os.Getenv("DB_PORT")),
		DBPassword: strings.TrimSpace(os.Getenv("DB_PASSWORD")),
	}

	if cfg.DBHost == "" {
		cfg.DBHost = "localhost"
	}
	if cfg.DBPort == "" {
		cfg.DBPort = "5432"
	}
	if cfg.DBName == "" {
		return nil, fmt.Errorf("DB_NAME: %w", ErrValidation)
	}
	if cfg.DBUser == "" {
		return nil, fmt.Errorf("DB_USER: %w", ErrValidation)
	}
	var creds string
	if cfg.DBPassword == "" {
		creds = cfg.DBUser
	} else {
		creds = fmt.Sprintf("%s:%s", cfg.DBUser, cfg.DBPassword)
	}

	cfg.DBURL = fmt.Sprintf(
		"postgres://%s@%s:%s/%s?sslmode=disable",
		creds,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	cfg.AdminURL = fmt.Sprintf(
		"postgres://%s@%s:%s/postgres?sslmode=disable",
		creds,
		cfg.DBHost,
		cfg.DBPort,
	)

	return cfg, nil
}
