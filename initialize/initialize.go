package initialize

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	hashing "github.com/thomasvvugt/fiber-hashing"
	"tinamic/config"
	"tinamic/pkg/database"
)

type App struct {
	*fiber.App
	Hasher hashing.Driver
	//Session *session.Session
	Config *config.Config
}

func InitApp() *App {
	app := App{
		App:    fiber.New(*config.Conf.GetFiberConfig()),
		Hasher: hashing.New(config.Conf.GetHasherConfig()),
		//Session: session.New(CONFIGFILE.GetSessionConfig()),
		Config: config.Conf,
	}
	err := app.Listen(string(config.Conf.GetInt32("server.port")))
	if err != nil {
		log.Fatal().Msgf("port is already in use %s", err)
	}
	initPg()

	return &app
}

func initPg() {
	err := database.DbConnect(config.Conf.GetPgConfig())
	if err != nil {
		log.Fatal().Msgf("connect postgresql failed %s", err)
	}
}
