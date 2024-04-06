package router

import (
	"github.com/gofiber/fiber/v2"
	"tinamic/handler/vector"
)

func registerUpload(api fiber.Router) {
	vec := api.Group("/vector")
	vec.Get("/upload", vector.UploadToMinio)
}
