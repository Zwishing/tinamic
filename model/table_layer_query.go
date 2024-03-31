package model

//
//import (
//	"context"
//	"fmt"
//	"github.com/jackc/pgconn"
//	"github.com/jackc/pgtype"
//	log "github.com/sirupsen/logrus"
//	"os"
//	"runtime"
//	"time"
//	"tinamic/common/database"
//	"tinamic/common/query"
//	"tinamic/common/utils"
//
//	//"github.com/spf13/viper"
//	//"sort"
//	"tinamic/app/models"
//	"tinamic/common/geos"
//)
//
//// 添加一个tablelayer
//func InsertTableLayer(tableLayer models.TableLayer)(pgconn.CommandTag,error){
//	sql:=`INSERT INTO layers.table_layer(
//			uid,schema,name,attr,geometry_column,geometry_type,srid,bounds,center,create_at,update_at)
//         VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`
//	tag,err:= database.Db.Exec(context.Background(),sql,
//		tableLayer.UID,tableLayer.Schema,tableLayer.Name,tableLayer.Attr,tableLayer.GeometryColumn,
//		tableLayer.GeometryType,tableLayer.Srid,tableLayer.Bounds,tableLayer.Center,
//		tableLayer.CreatedAt, tableLayer.UpdatedAt)
//	if err!=nil{
//		return nil, err
//	}
//	return tag,nil
//}
//
//// QueryTableLayers 查询所有的图层
//func QueryTableLayers() ([]models.TableLayer, error) {
//
//	var tableLayers []models.TableLayer
//	// Send query to database.
//	rows, err := database.Db.Query(context.Background(), sqlLayerInfo)
//	if err != nil {
//		// Return empty object and error.
//		return tableLayers, err
//	}
//
//	for rows.Next() {
//		var tableLayer models.TableLayer
//		err := rows.Scan(&tableLayer.UID, &tableLayer.Schema,
//			             &tableLayer.Name, &tableLayer.Attr,&tableLayer.Srid)
//		if err != nil {
//			return nil, err
//		}
//		tableLayers = append(tableLayers, tableLayer)
//	}
//	// Return query result.
//	return tableLayers, nil
//}
//
//// QueryTableLayerByUid lyr chan *models.TableLayer
//func QueryTableLayerByUid(uuid string,) *models.TableLayer{
//	var tableLayer models.TableLayer
//	sql:=fmt.Sprintf(`SELECT schema,name,srid,geometry_column FROM layers.table_layer WHERE table_layer.uid='%s'`,uuid)
//	err := database.Db.QueryRow(context.Background(), sql).Scan(
//			&tableLayer.Schema,&tableLayer.Name,&tableLayer.Srid,&tableLayer.GeometryColumn)
//	if err != nil {
//		panic("")
//	}
//	return &tableLayer
//	//lyr<- &tableLayer
//}
//
//func QueryLayerTile(tr *geos.TileRequest) ([]byte, error) {
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//	row := database.Db.QueryRow(ctx, tr.SQL, tr.Args...)
//	var mvtTile []byte
//	err := row.Scan(&mvtTile)
//
//	if err != nil {
//		log.Warn(err)
//	}
//	return mvtTile,nil
//}
//
//func getBoundsExact(lyr *models.TableLayer) (geos.Bounds, error) {
//	bounds := geos.Bounds{}
//	extentSQL := fmt.Sprintf(sqlExtent, lyr.GeometryColumn, lyr.Srid, lyr.Schema, lyr.Name)
//
//	var (
//		xmin pgtype.Float8
//		xmax pgtype.Float8
//		ymin pgtype.Float8
//		ymax pgtype.Float8
//	)
//	err := database.Db.QueryRow(context.Background(), extentSQL).Scan(&xmin, &ymin, &xmax, &ymax)
//	if err != nil {
//		return bounds, err
//	}
//
//	bounds.SRID = 4326
//	bounds.Xmin = xmin.Float
//	bounds.Ymin = ymin.Float
//	bounds.Xmax = xmax.Float
//	bounds.Ymax = ymax.Float
//	bounds.Sanitize()
//	return bounds, nil
//}
//
//func getBounds(lyr *models.TableLayer) (geos.Bounds, error) {
//	bounds := geos.Bounds{}
//	extentSQL := fmt.Sprintf(`
//		WITH ext AS (
//			SELECT ST_Transform(ST_SetSRID(ST_EstimatedExtent('%s', '%s', '%s'), %d), 4326) AS geom
//		)
//		SELECT
//			ST_XMin(ext.geom) AS xmin,
//			ST_YMin(ext.geom) AS ymin,
//			ST_XMax(ext.geom) AS xmax,
//			ST_YMax(ext.geom) AS ymax
//		FROM ext
//		`, lyr.Schema, lyr.Name, lyr.GeometryColumn, lyr.Srid)
//
//	var (
//		xmin pgtype.Float8
//		xmax pgtype.Float8
//		ymin pgtype.Float8
//		ymax pgtype.Float8
//	)
//	err := database.Db.QueryRow(context.Background(), extentSQL).Scan(&xmin, &ymin, &xmax, &ymax)
//	if err != nil {
//		return bounds,err
//	}
//
//	// Failed to get estimate? Get the exact bounds.
//	if xmin.Status == pgtype.Null {
//		warning := fmt.Sprintf("Estimated extent query failed, run 'ANALYZE %s.%s'", lyr.Schema, lyr.Name)
//		log.WithFields(log.Fields{
//			"event": "request",
//			"topic": "detail",
//			"key":   warning,
//		}).Warn(warning)
//		return getBoundsExact(lyr)
//	}
//
//	bounds.SRID = 4326
//	bounds.Xmin = xmin.Float
//	bounds.Ymin = ymin.Float
//	bounds.Xmax = xmax.Float
//	bounds.Ymax = ymax.Float
//	bounds.Sanitize()
//	return bounds, nil
//}
//
//func WriteLayerJSON() error {
//	return nil
//}
//
//
//func Shp2PgSql(schema,fname, fp string,args ...string){
//	sp:=query.NewShpParameters(schema,fname,fp)
//	if runtime.GOOS == "windows" {
//		err := os.Setenv("PGPASSWORD", database.Db.Config().ConnConfig.Password)
//		if err != nil {
//			return
//		}
//		cmd:=fmt.Sprintf("shp2pgsql -s %s -W %s -I %s %s.%s | psql -h %s -d %s -p %d -U %s",
//											sp.Epsg,sp.ENCD,sp.ShpPath,sp.Schema,sp.LayerName,
//											database.Db.Config().ConnConfig.Host,
//											database.Db.Config().ConnConfig.Database,
//											database.Db.Config().ConnConfig.Port,
//											database.Db.Config().ConnConfig.User)
//
//		err = utils.ExeCommand("cmd","/C",cmd)
//		if err != nil {
//			return
//		}
//	}
//	//else {
//	//	dpBash:=fmt.Sprintf(`DATABASE=%s HOST=%s PASSWD=%s PORT=%d
//	//					`,dp.Database,dp.Host,dp.Password,dp.Port)
//	//	spBash:=fmt.Sprintf(`CRS_EPSG=%s ENCD=%s SHP=%s GEODATAFOLDER=%s SQLFILE=%s LAYERNAME=%s
//   //                 `,sp.Epsg,sp.ENCD,sp.Shp,sp.GeoDataFolder,sp.SqlFile,sp.LayerName)
//	//	err:=utils.ExeCommand("/bin/bash","-c",dpBash,spBash)
//	//	if err!=nil{
//	//		return
//	//	}
//	//}
//}
//
//
////func getTableDetailJSON(lyr *LayerTable,db *pgxpool.Pool) (TableDetailJSON, error) {
////	td := TableDetailJSON{
////		ID:           lyr.ID,
////		Schema:       lyr.Schema,
////		Name:         lyr.Table,
////		Description:  lyr.Description,
////		GeometryType: lyr.GeometryType,
////		MinZoom:      viper.GetInt("DefaultMinZoom"),
////		MaxZoom:      viper.GetInt("DefaultMaxZoom"),
////	}
////	// TileURL is relative to server base
////	td.TileURL = fmt.Sprintf("%s/%s/{z}/{x}/{y}.pbf", serverURLBase(req), lyr.ID)
////
////	// Want to add the properties to the Json representation
////	// in table order, which is fiddly
////	tmpMap := make(map[int]TableProperty)
////	tmpKeys := make([]int, 0, len(lyr.Properties))
////	for _, v := range lyr.Properties {
////		tmpMap[v.Order] = v
////		tmpKeys = append(tmpKeys, v.Order)
////	}
////	sort.Ints(tmpKeys)
////	for _, v := range tmpKeys {
////		td.Properties = append(td.Properties, tmpMap[v])
////	}
////
////	// Read table bounds and convert to Json
////	// which prefers an array form
////	bnds, err := getBounds(lyr,db)
////	if err != nil {
////		return td, err
////	}
////	td.Bounds[0] = bnds.Xmin
////	td.Bounds[1] = bnds.Ymin
////	td.Bounds[2] = bnds.Xmax
////	td.Bounds[3] = bnds.Ymax
////	td.Center[0] = (bnds.Xmin + bnds.Xmax) / 2.0
////	td.Center[1] = (bnds.Ymin + bnds.Ymax) / 2.0
////	return td, nil
////}
