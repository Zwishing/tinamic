package query

import (
	"path/filepath"
)

type ShpParameters struct {
	Schema string
	LayerName string
	ShpPath string
	Epsg string
	ENCD string
}

func NewShpParameters(schema,fname,shpPath string) *ShpParameters {
	path,err:=filepath.Abs(shpPath)
	if err!=nil{
		panic(err)
	}
	return &ShpParameters{
		Schema:schema,
		LayerName: fname,
		Epsg: "4326",
		ShpPath: path,
		ENCD: "UTF-8",
	}
}
