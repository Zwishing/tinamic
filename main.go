package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/rs/zerolog/log"
	"tinamic/config"
	"tinamic/initialize"
	"tinamic/router"
)

func main() {
	app := initialize.InitApp()
	app.Use(cors.New())
	router.SwaggerRoute(app.App)
	api := app.Group("/api/v1")
	router.RegisterAPI(api)
	//设置端口监听
	err := app.Listen(fmt.Sprintf(":%d", config.Conf.GetInt("server.port")))
	if err != nil {
		log.Fatal().Msgf("%s", err)
	}
}
