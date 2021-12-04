/*
 * @Author: zhang
 * @Date: 2021-11-28 10:17:19
 * @LastEditTime: 2021-11-28 12:50:54
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \tinamic\app\queries\sql_template.go
 */
package queries

import (
	"fmt"
	"strings"

	"tinamic/app/models"
	"tinamic/common/geos"
	"tinamic/common/query"
)

//查询的sql语句
const (
	//
	sqlLayerInfo = `SELECT uid,schema,name,attr,layertype FROM layer.layerinfo;`
	sqlTableLayerByName = `SELECT * FROM layer.layerinfo WHERE name="%s";`

	sqlTableLayer = `
		SELECT ST_AsMVT(mvtgeom, {{ .MvtParams }}) FROM (
			SELECT ST_AsMVTGeom(
				ST_Transform(ST_Force2D(t."{{ .GeometryColumn }}"), {{ .TileSrid }}),
				bounds.geom_clip,
				{{ .Resolution }},
				{{ .Buffer }}
		  	) AS "{{ .GeometryColumn }}"
		  	{{ if .Properties }}
		  	, {{ .Properties }}
		  	{{ end }}
			FROM "{{ .Schema }}"."{{ .Table }}" t, (
				SELECT {{ .TileSQL }}  AS geom_clip,
						{{ .QuerySQL }} AS geom_query
				) bounds
			WHERE ST_Intersects(t."{{ .GeometryColumn }}",
								ST_Transform(bounds.geom_query, {{ .Srid }}))
			{{ .Limit }}
		) mvtgeom`
	sqlExtent =`
		WITH ext AS (
		SELECT
			coalesce(
				ST_Transform(ST_SetSRID(ST_Extent("%s"), %d), 4326),
				ST_MakeEnvelope(-180, -90, 180, 90, 4326)
			) AS geom
		FROM "%s"."%s"
	)
	SELECT
		ST_XMin(ext.geom) AS xmin,
		ST_YMin(ext.geom) AS ymin,
		ST_XMax(ext.geom) AS xmax,
		ST_YMax(ext.geom) AS ymax
	FROM ext
    `


)

func RequestSQL(lyr *models.TableLayer,tile *geos.Tile, qp *query.QueryParameters) (string, error) {

	type sqlParameters struct {
		TileSQL        string
		QuerySQL       string
		TileSrid       int
		Resolution     int
		Buffer         int
		Properties     string
		MvtParams      string
		Limit          string
		Schema         string
		Table          string
		GeometryColumn string
		Srid           int
	}

	// need both the exact tile boundary for clipping and an
	// expanded version for querying
	tileBounds := tile.Bounds
	queryBounds := tile.Bounds
	queryBounds.Expand(tile.Width() * float64(qp.Buffer) / float64(qp.Resolution))
	tileSQL := tileBounds.SQL()
	tileQuerySQL := queryBounds.SQL()

	// SRID of the tile we are going to generate, which might be different
	// from the layer SRID in the database
	tileSrid := tile.Bounds.SRID

	// preserve case and special characters in column names
	// of SQL query by double quoting names
	attrNames := make([]string, 0, len(qp.Properties))
	for _, a := range qp.Properties {
		attrNames = append(attrNames, fmt.Sprintf("\"%s\"", a))
	}

	// only specify MVT format parameters we have configured
	mvtParams := make([]string, 0)
	mvtParams = append(mvtParams, fmt.Sprintf("'%s', %d", lyr.UID, qp.Resolution))
	if lyr.GeometryColumn != "" {
		mvtParams = append(mvtParams, fmt.Sprintf("'%s'", lyr.GeometryColumn))
	}
	// The idColumn parameter is PostGIS3+ only
	if lyr.IDColumn != "" {
		mvtParams = append(mvtParams, fmt.Sprintf("'%s'", lyr.IDColumn))
	}

	sp := sqlParameters{
		TileSQL:        tileSQL,
		QuerySQL:       tileQuerySQL,
		TileSrid:       tileSrid,
		Resolution:     qp.Resolution,
		Buffer:         qp.Buffer,
		Properties:     strings.Join(attrNames, ", "),
		MvtParams:      strings.Join(mvtParams, ", "),
		Schema:         lyr.Schema,
		Table:          lyr.Name,
		GeometryColumn: lyr.GeometryColumn,
		Srid:           lyr.Srid,
	}

	if qp.Limit > 0 {
		sp.Limit = fmt.Sprintf("LIMIT %d", qp.Limit)
	}

	// TODO: Remove ST_Force2D when fixes to line clipping are common
	// in GEOS. See https://trac.osgeo.org/postgis/ticket/4690

	sql, err := query.RenderSQLTemplate("tabletilesql", sqlTableLayer, sp)
	if err != nil {
		return "", err
	}
	return sql, err
}