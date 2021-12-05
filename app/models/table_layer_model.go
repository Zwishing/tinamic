package models

import (
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
	"time"
)

type TableLayer struct {
	ID             int
	UID            uuid.UUID
	Schema         string
	Name           string
	Attr           map[string]TableProperty
	GeometryType   string
	IDColumn       string
	GeometryColumn string
	Srid           int
	Center         [2]float64
	Bounds         [4]float64
	MinZoom        int
	MaxZoom        int
	TileURL        string
	Description    string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type TableProperty struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Order       int
}

func NewTableLayer()*TableLayer{
	v4, err := uuid.NewV4()
	if err != nil {
		log.Fatal(err)
	}
	return &TableLayer{
		UID: v4,
		Schema:"layers",
		GeometryColumn:"geom",
		Srid:4326,
		MinZoom: 0,
		MaxZoom: 22,
	}
}

func (lyr *TableLayer) GetType() LayerType {
	return LayerTypeTable
}

func (lyr *TableLayer) GetUID() uuid.UUID {
	return lyr.UID
}

func (lyr *TableLayer) GetDescription() string {
	return lyr.Description
}

func (lyr *TableLayer) GetName() string {
	return lyr.Name
}

func (lyr *TableLayer) GetSchema() string {
	return lyr.Schema
}

//func (lyr *LayerTable) GetQueryPropertiesParameter(q url.Values) []string {
//	sAtts := make([]string, 0)
//	haveProperties := false
//
//	for k, v := range q {
//		if strings.EqualFold(k, "properties") {
//			sAtts = v
//			haveProperties = true
//			break
//		}
//	}
//
//	lyrAtts := (*lyr).Properties
//	queryAtts := make([]string, 0, len(lyrAtts))
//	haveIDColumn := false
//
//	if haveProperties {
//		aAtts := strings.Split(sAtts[0], ",")
//		for _, att := range aAtts {
//			decAtt, err := url.QueryUnescape(att)
//			if err == nil {
//				decAtt = strings.Trim(decAtt, " ")
//				att, ok := lyrAtts[decAtt]
//				if ok {
//					if att.Name == lyr.IDColumn {
//						haveIDColumn = true
//					}
//					queryAtts = append(queryAtts, att.Name)
//				}
//			}
//		}
//	}
//	// No request parameter or no matches, so we want to
//	// return all the properties in the table layer
//	if len(queryAtts) == 0 {
//		for _, v := range lyrAtts {
//			queryAtts = append(queryAtts, v.Name)
//		}
//	}
//	if (!haveIDColumn) && lyr.IDColumn != "" {
//		queryAtts = append(queryAtts, lyr.IDColumn)
//	}
//	return queryAtts
//}
//
//func (lyr *LayerTable) GetQueryParameters(ctx *fiber.Ctx) query.QueryParameters {
//	var rp query.QueryParameters
//	err := rp.GetQueryParameter(ctx)
//	if err != nil {
//		return query.QueryParameters{}
//	}
//	if rp.Limit < 0 {
//		rp.Limit = viper.GetInt("MaxFeaturesPerTile")
//	}
//	if rp.Resolution < 0 {
//		rp.Resolution = viper.GetInt("DefaultResolution")
//	}
//	if rp.Buffer < 0 {
//		rp.Buffer = viper.GetInt("DefaultBuffer")
//	}
//	return rp
//}