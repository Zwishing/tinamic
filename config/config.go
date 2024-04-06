package config

import (
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"os"

	"github.com/rs/zerolog/log"
	"strings"
	"time"

	"github.com/spf13/viper"

	hashing "github.com/thomasvvugt/fiber-hashing"
	argon_driver "github.com/thomasvvugt/fiber-hashing/driver/argon2id"
	bcrypt_driver "github.com/thomasvvugt/fiber-hashing/driver/bcrypt"

	"github.com/alexedwards/argon2id"
	"github.com/gofiber/fiber/v2"
	//"github.com/gofiber/session/v2/provider/sqlite3"
)

type Config struct {
	*viper.Viper
	fiber *fiber.Config
}

var (
	Conf *Config //配置文件
)

func New() *Config {
	config := new(Config)
	config.Viper = viper.New()
	// Set default configurations
	//config.setDefaults()

	config.AddConfigPath("./conf")
	config.AddConfigPath("../conf")
	config.SetConfigName("tinamic")
	config.SetConfigType("toml")

	// Read configuration
	if err := config.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			log.Error().Msgf("failed to read configuration:%s", err.Error())
			os.Exit(1)
		}
	}

	//config.SetErrorHandler(defaultErrorHandler)

	// TODO: Logger (Maybe a different zap object)

	// TODO: Add APP_KEY generation

	// TODO: Write changes to configuration file

	// Set Fiber configurations
	config.setFiberConfig()

	return config
}

//func (config *Config) SetErrorHandler(errorHandler fiber.ErrorHandler) {
//	config.errorHandler = errorHandler
//}

func (config *Config) setDefaults() {
	// 1d, 1h, 1m, 1s, see https://golang.org/pkg/time/#ParseDuration
	config.SetDefault("poolMaxConnLifeTime", "1h")
	config.SetDefault("poolMaxConns", 4)
	config.SetDefault("timeout", 10)

	// Set default Fiber configuration
	config.SetDefault("FIBER_PREFORK", false)
	config.SetDefault("FIBER_SERVERHEADER", "")
	config.SetDefault("FIBER_STRICTROUTING", false)
	config.SetDefault("FIBER_CASESENSITIVE", false)
	config.SetDefault("FIBER_IMMUTABLE", false)
	config.SetDefault("FIBER_UNESCAPEPATH", false)
	config.SetDefault("FIBER_ETAG", false)
	config.SetDefault("FIBER_BODYLIMIT", 209715200) //200*1024*1024
	config.SetDefault("FIBER_CONCURRENCY", 262144)
	config.SetDefault("FIBER_VIEWS", "html")
	config.SetDefault("FIBER_VIEWS_DIRECTORY", "resources/views")
	config.SetDefault("FIBER_VIEWS_RELOAD", false)
	config.SetDefault("FIBER_VIEWS_DEBUG", false)
	config.SetDefault("FIBER_VIEWS_LAYOUT", "embed")
	config.SetDefault("FIBER_VIEWS_DELIMS_L", "{{")
	config.SetDefault("FIBER_VIEWS_DELIMS_R", "}}")
	config.SetDefault("FIBER_READTIMEOUT", 0)
	config.SetDefault("FIBER_WRITETIMEOUT", 0)
	config.SetDefault("FIBER_IDLETIMEOUT", 0)
	config.SetDefault("FIBER_READBUFFERSIZE", 4096)
	config.SetDefault("FIBER_WRITEBUFFERSIZE", 4096)
	config.SetDefault("FIBER_COMPRESSEDFILESUFFIX", ".fiber.gz")
	config.SetDefault("FIBER_PROXYHEADER", "")
	config.SetDefault("FIBER_GETONLY", false)
	config.SetDefault("FIBER_DISABLEKEEPALIVE", false)
	config.SetDefault("FIBER_DISABLEDEFAULTDATE", false)
	config.SetDefault("FIBER_DISABLEDEFAULTCONTENTTYPE", false)
	config.SetDefault("FIBER_DISABLEHEADERNORMALIZING", false)
	config.SetDefault("FIBER_DISABLESTARTUPMESSAGE", false)
	config.SetDefault("FIBER_REDUCEMEMORYUSAGE", false)

	// Set default Fiber CORS middlewares configuration
	config.SetDefault("MW_FIBER_CORS_ENABLED", false)
	config.SetDefault("MW_FIBER_CORS_ALLOWORIGINS", "*")
	config.SetDefault("MW_FIBER_CORS_ALLOWMETHODS", "GET,POST,HEAD,PUT,DELETE,PATCH")
	config.SetDefault("MW_FIBER_CORS_ALLOWHEADERS", "")
	config.SetDefault("MW_FIBER_CORS_ALLOWCREDENTIALS", false)
	config.SetDefault("MW_FIBER_CORS_EXPOSEHEADERS", "")
	config.SetDefault("MW_FIBER_CORS_MAXAGE", 0)

}

