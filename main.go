package main

import (
	"fmt"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/rs/zerolog/log"
	"tinamic/conf"
	"tinamic/router"
	"tinamic/wire"
)

// @title Tinamic服务API
// @version 1.0
// @description
// @termsOfService

// @contact.name API Support
// @contact.url
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /api/v1
func main() {
	app := wire.InitApp()
	app.Use(cors.New())

	// swagger 配置
	cfg := swagger.Config{
		BasePath: "/",
		FilePath: "./docs/swagger.json",
		Path:     "swagger",
		Title:    "Tinamic API Docs",
	}
	app.Use(swagger.New(cfg))

	api := app.Group("/v1")

	router.RegisterAPI(api)
	//设置端口监听
	err := app.Listen(fmt.Sprintf(":%d", conf.GetConfigInstance().GetInt("server.port")))
	if err != nil {
		log.Fatal().Msgf("%s", err)
	}
}
