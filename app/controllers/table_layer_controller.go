package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"

	"tinamic/app/models"
	"tinamic/app/queries"
	"tinamic/common/geos"
	"tinamic/common/query"
	"tinamic/common/response"
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
func AddTableLayer(db *pgxpool.Pool)fiber.Handler{
	return func(ctx *fiber.Ctx) error {
		tableLayer:=models.NewTableLayer()
		{
			tableLayer.Name="river"
			tableLayer.Attr= map[string]models.TableProperty{}
			tableLayer.Srid=queries.QuerySrid(db,tableLayer.GeometryColumn)
			tableLayer.Bounds= queries.QueryBounds(db,tableLayer.GeometryColumn)
			tableLayer.GeometryType=queries.QueryGeometryType(db,tableLayer.GeometryColumn)
			tableLayer.Center=queries.QueryCenter(db,tableLayer.GeometryColumn)
			tableLayer.CreatedAt=time.Now()
			tableLayer.UpdatedAt=time.Now()
		}

		tag, err := queries.InsertTableLayer(db,*tableLayer)
		if err != nil {
			return response.Fail(ctx,"",err.Error())
		}
		return response.Success(ctx,"", string(tag))
	}
}

func GetAllTableLayers(db *pgxpool.Pool) fiber.Handler{
	return func(ctx *fiber.Ctx) error{
		tableLayers, err := queries.QueryTableLayers(db)
		if err != nil {
			err := response.Fail(ctx, fiber.Map{}, "layers were not found")
			if err != nil {
				return err
			}
		}
		// Return status 200 OK.
		err = response.Success(ctx,tableLayers,"layers query success")
		if err != nil {
			return err
		}
		return nil
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