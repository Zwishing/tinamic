package config

import (
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/logrusadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
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
	errorHandler fiber.ErrorHandler
	fiber *fiber.Config
}

var defaultErrorHandler = func (c *fiber.Ctx, err error) error {
	// Status code defaults to 500
	code := fiber.StatusInternalServerError

	// Set error message
	message := err.Error()

	// Check if it's a fiber.Error type
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	// TODO: Check return type for the client, JSON, HTML, YAML or any other (API vs web)

	// Return HTTP response
	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
	c.Status(code)

	// Render default error view
	err = c.Render("errors/" + strconv.Itoa(code), fiber.Map{"message": message})
	if err != nil {
		return c.SendString(message)
	}
	return err
}

func New() *Config {
	config := &Config{
		Viper: viper.New(),
	}

	// Set default configurations
	config.setDefaults()

	// Select the .env file
	config.AddConfigPath("./config")
	config.SetConfigName("tinamic")
	config.SetConfigType("toml")


	// Automatically refresh environment variables
	config.AutomaticEnv()

	// Read configuration
	if err := config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Println("failed to read configuration:", err.Error())
			os.Exit(1)
		}
	}

	config.SetErrorHandler(defaultErrorHandler)

	// TODO: Logger (Maybe a different zap object)

	// TODO: Add APP_KEY generation

	// TODO: Write changes to configuration file

	// Set Fiber configurations
	config.setFiberConfig()

	return config
}

func (config *Config) SetErrorHandler(errorHandler fiber.ErrorHandler) {
	config.errorHandler = errorHandler
}

