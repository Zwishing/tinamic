package router

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"time"
	"tinamic/handler"
	"tinamic/handler/vector"
	"tinamic/pkg/middlewares"
)

func RegisterAPI(api fiber.Router) {
	registerMvtService(api)
	registerFeatureService(api)
	registerRasterService(api)
	registerData(api)
	registerUser(api)
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

func registerData(api fiber.Router) {
	//data := api.Group("/data")

	//data.Post("/upload", vector.Upload)
	//data.Post("/publish", vector.Publish)
	//data.Get("/get_spatial_data", vector.QuerySpatialData)

}

func registerUser(api fiber.Router) {
	user := api.Group("/user")

	user.Post("/register", handler.Register)
	user.Post("/login", handler.Login)
	user.Get("/profile", middlewares.Protected, handler.Profile)
}
