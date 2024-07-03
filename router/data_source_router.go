package router

import (
	"github.com/gofiber/fiber/v2"
	"tinamic/handler"
)

func registerDataSource(api fiber.Router, dataSourceHandler *handler.DataSourceHandler) {
	dataSource := api.Group("/data-source/:sourceType")
	dataSource.Get("/store-items", dataSourceHandler.GetStoreItems)
	dataSource.Get("/direct-upload", dataSourceHandler.GeneratePutPreSignedUrl)
	dataSource.Post("/upload", dataSourceHandler.Upload)
}
