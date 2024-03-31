package test

import (
	"fmt"
	"testing"
	"tinamic/model"
)

func TestQueryGeometryType(t *testing.T) {
	gt := model.QueryGeometryType(db, "geom")
	fmt.Println(gt)
}

func TestQueryBounds(t *testing.T) {
	gt := model.QueryBounds(db, "geom")
	fmt.Println(gt)
}
