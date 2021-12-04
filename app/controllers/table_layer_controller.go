package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"
	"tinamic/app/models"
	"tinamic/app/queries"
	"tinamic/common/geos"
	"tinamic/common/query"
)

//func GetTableLayer(ctx fiber.Ctx) error{
//	db,err:= database.DbConnect()
//	if err != nil {
//		// Return status 500 and database connection error.
//		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
//			"error": true,
//			"msg":   err.Error(),
//		})
//	}
//	layerInfo, err := queries.QueryLayerInfo(db,)
//	if err != nil {
//		// Return, if books not found.
//		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
//			"error": true,
//			"msg":   "books were not found",
//		})
//	}
//
//	// Return status 200 OK.
//	return ctx.JSON(fiber.Map{
//		"error":       false,
//		"msg":         nil,
//		"layerInfo": layerInfo,
//	})
//
//}

func GetAllTableLayers(db *pgxpool.Pool) fiber.Handler{
	return func(ctx *fiber.Ctx) error{
		tableLayers, err := queries.QueryTableLayers(db)
		if err != nil {
			err := ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": true,
				"msg":   "layers were not found",
			})
			if err != nil {
				return err
			}
		}

		// Return status 200 OK.
		err = ctx.JSON(fiber.Map{
			"error":     false,
			"msg":       nil,
			"Layers": tableLayers,
		})
		if err != nil {
			return err
		}
		return err
	}
}

func getTableLayerByName(db *pgxpool.Pool) models.TableLayer{
	return models.TableLayer{}
}

func GetTableLayerTile(db *pgxpool.Pool) fiber.Handler{
	return func(ctx *fiber.Ctx) error {

		z,_:=ctx.ParamsInt("z")
		x,_:=ctx.ParamsInt("x")
		y,_:=ctx.ParamsInt("y")

		lyr:=queries.QueryTableLayerByName(db,ctx.Params("name"))
		tile, _ :=geos.MakeTile(uint8(z),int32(x),int32(y))
		tr:=getTileRequest(&lyr,*tile,ctx)

		mvt, _ := queries.QueryLayerTile(db,&tr)
		err := ctx.Send(mvt)
		if err != nil {
			return err
		}
		return nil
	}
}

func getTileRequest(lyr *models.TableLayer,tile geos.Tile, ctx *fiber.Ctx) geos.TileRequest {
	var qp query.QueryParameters
	_ = qp.GetQueryParameter(ctx)
	sql, _ := queries.RequestSQL(lyr, &tile, &qp)

	tr := geos.TileRequest{
		LayerID: lyr.Name,
		Tile:    tile,
		SQL:     sql,
		Args:    nil,
	}
	return tr
}