package routers

import "github.com/gofiber/fiber/v2"

func RegisterAPI(api fiber.Router) {
	registerRoles(api)
	registerUsers(api)
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