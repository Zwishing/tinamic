package routers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"time"
	"tinamic/app/controllers"
)

func RegisterAPI(api fiber.Router) {
	registerMvtServices(api)
	registerFeatureServices(api)
	registerRasterServices(api)
	registerData(api)
	registerUsers(api)
}

func registerMvtServices(api fiber.Router) {
	layer := api.Group("/mvt-services")

	//layer.Get("/get_table_layers", controllers.GetTableLayers)
	start:=time.Now()
	layer.Get("/:uuid/:z/:x/:y.pbf",controllers.GetTableLayerTile)
	fmt.Println(start.Sub(time.Now()))

}

func registerFeatureServices(api fiber.Router) {
	layer := api.Group("/feature-services")

	//layer.Get("/get_table_layers", controllers.GetTableLayers)
	layer.Get("/:uuid/:z/:x/:y.pbf",controllers.GetTableLayerTile)

}

func registerRasterServices(api fiber.Router) {
	layer := api.Group("/raster-services")

	//layer.Get("/get_table_layers", controllers.GetTableLayers)
	layer.Get("/:uuid/:z/:x/:y.pbf",controllers.GetTableLayerTile)

}

func registerData(api fiber.Router) {
	data := api.Group("/data")

	data.Post("/upload",controllers.Upload)
	data.Post("/publish",controllers.Publish)
	data.Get("/get_spatial_data", controllers.QuerySpatialData)

}

func registerUsers(api fiber.Router,) {
	users := api.Group("/users")

	users.Post("/register",controllers.Register)
	users.Post("/login",controllers.Login)
}