func (config *Config) setDefaults()  {
	// Set default database connect config
	config.SetDefault("DbConnection", "sslmode=disable")
	// 1d, 1h, 1m, 1s, see https://golang.org/pkg/time/#ParseDuration
	config.SetDefault("DbPoolMaxConnLifeTime", "1h")
	config.SetDefault("DbPoolMaxConns", 4)
	config.SetDefault("DbTimeout", 10)
	config.SetDefault("CORSOrigins", []string{"*"})

	// Set default tile server config
	config.SetDefault("DefaultResolution", 4096)
	config.SetDefault("DefaultBuffer", 256)
	config.SetDefault("MaxFeaturesPerTile", 10000)
	config.SetDefault("DefaultMinZoom", 0)
	config.SetDefault("DefaultMaxZoom", 22)

	// Set default SRID
	config.SetDefault("CoordinateSystem.SRID", 3857)
	// XMin, YMin, XMax, YMax, must be square
	config.SetDefault("CoordinateSystem.Xmin", -20037508.3427892)
	config.SetDefault("CoordinateSystem.Ymin", -20037508.3427892)
	config.SetDefault("CoordinateSystem.Xmax", 20037508.3427892)
	config.SetDefault("CoordinateSystem.Ymax", 20037508.3427892)

	// Set default Fiber configuration
	config.SetDefault("FIBER_PREFORK", false)
	config.SetDefault("FIBER_SERVERHEADER", "")
	config.SetDefault("FIBER_STRICTROUTING", false)
	config.SetDefault("FIBER_CASESENSITIVE", false)
	config.SetDefault("FIBER_IMMUTABLE", false)
	config.SetDefault("FIBER_UNESCAPEPATH", false)
	config.SetDefault("FIBER_ETAG", false)
	config.SetDefault("FIBER_BODYLIMIT", 4194304)
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

	// Set default Fiber CORS middleware configuration
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
		Prefork:                   config.GetBool("FIBER_PREFORK"),
		ServerHeader:              config.GetString("FIBER_SERVERHEADER"),
		StrictRouting:             config.GetBool("FIBER_STRICTROUTING"),
		CaseSensitive:             config.GetBool("FIBER_CASESENSITIVE"),
		Immutable:                 config.GetBool("FIBER_IMMUTABLE"),
		UnescapePath:              config.GetBool("FIBER_UNESCAPEPATH"),
		ETag:                      config.GetBool("FIBER_ETAG"),
		BodyLimit:                 config.GetInt("FIBER_BODYLIMIT"),
		Concurrency:               config.GetInt("FIBER_CONCURRENCY"),
		Views:                     nil,
		ReadTimeout:               config.GetDuration("FIBER_READTIMEOUT"),
		WriteTimeout:              config.GetDuration("FIBER_WRITETIMEOUT"),
		IdleTimeout:               config.GetDuration("FIBER_IDLETIMEOUT"),
		ReadBufferSize:            config.GetInt("FIBER_READBUFFERSIZE"),
		WriteBufferSize:           config.GetInt("FIBER_WRITEBUFFERSIZE"),
		CompressedFileSuffix:      config.GetString("FIBER_COMPRESSEDFILESUFFIX"),
		ProxyHeader:               config.GetString("FIBER_PROXYHEADER"),
		GETOnly:                   config.GetBool("FIBER_GETONLY"),
		ErrorHandler:              config.errorHandler,
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

//func (config *Config) GetSessionConfig() session.Config {
//	var provider fsession.Provider
//	switch strings.ToLower(config.GetString("SESSION_PROVIDER")) {
//	case "memcache":
//		sessionProvider, err := memcache.New(memcache.Config{
//			KeyPrefix:    config.GetString("SESSION_KEYPREFIX"),
//			ServerList:   []string {
//				config.GetString("SESSION_HOST") + ":" + config.GetString("SESSION_PORT"),
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
//			Host:            config.GetString("SESSION_HOST"),
//			Port:            config.GetInt("SESSION_PORT"),
//			Username:        config.GetString("SESSION_USERNAME"),
//			Password:        config.GetString("SESSION_PASSWORD"),
//			Database:        config.GetString("SESSION_DATABASE"),
//			TableName:       config.GetString("SESSION_TABLENAME"),
//		})
//		if err != nil {
//			fmt.Println("failed to initialized mysql session provider:", err.Error())
//			break
//		}
//		provider = sessionProvider
//		break
//	case "postgresql", "postgres":
//		sessionProvider, err := postgres.New(postgres.Config{
//			Host:            config.GetString("SESSION_HOST"),
//			Port:            config.GetInt64("SESSION_PORT"),
//			Username:        config.GetString("SESSION_USERNAME"),
//			Password:        config.GetString("SESSION_PASSWORD"),
//			Database:        config.GetString("SESSION_DATABASE"),
//			TableName:       config.GetString("SESSION_TABLENAME"),
//		})
//		if err != nil {
//			fmt.Println("failed to initialized postgresql session provider:", err.Error())
//			break
//		}
//		provider = sessionProvider
//		break
//	case "redis":
//		sessionProvider, err := redis.New(redis.Config{
//			KeyPrefix:          config.GetString("SESSION_KEYPREFIX"),
//			Addr:               config.GetString("SESSION_HOST") + ":" + config.GetString("SESSION_PORT"),
//			Password:           config.GetString("SESSION_PASSWORD"),
//			DB:                 config.GetInt("SESSION_DATABASE"),
//		})
//		if err != nil {
//			fmt.Println("failed to initialized redis session provider:", err.Error())
//			break
//		}
//		provider = sessionProvider
//		break
//	case "sqlite3":
//		sessionProvider, err := sqlite3.New(sqlite3.Config{
//			DBPath:          config.GetString("SESSION_DATABASE"),
//			TableName:       config.GetString("SESSION_TABLENAME"),
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
//		Lookup:     config.GetString("SESSION_LOOKUP"),
//		Secure:     config.GetBool("SESSION_SECURE"),
//		Domain:     config.GetString("SESSION_DOMAIN"),
//		SameSite:   config.GetString("SESSION_SAMESITE"),
//		Expiration: config.GetDuration("SESSION_EXPIRATION"),
//		Provider:   provider,
//		GCInterval: config.GetDuration("SESSION_GCINTERVAL"),
//	}
//}

func (config *Config)GetPgConfig() *pgxpool.Config {
	dbConnection := config.GetString("DbConnection")
	pgconfig, err := pgxpool.ParseConfig(dbConnection)
	if err != nil {
		log.Fatal(err)
	}
	dbPoolMaxLifeTime, errt := time.ParseDuration(config.GetString("DbPoolMaxConnLifeTime"))
	if errt != nil {
		log.Fatal(errt)
	}

	pgconfig.MaxConnLifetime = dbPoolMaxLifeTime
	dbPoolMaxConns := config.GetInt32("DbPoolMaxConns")
	if dbPoolMaxConns > 0 {
		pgconfig.MaxConns = dbPoolMaxConns
	}

	// Read current log level and use one less-fine level
	// below that
	pgconfig.ConnConfig.Logger = logrusadapter.NewLogger(log.New())
	levelString, _ := (log.GetLevel() - 1).MarshalText()
	pgxLevel, _ := pgx.LogLevelFromString(string(levelString))
	pgconfig.ConnConfig.LogLevel = pgxLevel
	return pgconfig
}