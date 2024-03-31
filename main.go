package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	log "github.com/sirupsen/logrus"
	hashing "github.com/thomasvvugt/fiber-hashing"
	"tinamic/common/database"
	"tinamic/config"
	"tinamic/router"
)

func main() {

	app := InitApp()
	app.Use(cors.New())
	err := database.DbConnect(config.Conf.GetPgConfig())
	if err != nil {
		fmt.Println("failed to connect to database:", err.Error())
	}
	router.SwaggerRoute(app.App)
	api := app.Group("/api/v1")
	router.RegisterAPI(api)

	log.Fatal(app.Listen(":8083"))
}

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
	return &app
}
