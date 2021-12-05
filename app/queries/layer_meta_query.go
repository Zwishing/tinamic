package queries

import (
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/net/context"
	"strings"
	"tinamic/common/utils"
)

func QueryGeometryType(db *pgxpool.Pool,geometryColumn string) (geometryType string){
	sql:=fmt.Sprintf(`select ST_GeometryType("%s") from public.river`,geometryColumn)
	err := db.QueryRow(context.Background(), sql).Scan(&geometryType)
	if err != nil {
		return err.Error()
	}
	return strings.Split(geometryType,"_")[1]
}

func QuerySrid(db *pgxpool.Pool,geometryColumn string) (srid int){
	sql:=fmt.Sprintf(`select ST_SRID("%s") from public.river`,geometryColumn)
	err := db.QueryRow(context.Background(), sql).Scan(&srid)
	if err != nil {
		return 0
	}
	return srid
}

func QueryCenter(db *pgxpool.Pool,geometryColumn string) (center [2]float64) {
	bounds:=QueryBounds(db,geometryColumn)
	center[0]=(bounds[0]+bounds[2])/2
	center[1]=(bounds[1]+bounds[3])/2
	return center
}

func QueryBounds(db *pgxpool.Pool,geometryColumn string) [4]float64 {
	sql := fmt.Sprintf(`select ST_Extent("%s") from public.river`,geometryColumn)
	var bounds string
	err := db.QueryRow(context.Background(), sql).Scan(&bounds)
	if err != nil {
		return [4]float64{}
	}
	return utils.BoxStringToArray(bounds)
}