func (config *Config) setFiberConfig() {
	config.fiber = &fiber.Config{
		Prefork:              config.GetBool("FIBER_PREFORK"),
		ServerHeader:         config.GetString("FIBER_SERVERHEADER"),
		StrictRouting:        config.GetBool("FIBER_STRICTROUTING"),
		CaseSensitive:        config.GetBool("FIBER_CASESENSITIVE"),
		Immutable:            config.GetBool("FIBER_IMMUTABLE"),
		UnescapePath:         config.GetBool("FIBER_UNESCAPEPATH"),
		ETag:                 config.GetBool("FIBER_ETAG"),
		BodyLimit:            config.GetInt("FIBER_BODYLIMIT"),
		Concurrency:          config.GetInt("FIBER_CONCURRENCY"),
		Views:                nil,
		ReadTimeout:          config.GetDuration("FIBER_READTIMEOUT"),
		WriteTimeout:         config.GetDuration("FIBER_WRITETIMEOUT"),
		IdleTimeout:          config.GetDuration("FIBER_IDLETIMEOUT"),
		ReadBufferSize:       config.GetInt("FIBER_READBUFFERSIZE"),
		WriteBufferSize:      config.GetInt("FIBER_WRITEBUFFERSIZE"),
		CompressedFileSuffix: config.GetString("FIBER_COMPRESSEDFILESUFFIX"),
		ProxyHeader:          config.GetString("FIBER_PROXYHEADER"),
		GETOnly:              config.GetBool("FIBER_GETONLY"),
		//ErrorHandler:              config.errorHandler,
		DisableKeepalive:          config.GetBool("FIBER_DISABLEKEEPALIVE"),
		DisableDefaultDate:        config.GetBool("FIBER_DISABLEDEFAULTDATE"),
		DisableDefaultContentType: config.GetBool("FIBER_DISABLEDEFAULTCONTENTTYPE"),
		DisableHeaderNormalizing:  config.GetBool("FIBER_DISABLEHEADERNORMALIZING"),
		DisableStartupMessage:     config.GetBool("FIBER_DISABLESTARTUPMESSAGE"),
		ReduceMemoryUsage:         config.GetBool("FIBER_REDUCEMEMORYUSAGE"),
	}
}

func (config *Config) GetFiberConfig() *fiber.Config {
	return config.fiber
}

func (config *Config) GetHasherConfig() hashing.Config {
	if strings.ToLower(config.GetString("HASHER_DRIVER")) == "bcrypt" {
		return hashing.Config{
			Driver: bcrypt_driver.New(bcrypt_driver.Config{
				Complexity: config.GetInt("HASHER_ROUNDS"),
			})}
	} else {
		return hashing.Config{
			Driver: argon_driver.New(argon_driver.Config{
				Params: &argon2id.Params{
					Memory:      config.GetUint32("HASHER_MEMORY"),
					Iterations:  config.GetUint32("HASHER_ITERATIONS"),
					Parallelism: uint8(config.GetInt("HASHER_PARALLELISM")),
					SaltLength:  config.GetUint32("HASHER_SALTLENGTH"),
					KeyLength:   config.GetUint32("HASHER_KEYLENGTH"),
				}})}
	}
}

