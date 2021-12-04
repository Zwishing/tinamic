package routers

import (
	"github.com/gofiber/fiber/v2"

	."tinamic/app/controllers"
)

func RegisterAPI(api fiber.Router) {
	registerLayer(api)
	//registerRoles(api)
	//registerUsers(api)
}

func registerLayer(api fiber.Router) {
	layer := api.Group("/layer")

	layer.Get("/layerinfo", GetLayerInfo)
	layer.Get("/tablelayer/:name/:z/:x/:y.pbf",GetLayerInfo)
	//layer.Post("/")
	//layer.Put("/:id")
	//layer.Delete("/:id")
}

func registerRoles(api fiber.Router) {
	roles := api.Group("/roles")

	roles.Get("/", )
	roles.Get("/:id")
	roles.Post("/")
	roles.Put("/:id")
	roles.Delete("/:id")
}

func registerUsers(api fiber.Router) {
	users := api.Group("/users")

	users.Get("/")
	users.Get("/:id")
	users.Post("/")
	users.Put("/:id")
	users.Delete("/:id")
}