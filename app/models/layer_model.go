package models

import (
	"net/http"
)

// LayerType 用来定义是基本图层还是函数图层
type LayerType int

const (
	// LayerTypeTable is a table layer
	LayerTypeTable = 1
	// LayerTypeFunction is a function layer
	LayerTypeFunction = 2
)

func (lt LayerType) String() string {
	switch lt {
	case LayerTypeTable:
		return "table"
	case LayerTypeFunction:
		return "function"
	default:
		return "unknown"
	}
}

// Layer 接口用来实现是基本图层还是函数图层
type Layer interface {
	GetType() LayerType
	GetID() string
	GetDescription() string
	GetName() string
	GetSchema() string
	GetTileRequest(tile Tile, r *http.Request) TileRequest
	WriteLayerJSON(w http.ResponseWriter, req *http.Request) error
}



