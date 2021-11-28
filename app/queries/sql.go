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
	sqlLyrBaseInfo = `SELECT uid,schema,name,attr,layertype FROM layer.layerinfo`

	//
	sqlPostGISVersion = `SELECT postgis_full_version()`
)

// 创建表的sql语句
const (
	createLayerBaseInfo = `CREATE TABLE layer.layerinfo(
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

	insertLayerBaseInfoTest = ` INSERT INTO layer.layerinfo(
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
