package initialize

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	hashing "github.com/thomasvvugt/fiber-hashing"
	"tinamic/config"
	"tinamic/pkg/database"
	"tinamic/pkg/storage"
)

type App struct {
	*fiber.App
	Hasher hashing.Driver
	//Session *session.Session
}

func InitApp() *App {
	app := &App{
		App:    fiber.New(),
		Hasher: hashing.New(config.Conf.GetHasherConfig()),
		//Session: session.New(CONFIGFILE.GetSessionConfig()),
	}
	// 初始化数据库
	initPg()
	// 初始化minio连接
	initMinio()

	return app
}

func initPg() {
	err := database.DbConnect(config.Conf.GetPgConfig())
	if err != nil {
		log.Fatal().Msgf("connect postgresql failed %s", err)
	}
}

func initMinio() {
	cfg := config.Conf.GetMinioConfig()
	minio, err := storage.New(cfg)
	if err != nil {
		log.Fatal().Msgf("connect minio failed %s", err)
	}
	storage.Minio = minio

}
