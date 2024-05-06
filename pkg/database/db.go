package database

import (
	"context"
	"fmt"
	"reflect"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"regexp"
	"strconv"
	"strings"

	// Logging
	"github.com/rs/zerolog/log"
)

var (
	Db *pgxpool.Pool //定义一个连接池
)

// DbConnect 连接数据库
func DbConnect(config *pgxpool.Config) (err error) {
	// Connect!
	Db, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatal().Msgf("%s", err)
	}
	dbName := config.ConnConfig.Config.Database
	dbUser := config.ConnConfig.Config.User
	dbHost := config.ConnConfig.Config.Host
	log.Info().Msgf("Connected as '%s' to '%s' @ '%s'", dbUser, dbName, dbHost)
	return nil
}

//	func DBTileRequest(ctx context.Context, tr *TileRequest) ([]byte, error) {
//		db, err := DbConnect()
//		if err != nil {
//			log.Error(err)
//			return nil, err
//		}
//		row := db.QueryRow(ctx, tr.SQL, tr.Args...)
//		var mvtTile []byte
//		err = row.Scan(&mvtTile)
//		if err != nil {
//			log.Warn(err)
//
//			// check for errors retrieving the rendered tile from the database
//			// Timeout errors can occur if the context deadline is reached
//			// or if the context is canceled during/before a database query.
//			if pgconn.Timeout(err) {
//				return nil,err
//			}
//
//			return nil,err
//		}
//		return mvtTile, nil
//	}
//
// QueryVersion 获取postgis版本
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
	PostGISVersion := pgisNum

	return vers, PostGISVersion, nil
}

func Raster2pgsql() {
}

//func Insert(schema, table string, source interface{}) (pgconn.CommandTag, error) {
//	t := reflect.TypeOf(source)
//	var field strings.Builder
//	var value strings.Builder
//	var values []interface{}
//	if t.Kind() == reflect.Struct {
//		for index := 0; index < t.NumField(); index++ {
//			values = append(values, reflect.ValueOf(source).Field(index).Interface())
//			if index == 0 {
//				field.WriteString("(")
//				field.WriteString(t.Field(index).Tag.Get("json"))
//				field.WriteString(",")
//
//				value.WriteString("($")
//				value.WriteString(strconv.Itoa(index + 1))
//				value.WriteString(",")
//				continue
//			}
//			if index == t.NumField()-1 {
//				field.WriteString(t.Field(index).Tag.Get("json"))
//				field.WriteString(")")
//
//				value.WriteString("$")
//				value.WriteString(strconv.Itoa(index + 1))
//				value.WriteString(")")
//				break
//			}
//			field.WriteString(t.Field(index).Tag.Get("json"))
//			field.WriteString(",")
//
//			value.WriteString("$")
//			value.WriteString(strconv.Itoa(index + 1))
//			value.WriteString(",")
//		}
//
//		var sql bytes.Buffer
//		sql.WriteString("INSERT INTO ")
//		sql.WriteString(schema)
//		sql.WriteString(".")
//		sql.WriteString(table)
//		sql.WriteString(field.String())
//		sql.WriteString("VALUES")
//		sql.WriteString(value.String())
//		sqlString := sql.String()
//
//		log.Info().Msg(sqlString)
//
//		tag, err := Db.Exec(context.Background(), sqlString, values...)
//
//		if err != nil {
//			return nil, err
//		}
//		return tag, nil
//	}
//	return nil, errors.New("")
//}

//func Select(sql string, dest interface{}) {
//	pgxscan.Select(context.Background(), Db, &dest, sql)
//}

func Create(schema, table string, source interface{}) {
	t := reflect.TypeOf(source)
	if t.Kind() != reflect.Struct {
		return
	}
	createSchemaSql := fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS %s`, schema)
	fmt.Println(createSchemaSql)
}
