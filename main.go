package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	hashing "github.com/thomasvvugt/fiber-hashing"
	"log"
	configuration "tinamic/config"
	"tinamic/database"
	. "tinamic/routers"
)

func main() {

	app:=NewApp()

	db,err :=database.DbConnect(app.Config.GetPgConfig())
	if err!=nil {
		fmt.Println("failed to connect to database:", err.Error())
	}

	api := app.Group("/api/v1")
	RegisterAPI(api,db)

	log.Fatal(app.Listen(":3001"))
}


type App struct {
	*fiber.App
	Hasher hashing.Driver
	//Session *session.Session
	Config *configuration.Config
}

func NewApp() *App{
	config:=configuration.New()
	app := App{
		App:     fiber.New(*config.GetFiberConfig()),
		Hasher:  hashing.New(config.GetHasherConfig()),
		//Session: session.New(config.GetSessionConfig()),
		Config: config,
	}
	return &app
}