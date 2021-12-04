package controllers

import (
	"github.com/gofiber/fiber/v2"
	"tinamic/app/queries"
	"tinamic/database"
)

func GetTableLayer(ctx fiber.Ctx) error{
	db,err:= database.DbConnect()
	if err != nil {
		// Return status 500 and database connection error.
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	layerInfo, err := queries.QueryLayerInfo(db,)
	if err != nil {
		// Return, if books not found.
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "books were not found",
		})
	}

	// Return status 200 OK.
	return ctx.JSON(fiber.Map{
		"error":       false,
		"msg":         nil,
		"layerInfo": layerInfo,
	})

}
