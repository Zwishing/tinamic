package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *Storage
)

// Storage interface that is implemented by storage providers
type Storage struct {
	db redis.UniversalClient
}

// New creates a new redis storage
func New(cfg *RedisConfig) *Storage {

	// Create new redis universal client
	var db redis.UniversalClient

	// Parse the URL and update config values accordingly
	if cfg.URL != "" {
		options, err := redis.ParseURL(cfg.URL)
		if err != nil {
			log.Error().Msgf(err.Error())
		}

		// Update the config values with the parsed URL values
		cfg.Username = options.Username
		cfg.Password = options.Password
		cfg.Database = options.DB
		cfg.Addrs = []string{options.Addr}

		// If cfg.TLSConfig is not provided, and options returns one, use it.
		if cfg.TLSConfig == nil && options.TLSConfig != nil {
			cfg.TLSConfig = options.TLSConfig
		}
	} else if len(cfg.Addrs) == 0 {
		// Fallback to Host and Port values if Addrs is empty
		cfg.Addrs = []string{fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)}
	}

	// Create Universal Client
	db = redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:            cfg.Addrs,
		MasterName:       cfg.MasterName,
		ClientName:       cfg.ClientName,
		SentinelUsername: cfg.SentinelUsername,
		SentinelPassword: cfg.SentinelPassword,
		DB:               cfg.Database,
		Username:         cfg.Username,
		Password:         cfg.Password,
		TLSConfig:        cfg.TLSConfig,
		PoolSize:         cfg.PoolSize,
	})

	// Test connection
	if err := db.Ping(context.Background()).Err(); err != nil {
		log.Error().Msgf(err.Error())
	}

	// Empty collection if Clear is true
	if cfg.Reset {
		if err := db.FlushDB(context.Background()).Err(); err != nil {
			log.Error().Msgf(err.Error())
		}
	}
	log.Info().Msgf("Connected to redis @ '%s'", cfg.Host)
	// Create new store
	return &Storage{
		db: db,
	}
}

// Get value by key
func (s *Storage) Get(key string) (string, error) {
	if len(key) <= 0 {
		return "", nil
	}
	val, err := s.db.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

// Set key with value
func (s *Storage) Set(key string, val []byte, exp time.Duration) error {
	if len(key) <= 0 || len(val) <= 0 {
		return nil
	}
	return s.db.Set(context.Background(), key, val, exp).Err()
}

func (s *Storage) SetMap(key string, val map[string]string, exp time.Duration) error {
	if len(key) <= 0 || len(val) <= 0 {
		return nil
	}

	marshalProfile, err := json.Marshal(val)
	if err != nil {
		return err
	}
	return s.db.Set(context.Background(), key, marshalProfile, exp).Err()
}

// Delete key by key
func (s *Storage) Delete(key string) error {
	if len(key) <= 0 {
		return nil
	}
	return s.db.Del(context.Background(), key).Err()
}

// Reset all keys
func (s *Storage) Reset() error {
	return s.db.FlushDB(context.Background()).Err()
}

// Close the database
func (s *Storage) Close() error {
	return s.db.Close()
}

// Return database client
func (s *Storage) Conn() redis.UniversalClient {
	return s.db
}

// Return all the keys
func (s *Storage) Keys() ([][]byte, error) {
	var keys [][]byte
	var cursor uint64
	var err error

	for {
		var batch []string

		if batch, cursor, err = s.db.Scan(context.Background(), cursor, "*", 10).Result(); err != nil {
			return nil, err
		}

		for _, key := range batch {
			keys = append(keys, []byte(key))
		}

		if cursor == 0 {
			break
		}
	}

	if len(keys) == 0 {
		return nil, nil
	}

	return keys, nil
}
