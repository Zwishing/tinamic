package router

import (
	"github.com/gofiber/fiber/v2"
	"tinamic/handler"
	"tinamic/middleware"
)

func registerUpload(api fiber.Router) {
	api.Get("/:dtype<range(0,1)>/post-upload", handler.CreatePostPresignedUrl)
	api.Get("/:dtype<range(0,1)>/put-upload", middleware.Protected(), handler.CreatePutPresignedUrl)
}
