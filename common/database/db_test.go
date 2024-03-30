package database

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"testing"
	"time"
	"tinamic/app/models"
)

var DbConnection = "postgresql://postgres:admin@localhost/postgres"

func TestDbConnect(t *testing.T) {
	db, _ := pgxpool.Connect(context.Background(), DbConnection)
	fmt.Println(db)
}

func TestInsert(t *testing.T) {
	uid, _ :=uuid.NewV4()
	data:=models.SpatialData{
		Uid:       uid,
		Name:      "shp",
		IsPublish: false,
		FileType:  "shapefile",
		Size:      1000,
		FilePath:  "",
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	}
	_, err := Insert("info","spatial_data",data)
	if err != nil {
		return 
	}
}
