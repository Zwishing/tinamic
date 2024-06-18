package cache

import (
	"crypto/tls"
)

// Config defines the config for storage.
type RedisConfig struct {
	// Host name where the DB is hosted
	//
	// Optional. Default is "127.0.0.1"
	Host string

	// Port where the DB is listening on
	//
	// Optional. Default is 6379
	Port int

	// Server username
	//
	// Optional. Default is ""
	Username string

	// Server password
	//
	// Optional. Default is ""
	Password string

	// Database to be selected after connecting to the server.
	//
	// Optional. Default is 0
	Database int

	// URL standard format Redis URL. If this is set all other config options, Host, Port, Username, Password, Database have no effect.
	//
	// Example: redis://<user>:<pass>@localhost:6379/<db>
	// Optional. Default is ""
	URL string

	// Either a single address or a seed list of host:port addresses, this enables FailoverClient and ClusterClient
	//
	// Optional. Default is []string{}
	Addrs []string

	// MasterName is the sentinel master's name
	//
	// Optional. Default is ""
	MasterName string

	// ClientName will execute the `CLIENT SETNAME ClientName` command for each conn.
	//
	// Optional. Default is ""
	ClientName string

	// SentinelUsername
	//
	// Optional. Default is ""
	SentinelUsername string

	// SentinelPassword
	//
	// Optional. Default is ""
	SentinelPassword string

	// Reset clears any existing keys in existing Collection
	//
	// Optional. Default is false
	Reset bool

	// TLS Config to use. When set TLS will be negotiated.
	//
	// Optional. Default is nil
	TLSConfig *tls.Config

	// Maximum number of socket connections.
	//
	// Optional. Default is 10 connections per every available CPU as reported by runtime.GOMAXPROCS.
	PoolSize int
}
