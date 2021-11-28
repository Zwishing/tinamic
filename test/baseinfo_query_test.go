package test

import (
	"fmt"
	"testing"
	"tinamic/app/queries"
	"tinamic/database"
)

var db, _ = database.DbConnect()

func TestQueryVersion(t *testing.T) {
	ver,_,_:=queries.QueryVersion(db)
	fmt.Println(ver)
}