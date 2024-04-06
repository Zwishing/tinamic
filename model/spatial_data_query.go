package model

import (
	"context"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"tinamic/pkg/database"
)

func InsertSpatialData(data *SpatialData) (pgconn.CommandTag, error) {
	tag, err := database.Insert("info", "spatial_data", *data)
	if err != nil {
		return nil, err
	}
	return tag, nil
}

func QuerySpatialData() (spatialData []*SpatialData, err error) {
	err = pgxscan.Select(context.Background(), database.Db, &spatialData,
		`SELECT uid, name, is_publish, file_type,size,file_path,update_at FROM info.spatial_data order by update_at desc `)
	if err != nil {
		return nil, err
	}
	return spatialData, nil
}
