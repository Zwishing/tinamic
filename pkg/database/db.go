package database

import (
	"bytes"
	"context"
	"fmt"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"tinamic/conf"

	// Logging
	"github.com/rs/zerolog/log"
)

var (
	db *Dbpool //定义一个连接池

	once sync.Once
)

type Dbpool struct {
	*pgxpool.Pool
}

func New(config *pgxpool.Config) *Dbpool {
	dbPool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Error().Msgf("%s", err)
		return nil
	}
	dbName := config.ConnConfig.Config.Database
	dbUser := config.ConnConfig.Config.User
	dbHost := config.ConnConfig.Config.Host
	log.Info().Msgf("Connected as '%s' to '%s' @ '%s'", dbUser, dbName, dbHost)
	return &Dbpool{
		dbPool,
	}
}

// QueryVersion 查询版本
func (db *Dbpool) QueryVersion() (map[string]string, int, error) {
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

// Insert 插入数据
func (db *Dbpool) Insert(schema, table string, source interface{}) (pgconn.CommandTag, error) {
	t := reflect.TypeOf(source)
	var field strings.Builder
	var value strings.Builder
	var values []interface{}
	if t.Kind() == reflect.Struct {
		for index := 0; index < t.NumField(); index++ {
			values = append(values, reflect.ValueOf(source).Field(index).Interface())
			if index == 0 {
				field.WriteString("(")
				field.WriteString(t.Field(index).Tag.Get("json"))
				field.WriteString(",")

				value.WriteString("($")
				value.WriteString(strconv.Itoa(index + 1))
				value.WriteString(",")
				continue
			}
			if index == t.NumField()-1 {
				field.WriteString(t.Field(index).Tag.Get("json"))
				field.WriteString(")")

				value.WriteString("$")
				value.WriteString(strconv.Itoa(index + 1))
				value.WriteString(")")
				break
			}
			field.WriteString(t.Field(index).Tag.Get("json"))
			field.WriteString(",")

			value.WriteString("$")
			value.WriteString(strconv.Itoa(index + 1))
			value.WriteString(",")
		}

		var sql bytes.Buffer
		sql.WriteString("INSERT INTO ")
		sql.WriteString(schema)
		sql.WriteString(".")
		sql.WriteString(table)
		sql.WriteString(field.String())
		sql.WriteString("VALUES")
		sql.WriteString(value.String())
		sqlString := sql.String()

		log.Info().Msg(sqlString)

		tag, err := db.Exec(context.Background(), sqlString, values...)

		if err != nil {
			return nil, err
		}
		return tag, nil
	}
	return nil, errors.New("")
}

func (db *Dbpool) Select(sql string, dest interface{}) error {
	err := pgxscan.Select(context.Background(), db, &dest, sql)
	if err != nil {
		return err
	}
	return nil
}

// SelectRow 查询单行数据
func (db *Dbpool) SelectRow(sql string, dest ...interface{}) error {
	err := db.QueryRow(context.Background(), sql).Scan(dest...)
	if err != nil {
		return err
	}
	return nil
}

func (db *Dbpool) Create(schema, table string, source interface{}) {
	t := reflect.TypeOf(source)
	if t.Kind() != reflect.Struct {
		return
	}
	createSchemaSql := fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS %s`, schema)
	log.Info().Msgf(createSchemaSql)
}

func GetDbPoolInstance() *Dbpool {
	once.Do(func() {
		cfg := conf.GetConfigInstance()
		constr := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s",
			cfg.GetString("database.postgresql.user"),
			cfg.GetString("database.postgresql.password"),
			cfg.GetString("database.postgresql.host"),
			cfg.GetInt32("database.postgresql.port"),
			cfg.GetString("database.postgresql.database"),
			cfg.GetString("database.postgresql.sslmode"))
		dbConfig := NewPgConfig(WithConnString(constr))
		db = New(dbConfig.Config)
	})
	return db
}
