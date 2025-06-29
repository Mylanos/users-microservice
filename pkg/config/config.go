package config

import (
	"os"
	"time"
	"users-backend/pkg/models"
)

type Config struct {
	DatabaseURL     string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

func Load() (*Config, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, models.ErrDatabaseEnvConfigNotSet
	}

	cfg := &Config{
		DatabaseURL:     dbURL,
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 1 * time.Minute,
	}

	return cfg, nil
}
