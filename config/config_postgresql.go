package config

import pg "github.com/jackc/pgx/v4/pgxpool"

type PgConfig struct {
	*pg.Config
}
