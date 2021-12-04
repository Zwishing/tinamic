/*
 * @Author: zhang
 * @Date: 2021-11-28 10:17:19
 * @LastEditTime: 2021-11-28 12:50:54
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \tinamic\app\queries\sql.go
 */
package queries

//查询的sql语句
const (

	//
	sqlLayerInfo = `SELECT uid,schema,name,attr,layertype FROM layer.layerinfo`

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

	//
	sqlPostGISVersion = `SELECT postgis_full_version()`
)

// 创建表的sql语句
const (
	createLayerInfo = `CREATE TABLE layer.layerinfo(
		id  serial PRIMARY KEY,
		uid uuid NOT NULL UNIQUE,
		schema varchar(255) NOT NULL,
		name varchar(255) NOT NULL,
		attr json NOT NULL,
		layertype smallint NOT NULL default 1,
		description text, 
		createat timestamptz,
		updateat timestamptz 
	)`

	insertLayerInfoTest = ` INSERT INTO layer.layerinfo(
		uid, 
		schema,
		name,
		attr,
		layertype
	) VALUES(
		'123e4567-e89b-12d3-a456-426655440000',
		'layer',
		'city',
		'{
			"name":"string"
		}',
		1
	)
	`
)
