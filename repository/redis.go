package repository

import (
	"sync"
	"tinamic/conf"
	. "tinamic/pkg/redis"
)

var (
	RedisClient *RedisStorage
	redisOnce   sync.Once
)

func GetRedisInstance() *RedisStorage {
	redisOnce.Do(func() {
		cfg := conf.GetConfigInstance()
		prefix := "database.redis."
		redisConfig := NewRedisConfig(
			WithHost(cfg.GetString(prefix+"host")),
			WithPort(cfg.GetInt(prefix+"port")),
			WithDatabase(cfg.GetInt(prefix+"database")),
			WithPassword(cfg.GetString(prefix+"password")),
			WithPoolSize(cfg.GetInt(prefix+"poolMaxConns")),
		)
		RedisClient = NewRedisStorage(redisConfig)
	})
	return RedisClient
}
