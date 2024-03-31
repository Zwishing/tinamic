package vector

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)
import "github.com/gofiber/fiber/v2/middleware/proxy"

// GetTableLayerTile
// import (
//
//	"github.com/gofiber/fiber/v2"
//	log "github.com/sirupsen/logrus"
//	"time"
//	"tinamic/app/models"
//	"tinamic/app/queries"
//	"tinamic/common/geos"
//	"tinamic/common/query"
//	"tinamic/common/response"
//
// )
//
// type TableLayerController struct {
// }
//
// //func GetTableLayer(ctx fiber.Ctx) error{
// //	db,err:= database.DbConnect()
// //	if err != nil {
// //		// Return status 500 and database connection error.
// //		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// //			"error": true,
// //			"msg":   err.Error(),
// //		})
// //	}
// //	layerInfo, err := queries.QueryLayerInfo(db,)
// //	if err != nil {
// //		// Return, if books not found.
// //		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
// //			"error": true,
// //			"msg":   "books were not found",
// //		})
// //	}
// //
// //	// Return status 200 OK.
// //	return ctx.JSON(fiber.Map{
// //		"error":       false,
// //		"msg":         nil,
// //		"layerInfo": layerInfo,
// //	})
// //
// //}
//
//	func addTableLayer(ctx *fiber.Ctx)error{
//			tableLayer:=models.NewTableLayer()
//			{
//				tableLayer.Name="river"
//				tableLayer.Attr= map[string]models.TableProperty{}
//				tableLayer.Srid=queries.QuerySrid(tableLayer.GeometryColumn)
//				tableLayer.Bounds= queries.QueryBounds(tableLayer.GeometryColumn)
//				tableLayer.GeometryType=queries.QueryGeometryType(tableLayer.GeometryColumn)
//				tableLayer.Center=queries.QueryCenter(tableLayer.GeometryColumn)
//				tableLayer.CreatedAt=time.Now()
//				tableLayer.UpdatedAt=time.Now()
//			}
//
//			tag, err := queries.InsertTableLayer(*tableLayer)
//			if err != nil {
//				return response.Fail(ctx,"",err.Error())
//			}
//			return response.Success(ctx,"", string(tag))
//	}
//
//	func GetTableLayers(ctx *fiber.Ctx) error{
//		tableLayers, err := queries.QueryTableLayers()
//		if err != nil {
//			err := response.Fail(ctx, fiber.Map{}, "layers were not found")
//			if err != nil {
//				return err
//			}
//		}
//		// Return status 200 OK.
//		err = response.Success(ctx,tableLayers,"layers query success")
//		if err != nil {
//			return err
//		}
//		return nil
//	}
//
//	func getTableLayerByName() models.TableLayer{
//		return models.TableLayer{}
//	}
//
// // 获取瓦片数据
func GetTableLayerTile(ctx *fiber.Ctx) error {

	z, _ := ctx.ParamsInt("z")
	x, _ := ctx.ParamsInt("x")
	y, _ := ctx.ParamsInt("y")
	//parms:=ctx.Params("params")
	//var mvt []byte
	////log.Info(tr)
	//go Mvt(mvt,x,y,z,"")
	ctx.Set(fiber.HeaderContentType, "application/vnd.mapbox-vector-tile")
	ctx.Set(fiber.HeaderAcceptEncoding, "gzip")
	//err := ctx.Send(mvt)
	//if err != nil {
	//	return err
	//}

	url := fmt.Sprintf("http://localhost:7800/dwh_gis.g_administrative_boundary/%d/%d/%d.pbf?)", z, x, y)
	if err := proxy.Do(ctx, url); err != nil {
		return err
	}
	// 从响应中删除 HeaderServer
	ctx.Response().Header.Del(fiber.HeaderServer)

	return nil

}

//
//func getTileRequest(lyr *models.TableLayer,tile geos.Tile, ctx *fiber.Ctx) geos.TileRequest {
//	qp:=query.NewQueryParameters()
//	_ = qp.GetQueryParameter(ctx)
//	sql, _ := queries.RequestSQL(lyr, &tile, qp)
//
//	tr := geos.TileRequest{
//		LayerID: lyr.Name,
//		Tile:    tile,
//		SQL:     sql,
//		Args:    nil,
//	}
//	return tr
//}
