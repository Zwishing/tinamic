package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	hashing "github.com/thomasvvugt/fiber-hashing"
	"log"
	configuration "tinamic/config"
	"tinamic/database"
	"tinamic/routers"
)

func main() {

	app:=InitApp()
	app.Use(cors.New())
	db,err :=database.DbConnect(app.Config.GetPgConfig())
	if err!=nil {
		fmt.Println("failed to connect to database:", err.Error())
	}
	routers.SwaggerRoute(app.App)
	api := app.Group("/api/v1")
	routers.RegisterAPI(api,db)

	log.Fatal(app.Listen(":8080"))
}


type App struct {
	*fiber.App
	Hasher hashing.Driver
	//Session *session.Session
	Config *configuration.Config
}

func InitApp() *App{
	config:=configuration.New()
	app := App{
		App:     fiber.New(*config.GetFiberConfig()),
		Hasher:  hashing.New(config.GetHasherConfig()),
		//Session: session.New(config.GetSessionConfig()),
		Config: config,
	}
	return &app
}