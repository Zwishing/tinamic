package router

import (
	"github.com/gofiber/fiber/v2"
	"tinamic/handler/user"
)

func registerUser(api fiber.Router) {
	//layer := api.Group("/user")
	api.Post("/login/account", user.Login)

}
