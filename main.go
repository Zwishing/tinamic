package main

import (
	"github.com/gofiber/fiber/v2/middleware/cors"
	"tinamic/initialize"
	"tinamic/router"
)

func main() {
	app := initialize.InitApp()
	app.Use(cors.New())
	router.SwaggerRoute(app.App)
	api := app.Group("/api/v1")
	router.RegisterAPI(api)
}
