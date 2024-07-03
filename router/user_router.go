package router

import (
	"github.com/gofiber/fiber/v2"
	"tinamic/handler"
	"tinamic/middleware"
)

func registerUser(api fiber.Router, handler *handler.UserHandler) {
	user := api.Group("/user")

	user.Post("/register", handler.Register)
	user.Post("/login", handler.Login)
	user.Get("/profile", middleware.Protected(), middleware.Authz().RoutePermission(), handler.Profile)
}