//func (CONFIGFILE *Config) GetSessionConfig() session.Config {
//	var provider fsession.Provider
//	switch strings.ToLower(CONFIGFILE.GetString("SESSION_PROVIDER")) {
//	case "memcache":
//		sessionProvider, err := memcache.New(memcache.Config{
//			KeyPrefix:    CONFIGFILE.GetString("SESSION_KEYPREFIX"),
//			ServerList:   []string {
//				CONFIGFILE.GetString("SESSION_HOST") + ":" + CONFIGFILE.GetString("SESSION_PORT"),
//			},
//		})
//		if err != nil {
//			fmt.Println("failed to initialized memcache session provider:", err.Error())
//			break
//		}
//		provider = sessionProvider
//		break
//	case "mysql":
//		sessionProvider, err := mysql.New(mysql.Config{
//			Host:            CONFIGFILE.GetString("SESSION_HOST"),
//			Port:            CONFIGFILE.GetInt("SESSION_PORT"),
//			Username:        CONFIGFILE.GetString("SESSION_USERNAME"),
//			Password:        CONFIGFILE.GetString("SESSION_PASSWORD"),
//			Database:        CONFIGFILE.GetString("SESSION_DATABASE"),
//			TableName:       CONFIGFILE.GetString("SESSION_TABLENAME"),
//		})
//		if err != nil {
//			fmt.Println("failed to initialized mysql session provider:", err.Error())
//			break
//		}
//		provider = sessionProvider
//		break
//	case "postgresql", "postgres":
//		sessionProvider, err := postgres.New(postgres.Config{
//			Host:            CONFIGFILE.GetString("SESSION_HOST"),
//			Port:            CONFIGFILE.GetInt64("SESSION_PORT"),
//			Username:        CONFIGFILE.GetString("SESSION_USERNAME"),
//			Password:        CONFIGFILE.GetString("SESSION_PASSWORD"),
//			Database:        CONFIGFILE.GetString("SESSION_DATABASE"),
//			TableName:       CONFIGFILE.GetString("SESSION_TABLENAME"),
//		})
//		if err != nil {
//			fmt.Println("failed to initialized postgresql session provider:", err.Error())
//			break
//		}
//		provider = sessionProvider
//		break
//	case "redis":
//		sessionProvider, err := redis.New(redis.Config{
//			KeyPrefix:          CONFIGFILE.GetString("SESSION_KEYPREFIX"),
//			Addr:               CONFIGFILE.GetString("SESSION_HOST") + ":" + CONFIGFILE.GetString("SESSION_PORT"),
//			Password:           CONFIGFILE.GetString("SESSION_PASSWORD"),
//			DB:                 CONFIGFILE.GetInt("SESSION_DATABASE"),
//		})
//		if err != nil {
//			fmt.Println("failed to initialized redis session provider:", err.Error())
//			break
//		}
//		provider = sessionProvider
//		break
//	case "sqlite3":
//		sessionProvider, err := sqlite3.New(sqlite3.Config{
//			DBPath:          CONFIGFILE.GetString("SESSION_DATABASE"),
//			TableName:       CONFIGFILE.GetString("SESSION_TABLENAME"),
//		})
//		if err != nil {
//			fmt.Println("failed to initialized sqlite3 session provider:", err.Error())
//			break
//		}
//		provider = sessionProvider
//		break
//	}
//
//	return session.Config{
//		Lookup:     CONFIGFILE.GetString("SESSION_LOOKUP"),
//		Secure:     CONFIGFILE.GetBool("SESSION_SECURE"),
//		Domain:     CONFIGFILE.GetString("SESSION_DOMAIN"),
//		SameSite:   CONFIGFILE.GetString("SESSION_SAMESITE"),
//		Expiration: CONFIGFILE.GetDuration("SESSION_EXPIRATION"),
//		Provider:   provider,
//		GCInterval: CONFIGFILE.GetDuration("SESSION_GCINTERVAL"),
//	}
//}

func (config *Config) GetPgConfig() *pgxpool.Config {
	connString := config.GetPgConnString()
	pgconfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		log.Fatal().Msgf("postgresql config is fail %s", err)
	}
	dbPoolMaxLifeTime, errt := time.ParseDuration(config.GetString("database.postgresql.poolMaxConnLifeTime"))
	if errt != nil {
		log.Fatal().Msgf("postgresql poolMaxConnLifeTime  error %s", errt)
	}

	pgconfig.MaxConnLifetime = dbPoolMaxLifeTime
	dbPoolMaxConns := config.GetInt32("database.postgresql.poolMaxConns")
	if dbPoolMaxConns > 0 {
		pgconfig.MaxConns = dbPoolMaxConns
	}

	// Read current log level and use one less-fine level
	// below that
	//pgconfig.ConnConfig.Logger = zerolog.New()
	//levelString, _ := (log.GetLevel() - 1).MarshalText()
	//pgxLevel, _ := pgx.LogLevelFromString(string(levelString))
	//pgconfig.ConnConfig.LogLevel = pgxLevel
	return pgconfig
}

// postgres://jack:secret@pg.example.com:5432/mydb?sslmode=verify-ca&pool_max_conns=10
func (config *Config) GetPgConnString() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		config.GetString("database.postgresql.user"),
		config.GetString("database.postgresql.password"),
		config.GetString("database.postgresql.host"),
		config.GetInt32("database.postgresql.port"),
		config.GetString("database.postgresql.database"),
		config.GetString("database.postgresql.sslmode"))
}
