package queries

import (
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/net/context"
	"strings"
	"tinamic/common/database"
	"tinamic/common/utils"
)

func QueryGeometryColumn(db *pgxpool.Pool) {

}

func QueryGeometryType(geometryColumn string) (geometryType string) {
	sql := fmt.Sprintf(`select ST_GeometryType("%s") from public.river`, geometryColumn)
	err := database.Db.QueryRow(context.Background(), sql).Scan(&geometryType)
	if err != nil {
		return err.Error()
	}
	return strings.Split(geometryType, "_")[1]
}

func QuerySrid(geometryColumn string) (srid int) {
	sql := fmt.Sprintf(`select ST_SRID("%s") from public.river`, geometryColumn)
	err := database.Db.QueryRow(context.Background(), sql).Scan(&srid)
	if err != nil {
		return 0
	}
	return srid
}

func QueryCenter(geometryColumn string) (center [2]float64) {
	bounds := QueryBounds(geometryColumn)
	center[0] = (bounds[0] + bounds[2]) / 2
	center[1] = (bounds[1] + bounds[3]) / 2
	return center
}

func QueryBounds(geometryColumn string) [4]float64 {
	sql := fmt.Sprintf(`select ST_Extent("%s") from public.river`, geometryColumn)
	var bounds string
	err := database.Db.QueryRow(context.Background(), sql).Scan(&bounds)
	if err != nil {
		return [4]float64{}
	}
	return utils.BoxStringToArray(bounds)
}
