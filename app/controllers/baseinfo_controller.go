package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"

	"tinamic/app/queries"
	"tinamic/database"
)

func GetLayerBaseInfo(c *fiber.Ctx) error {
	db,err:=database.DbConnect()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	layerinfo, err := queries.QueryLayerInfo(db)
	if err != nil {
		// Return, if books not found.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "books were not found",
			"count": 0,
			"books": nil,
		})
	}
	fmt.Println(layerinfo)
	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
	})
}