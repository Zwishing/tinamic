package wire

import (
	"github.com/gofiber/fiber/v2"
	hashing "github.com/thomasvvugt/fiber-hashing"
)

type App struct {
	*fiber.App
	Hasher hashing.Driver
	//Session *session.Session
}

func InitApp() *App {
	app := &App{
		App: fiber.New(),
		//Hasher: hashing.New(config.Conf.GetHasherConfig()),
		//Session: session.New(CONFIGFILE.GetSessionConfig()),
	}
	return app
}

func InitServices() (*UserComponents, *DataSourceComponents) {
	userHandler := InitializeUserService()
	dataSourceHandler, _, _ := InitializeDataSourceService()
	return userHandler, dataSourceHandler
}
