/*
 * @Author: your name
 * @Date: 2021-11-28 10:17:19
 * @LastEditTime: 2021-11-28 13:51:49
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \tinamic\app\controllers\baseinfo_controller.go
 */
package controllers

import (
	"github.com/gofiber/fiber/v2"

	"tinamic/app/queries"
	"tinamic/database"
)

func GetLyrBaseInfo(c *fiber.Ctx) error {
	db, err := database.DbConnect()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	lyrBaseInfo, err := queries.QueryLyrBaseInfo(db)
	if err != nil {
		// Return, if books not found.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "books were not found",
		})
	}
	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error":       false,
		"msg":         nil,
		"lyrBaseInfo": lyrBaseInfo,
	})
}
