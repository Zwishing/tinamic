package database

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"regexp"
	"strconv"
	"strings"

	// Logging
	log "github.com/sirupsen/logrus"
)

// DbConnect 连接数据库
func DbConnect(config *pgxpool.Config) (Db *pgxpool.Pool, err error) {
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
//QueryVersion 获取postgis版本
func QueryVersion(db *pgxpool.Pool) (map[string]string, int, error) {

	row := db.QueryRow(context.Background(), "SELECT postgis_full_version()")
	var verStr string
	err := row.Scan(&verStr)
	if err != nil {
		return nil, 0, err
	}
	// Parse full version string
	//   POSTGIS="3.0.0 r17983" [EXTENSION] PGSQL="110" GEOS="3.8.0-CAPI-1.11.0 "
	//   PROJ="6.2.0" LIBXML="2.9.4" LIBJSON="0.13" LIBPROTOBUF="1.3.2" WAGYU="0.4.3 (Internal)"
	re := regexp.MustCompile(`([A-Z]+)="(.+?)"`)
	vers := make(map[string]string)
	for _, mtch := range re.FindAllStringSubmatch(verStr, -1) {
		vers[mtch[1]] = mtch[2]
	}

	pgisVer, ok := vers["POSTGIS"]
	if !ok {
		return nil, 0, errors.New("POSTGIS key missing from postgis_full_version")
	}
	// Convert Postgis version string into a lexically (and/or numerically) sortable form
	// "3.1.1 r17983" => "3001001"
	pgisMajMinPat := strings.Split(strings.Split(pgisVer, " ")[0], ".")
	pgisMaj, _ := strconv.Atoi(pgisMajMinPat[0])
	pgisMin, _ := strconv.Atoi(pgisMajMinPat[1])
	pgisPat, _ := strconv.Atoi(pgisMajMinPat[2])
	pgisNum := 1000000*pgisMaj + 1000*pgisMin + pgisPat
	vers["POSTGISFULL"] = strconv.Itoa(pgisNum)
	//globalVersions = vers
	globalPostGISVersion := pgisNum

	return vers, globalPostGISVersion, nil
}
