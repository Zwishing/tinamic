package queries

//import (
//	"context"
//	"github.com/jackc/pgx/v4/pgxpool"
//	log "github.com/sirupsen/logrus"
//	."tinamic/app/models"
//)

//func QueryTableLayer(db *pgxpool.Pool,name string,z uint8,x int32,y int32) ([]LayerInfo, error) {
//
//	row := db.QueryRow(context.Background(), tr.SQL, tr.Args...)
//	var mvtTile []byte
//	err := row.Scan(&mvtTile)
//	if err != nil {
//		log.Warn(err)
//	}
//}

//func composeSql(name string,z uint8,x int32,y int32) (string,error) {
//	tile, err :=MakeTile(z,x,y)
//	if err!=nil{
//		return nil,err
//	}
//
//
//
//}