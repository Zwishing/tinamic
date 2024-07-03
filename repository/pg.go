package repository

import (
	"fmt"
	"sync"
	"tinamic/conf"
	"tinamic/pkg/pg"
)

var (
	db *pg.PGPool //定义一个连接池

	dbOnce sync.Once
)

func GetDbPoolInstance() *pg.PGPool {
	dbOnce.Do(func() {
		cfg := conf.GetConfigInstance()
		constr := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s",
			cfg.GetString("database.postgresql.user"),
			cfg.GetString("database.postgresql.password"),
			cfg.GetString("database.postgresql.host"),
			cfg.GetInt32("database.postgresql.port"),
			cfg.GetString("database.postgresql.database"),
			cfg.GetString("database.postgresql.sslmode"))
		dbConfig := pg.NewPgConfig(pg.WithConnString(constr))
		db = pg.New(dbConfig.Config)
	})
	return db
}
