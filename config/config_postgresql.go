package config

import pg "github.com/jackc/pgx/v5/pgxpool"

type PgConfig struct {
	*pg.Config
}
