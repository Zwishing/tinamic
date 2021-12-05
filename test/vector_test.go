package test

import (
	"fmt"
	"testing"
	"tinamic/app/queries"
)

func TestQueryGeometryType(t *testing.T) {
	gt:=queries.QueryGeometryType(db,"geom")
	fmt.Println(gt)
}

func TestQueryBounds(t *testing.T) {
	gt:=queries.QueryBounds(db,"geom")
	fmt.Println(gt)
}