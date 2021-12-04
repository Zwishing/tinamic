package database

import (
	"context"
	"github.com/spf13/viper"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/logrusadapter"
	"github.com/jackc/pgx/v4/pgxpool"

	// Logging
	log "github.com/sirupsen/logrus"

	."tinamic/utils"
)
func init(){
	LoadConfig()
}

func DbConnect() (*pgxpool.Pool, error) {
	var err error
	var config *pgxpool.Config
	var Db *pgxpool.Pool

	dbConnection := viper.GetString("DbConnection")
	config, err = pgxpool.ParseConfig(dbConnection)
	if err != nil {
		log.Fatal(err)
	}

	// Read and parse connection lifetime
	dbPoolMaxLifeTime, errt := time.ParseDuration(viper.GetString("DbPoolMaxConnLifeTime"))
	if errt != nil {
		log.Fatal(errt)
	}
	config.MaxConnLifetime = dbPoolMaxLifeTime

	// Read and parse max connections
	dbPoolMaxConns := viper.GetInt32("DbPoolMaxConns")
	if dbPoolMaxConns > 0 {
		config.MaxConns = dbPoolMaxConns
	}

	// Read current log level and use one less-fine level
	// below that
	config.ConnConfig.Logger = logrusadapter.NewLogger(log.New())
	levelString, _ := (log.GetLevel() - 1).MarshalText()
	pgxLevel, _ := pgx.LogLevelFromString(string(levelString))
	config.ConnConfig.LogLevel = pgxLevel

	// Connect!
	Db, err = pgxpool.ConnectConfig(context.Background(), config)

	if err != nil {
		log.Fatal(err)
	}
	dbName := config.ConnConfig.Config.Database
	dbUser := config.ConnConfig.Config.User
	dbHost := config.ConnConfig.Config.Host
	log.Infof("Connected as '%s' to '%s' @ '%s'", dbUser, dbName, dbHost)

	return Db, nil
}

//func DBTileRequest(ctx context.Context, tr *TileRequest) ([]byte, error) {
//	db, err := DbConnect()
//	if err != nil {
//		log.Error(err)
//		return nil, err
//	}
//	row := db.QueryRow(ctx, tr.SQL, tr.Args...)
//	var mvtTile []byte
//	err = row.Scan(&mvtTile)
//	if err != nil {
//		log.Warn(err)
//
//		// check for errors retrieving the rendered tile from the database
//		// Timeout errors can occur if the context deadline is reached
//		// or if the context is canceled during/before a database query.
//		if pgconn.Timeout(err) {
//			return nil,err
//		}
//
//		return nil,err
//	}
//	return mvtTile, nil
//}

