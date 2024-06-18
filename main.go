package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/rs/zerolog/log"
	"tinamic/conf"
	"tinamic/initialize"
	"tinamic/router"
)

func main() {
	app := initialize.InitApp()
	app.Use(cors.New())
	//router.SwaggerRoute(app.App)
	app.Get("/v1", func(ctx *fiber.Ctx) error {
		ctx.JSON("ok")
		return nil
	})
	api := app.Group("/v1")
	router.RegisterAPI(api)
	//设置端口监听
	err := app.Listen(fmt.Sprintf(":%d", conf.GetConfigInstance().GetInt("server.port")))
	if err != nil {
		log.Fatal().Msgf("%s", err)
	}
}
