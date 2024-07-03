package pg

import (
	"context"
	"fmt"
	//"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	// Logging
	"github.com/rs/zerolog/log"
)

type PGPool struct {
	*pgxpool.Pool
}

func New(config *pgxpool.Config) *PGPool {
	//dbPool, err := pgxpool.ConnectConfig(context.Background(), config)
	dbPool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Error().Msgf("%s", err)
		return nil
	}
	dbName := config.ConnConfig.Config.Database
	dbUser := config.ConnConfig.Config.User
	dbHost := config.ConnConfig.Config.Host
	log.Info().Msgf("Connected as '%s' to '%s' @ '%s'", dbUser, dbName, dbHost)
	return &PGPool{
		dbPool,
	}
}

func (db *PGPool) ConnString() string {
	return db.Config().ConnConfig.ConnString()
}

// QueryVersion 查询版本
func (db *PGPool) QueryVersion() (map[string]string, int, error) {
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

func CreateInsertSql(schema, table string, source any) (string, []any, error) {
	t := reflect.TypeOf(source)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return "", nil, fmt.Errorf("source is not a struct or struct pointer")
	}
	v := reflect.ValueOf(source)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	columns := make([]string, 0, t.NumField())
	values := make([]interface{}, 0, t.NumField())
	placeholders := make([]string, 0, t.NumField())
	placeholderIndex := 1 // 占位符索引从1开始
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		columnName := field.Tag.Get("json")
		value := v.Field(i).Interface()

		// 忽略omitempty字段和未指定列名的字段
		if columnName == "" || columnName == "-" {
			s := reflect.TypeOf(value)
			val := reflect.ValueOf(value)
			for j := 0; j < s.NumField(); j++ {
				f := s.Field(j)
				c := f.Tag.Get("json")
				va := val.Field(j).Interface()
				columns = append(columns, c)
				values = append(values, va)
				placeholders = append(placeholders, fmt.Sprintf("$%d", placeholderIndex))
				placeholderIndex++
			}
			continue
		}

		columns = append(columns, columnName)
		values = append(values, value)
		placeholders = append(placeholders, fmt.Sprintf("$%d", placeholderIndex))
		placeholderIndex++
	}

	// 构建 SQL 语句
	query := fmt.Sprintf("INSERT INTO %s.%s (%s) VALUES (%s)",
		schema, table,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "))

	return query, values, nil
}

// Insert 插入数据
func (db *PGPool) Insert(schema, table string, source interface{}) error {

	query, values, err := CreateInsertSql(schema, table, source)
	// 输出 SQL 语句，用于调试
	fmt.Println("Generated SQL:", query)

	// 执行数据库插入操作
	_, err = db.Exec(context.Background(), query, values...)
	if err != nil {
		return err
	}

	return nil
}

// SelectRow 查询单行数据
func (db *PGPool) SelectRow(sql string, dest ...interface{}) error {
	err := db.QueryRow(context.Background(), sql).Scan(dest...)
	if err != nil {
		return err
	}
	return nil
}

func (db *PGPool) Create(schema, table string, source interface{}) {
	t := reflect.TypeOf(source)
	if t.Kind() != reflect.Struct {
		return
	}
	createSchemaSql := fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS %s`, schema)
	log.Info().Msgf(createSchemaSql)
}
