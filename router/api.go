package router

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"time"
	"tinamic/handler"
	"tinamic/handler/vector"
)

func RegisterAPI(api fiber.Router, userhandler *handler.UserHandler, sourceHandler *handler.DataSourceHandler) {
	registerMvtService(api)
	registerFeatureService(api)
	registerRasterService(api)
	registerDataSource(api, sourceHandler)
	registerUser(api, userhandler)
	registerUpload(api)
}

func registerMvtService(api fiber.Router) {
	layer := api.Group("/mvt-services")

	//layer.Get("/get_table_layers", controllers.GetTableLayers)
	start := time.Now()
	layer.Get("/:uuid/:z/:x/:y.pbf", vector.GetTableLayerTile)
	fmt.Println(start.Sub(time.Now()))

}

func registerFeatureService(api fiber.Router) {
	layer := api.Group("/feature-services")

	//layer.Get("/get_table_layers", controllers.GetTableLayers)
	layer.Get("/:uuid/:z/:x/:y.pbf", vector.GetTableLayerTile)

}

func registerRasterService(api fiber.Router) {
	layer := api.Group("/raster-services")

	//layer.Get("/get_table_layers", controllers.GetTableLayers)
	layer.Get("/:uuid/:z/:x/:y.pbf", vector.GetTableLayerTile)

}
