package storage

import (
	"context"
	"github.com/jackc/pgconn"
	"time"
	"tinamic/model"
	"tinamic/pkg/database"
)

func InsertVectorSource(vectorSource model.VectorSource) (pgconn.CommandTag, error) {
	sql := `INSERT INTO source_info.vectors(
			uuid,name,data_type,size,layers,file_path,created)
          VALUES($1,$2,$3,$4,$5)`
	tag, err := database.Db.Exec(context.Background(), sql, vectorSource.Uuid,
		vectorSource.Name, vectorSource.DataType, vectorSource.Size, vectorSource.Layers,
		vectorSource.FilePath, time.Now())
	if err != nil {
		return nil, err
	}
	return tag, nil
}
