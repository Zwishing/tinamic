/*
 * @Author: your name
 * @Date: 2021-11-28 10:17:19
 * @LastEditTime: 2021-11-28 12:54:09
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \tinamic\test\layerinfo_query_test.go
 */
package test

import (
	"fmt"
	"testing"
	"tinamic/app/queries"
	"tinamic/database"
)

func TestQueryVersion(t *testing.T) {
	var db, _ = database.DbConnect()
	ver, _, _ := queries.QueryVersion(db)
	fmt.Println(ver)
}

func TestQueryLyrBaseInfo(t *testing.T) {
	var db, _ = database.DbConnect()
	_, err := queries.QueryLayerInfo(db)
	if err != nil {
		return 
	}
}
