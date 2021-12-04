package queries

import (
	"context"
	"fmt"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
	//"github.com/spf13/viper"
	//"sort"
	"tinamic/app/models"
	"tinamic/common/geos"

)

// QueryTableLayers 查询所有的图层
func QueryTableLayers(db *pgxpool.Pool) ([]models.TableLayer, error) {

	var tableLayers []models.TableLayer
	// Send query to database.
	rows, err := db.Query(context.Background(), sqlLayerInfo)
	if err != nil {
		// Return empty object and error.
		return tableLayers, err
	}

	for rows.Next() {
		var tableLayer models.TableLayer
		err := rows.Scan(&tableLayer.UID, &tableLayer.Schema,
			             &tableLayer.Name, &tableLayer.Attr)
		if err != nil {
			return nil, err
		}
		tableLayers = append(tableLayers, tableLayer)
	}
	// Return query result.
	return tableLayers, nil
}

func QueryTableLayerByName(db *pgxpool.Pool,name string) models.TableLayer{
	var tableLayer models.TableLayer
	err := db.QueryRow(context.Background(), sqlTableLayerByName, name).Scan(tableLayer)
	if err != nil {
		panic("")
	}
	return tableLayer
}

func QueryLayerTile(db *pgxpool.Pool,tr *geos.TileRequest) ([]byte, error) {

	row := db.QueryRow(context.Background(), tr.SQL, tr.Args...)
	var mvtTile []byte
	err := row.Scan(&mvtTile)
	if err != nil {
		log.Warn(err)
	}
	return mvtTile,nil
}

func getBoundsExact(lyr *models.TableLayer,db *pgxpool.Pool) (geos.Bounds, error) {
	bounds := geos.Bounds{}
	extentSQL := fmt.Sprintf(sqlExtent, lyr.GeometryColumn, lyr.Srid, lyr.Schema, lyr.Name)

	var (
		xmin pgtype.Float8
		xmax pgtype.Float8
		ymin pgtype.Float8
		ymax pgtype.Float8
	)
	err := db.QueryRow(context.Background(), extentSQL).Scan(&xmin, &ymin, &xmax, &ymax)
	if err != nil {
		return bounds, err
	}

	bounds.SRID = 4326
	bounds.Xmin = xmin.Float
	bounds.Ymin = ymin.Float
	bounds.Xmax = xmax.Float
	bounds.Ymax = ymax.Float
	bounds.Sanitize()
	return bounds, nil
}

func getBounds(lyr *models.TableLayer,db *pgxpool.Pool) (geos.Bounds, error) {
	bounds := geos.Bounds{}
	extentSQL := fmt.Sprintf(`
		WITH ext AS (
			SELECT ST_Transform(ST_SetSRID(ST_EstimatedExtent('%s', '%s', '%s'), %d), 4326) AS geom
		)
		SELECT
			ST_XMin(ext.geom) AS xmin,
			ST_YMin(ext.geom) AS ymin,
			ST_XMax(ext.geom) AS xmax,
			ST_YMax(ext.geom) AS ymax
		FROM ext
		`, lyr.Schema, lyr.Name, lyr.GeometryColumn, lyr.Srid)

	var (
		xmin pgtype.Float8
		xmax pgtype.Float8
		ymin pgtype.Float8
		ymax pgtype.Float8
	)
	err := db.QueryRow(context.Background(), extentSQL).Scan(&xmin, &ymin, &xmax, &ymax)
	if err != nil {
		return bounds,err
	}

	// Failed to get estimate? Get the exact bounds.
	if xmin.Status == pgtype.Null {
		warning := fmt.Sprintf("Estimated extent query failed, run 'ANALYZE %s.%s'", lyr.Schema, lyr.Name)
		log.WithFields(log.Fields{
			"event": "request",
			"topic": "detail",
			"key":   warning,
		}).Warn(warning)
		return getBoundsExact(lyr,db)
	}

	bounds.SRID = 4326
	bounds.Xmin = xmin.Float
	bounds.Ymin = ymin.Float
	bounds.Xmax = xmax.Float
	bounds.Ymax = ymax.Float
	bounds.Sanitize()
	return bounds, nil
}

func WriteLayerJSON() error {
	return nil
}


//func getTableDetailJSON(lyr *LayerTable,db *pgxpool.Pool) (TableDetailJSON, error) {
//	td := TableDetailJSON{
//		ID:           lyr.ID,
//		Schema:       lyr.Schema,
//		Name:         lyr.Table,
//		Description:  lyr.Description,
//		GeometryType: lyr.GeometryType,
//		MinZoom:      viper.GetInt("DefaultMinZoom"),
//		MaxZoom:      viper.GetInt("DefaultMaxZoom"),
//	}
//	// TileURL is relative to server base
//	td.TileURL = fmt.Sprintf("%s/%s/{z}/{x}/{y}.pbf", serverURLBase(req), lyr.ID)
//
//	// Want to add the properties to the Json representation
//	// in table order, which is fiddly
//	tmpMap := make(map[int]TableProperty)
//	tmpKeys := make([]int, 0, len(lyr.Properties))
//	for _, v := range lyr.Properties {
//		tmpMap[v.Order] = v
//		tmpKeys = append(tmpKeys, v.Order)
//	}
//	sort.Ints(tmpKeys)
//	for _, v := range tmpKeys {
//		td.Properties = append(td.Properties, tmpMap[v])
//	}
//
//	// Read table bounds and convert to Json
//	// which prefers an array form
//	bnds, err := getBounds(lyr,db)
//	if err != nil {
//		return td, err
//	}
//	td.Bounds[0] = bnds.Xmin
//	td.Bounds[1] = bnds.Ymin
//	td.Bounds[2] = bnds.Xmax
//	td.Bounds[3] = bnds.Ymax
//	td.Center[0] = (bnds.Xmin + bnds.Xmax) / 2.0
//	td.Center[1] = (bnds.Ymin + bnds.Ymax) / 2.0
//	return td, nil
//}
