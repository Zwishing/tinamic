/*
 * @Author: your name
 * @Date: 2021-11-28 10:17:19
 * @LastEditTime: 2021-11-28 13:51:49
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \tinamic\app\controllers\layerinfo_controller.go
 */
package controllers

import (
	"github.com/gofiber/fiber/v2"
	"tinamic/app/queries"
	"tinamic/database"
)

//func GetLayerInfo(db *pgxpool.Pool) fiber.Handler {
//	return func(ctx *fiber.Ctx) error {
//		layerInfo, err := queries.QueryLayerInfo(db)
//		if err!=nil{
//			return err
//		}
//		ctx.Status(200).JSON(fiber.Map{
//			"error":       false,
//			"msg":         nil,
//			"layerInfo": layerInfo,
//		})
//		return nil
//	}
//}

func GetLayerInfo(ctx *fiber.Ctx) error{
	db,err:=database.DbConnect()

	if err != nil {
		// Return status 500 and database connection error.
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	layerInfo, err := queries.QueryLayerInfo(db)
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
