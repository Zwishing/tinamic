/*
 * @Author: zhang
 * @Date: 2021-11-28 10:17:19
 * @LastEditTime: 2021-11-28 13:37:18
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \tinamic\app\queries\layerinfo_query.go
 */
package queries

import (
	"context"
	"errors"
	"regexp"
	"strconv"
	"strings"
	. "tinamic/app/models"

	"github.com/jackc/pgx/v4/pgxpool"
)

// QueryLayerInfo 查询所有的图层
func QueryLayerInfo(db *pgxpool.Pool) ([]LayerInfo, error) {

	var layerInfo []LayerInfo
	// Send query to database.
	rows, err := db.Query(context.Background(), sqlLayerInfo)
	if err != nil {
		// Return empty object and error.
		return layerInfo, err
	}

	for rows.Next() {

		var info LayerInfo

		err := rows.Scan(&info.UID, &info.Schema, &info.Name, &info.Attr, &info.LayerType)
		if err != nil {
			return nil, err
		}
		layerInfo = append(layerInfo, info)
	}

	// Return query result.
	return layerInfo, nil
}

//QueryVersion 获取postgis版本
func QueryVersion(db *pgxpool.Pool) (map[string]string, int, error) {

	row := db.QueryRow(context.Background(), sqlPostGISVersion)
	var verStr string
	err := row.Scan(&verStr)
	if err != nil {
		return nil, 0, err
	}
	// Parse full version string
	//   POSTGIS="3.0.0 r17983" [EXTENSION] PGSQL="110" GEOS="3.8.0-CAPI-1.11.0 "
	//   PROJ="6.2.0" LIBXML="2.9.4" LIBJSON="0.13" LIBPROTOBUF="1.3.2" WAGYU="0.4.3 (Internal)"
	re := regexp.MustCompile(`([A-Z]+)="(.+?)"`)
	vers := make(map[string]string)
	for _, mtch := range re.FindAllStringSubmatch(verStr, -1) {
		vers[mtch[1]] = mtch[2]
	}

	pgisVer, ok := vers["POSTGIS"]
	if !ok {
		return nil, 0, errors.New("POSTGIS key missing from postgis_full_version")
	}
	// Convert Postgis version string into a lexically (and/or numerically) sortable form
	// "3.1.1 r17983" => "3001001"
	pgisMajMinPat := strings.Split(strings.Split(pgisVer, " ")[0], ".")
	pgisMaj, _ := strconv.Atoi(pgisMajMinPat[0])
	pgisMin, _ := strconv.Atoi(pgisMajMinPat[1])
	pgisPat, _ := strconv.Atoi(pgisMajMinPat[2])
	pgisNum := 1000000*pgisMaj + 1000*pgisMin + pgisPat
	vers["POSTGISFULL"] = strconv.Itoa(pgisNum)
	//globalVersions = vers
	globalPostGISVersion := pgisNum

	return vers, globalPostGISVersion, nil
}
