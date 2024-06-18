package database

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
	"time"
)

type PgConfig struct {
	*pgxpool.Config
}

type PgOptions func(cfg *PgConfig)

// NewPgConfig creates a new PgConfig with the provided options.
func NewPgConfig(opts ...PgOptions) *PgConfig {
	cfg := &PgConfig{
		&pgxpool.Config{},
	}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

// Option functions

func WithConnString(connString string) PgOptions {
	return func(cfg *PgConfig) {
		config, err := pgxpool.ParseConfig(connString)
		if err != nil {
			log.Error().Msgf("Invalid connection string: %s", err.Error())
		}
		cfg.Config = config
	}
}

func WithMaxConns(maxConns int32) PgOptions {
	return func(cfg *PgConfig) {
		cfg.MaxConns = maxConns
	}
}

func WithMinConns(minConns int32) PgOptions {
	return func(cfg *PgConfig) {
		cfg.MinConns = minConns
	}
}

func WithMaxConnLifetime(maxConnLifetime int64) PgOptions {
	return func(cfg *PgConfig) {
		cfg.MaxConnLifetime = time.Duration(maxConnLifetime)
	}
}
