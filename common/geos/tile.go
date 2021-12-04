package geos

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"math"
)

// Tile represents a single tile in a tile
// pyramid, usually referenced in a URL path
// of the form "Zoom/X/Y.Ext"
type Tile struct {
	Zoom   uint8    `json:"zoom"`
	X      int32    `json:"x"`
	Y      int32    `json:"y"`
	Ext    string `json:"ext"`
	Bounds Bounds `json:"bounds"`
}

// MakeTile 通过路由参数构造Tile对象
func MakeTile(zoom uint8,x int32,y int32) (tile *Tile, e error) {
	tile = &Tile{Zoom: zoom, X: x, Y: y, Ext: "pbf"}
	// No tile numbers outside the tile grid implied
	// by the zoom?
	if _,err:=tile.IsValid();err!=nil {
		return tile, err
	}
	//获取 tile bounds
	e = tile.CalculateBounds()
	return tile, nil
}

func (tile *Tile) Width() float64 {
	return math.Abs(tile.Bounds.Xmax - tile.Bounds.Xmin)
}

// IsValid tests that the tile contains
// only tile addresses that fit within the
// zoom level, and that the zoom level is
// not crazy large
func (tile *Tile) IsValid() (bool,error) {
	var err error
	if tile.Zoom > 32{
		return false, err
	}
	worldTileSize := int32(1) << int32(tile.Zoom)
	if tile.X >= worldTileSize ||
		tile.Y >= worldTileSize {
		return false,err
	}
	return true,nil
}

// CalculateBounds calculates the cartesian bounds that
// correspond to this tile
func (tile *Tile) CalculateBounds() (e error) {
	serverBounds, e := GetServerBounds()
	if e != nil {
		return e
	}

	worldWidthInTiles := float64(int(1) << uint(tile.Zoom))
	tileWidth := math.Abs(serverBounds.Xmax-serverBounds.Xmin) / worldWidthInTiles

	// Calculate geographic bounds from tile coordinates
	// XYZ tile coordinates are in "image space" so origin is
	// top-left, not bottom right
	xmin := serverBounds.Xmin + (tileWidth * float64(tile.X))
	xmax := serverBounds.Xmin + (tileWidth * float64(tile.X+1))
	ymin := serverBounds.Ymax - (tileWidth * float64(tile.Y+1))
	ymax := serverBounds.Ymax - (tileWidth * float64(tile.Y))
	tile.Bounds = Bounds{serverBounds.SRID, xmin, ymin, xmax, ymax}

	return nil
}

// String returns a path-like representation of the Tile
func (tile *Tile) String() string {
	return fmt.Sprintf("%d/%d/%d.%s", tile.Zoom, tile.X, tile.Y, tile.Ext)
}

type TileRequest struct {
	LayerID string
	Tile    Tile
	SQL     string
	Args    []interface{}
}


// GetServerBounds 获取全球服务边界
func GetServerBounds() (b *Bounds, e error) {

	srid := viper.GetInt("CoordinateSystem.SRID")
	xmin := viper.GetFloat64("CoordinateSystem.Xmin")
	ymin := viper.GetFloat64("CoordinateSystem.Ymin")
	xmax := viper.GetFloat64("CoordinateSystem.Xmax")
	ymax := viper.GetFloat64("CoordinateSystem.Ymax")

	log.Infof("Using CoordinateSystem.SRID %d with bounds [%g, %g, %g, %g]",
		srid, xmin, ymin, xmax, ymax)

	width := xmax - xmin
	height := ymax - ymin
	size := math.Min(width, height)

	/* Not square enough to just adjust */
	if math.Abs(width-height) > 0.01*size {
		return nil, errors.New("CoordinateSystem bounds must be square")
	}

	cx := xmin + width/2
	cy := ymin + height/2

	/* Perfectly square bounds please */
	xmin = cx - size/2
	ymin = cy - size/2
	xmax = cx + size/2
	ymax = cy + size/2

	globalServerBounds := &Bounds{SRID: srid, Xmin: xmin, Ymin: ymin, Xmax: xmax, Ymax: ymax}
	return globalServerBounds, nil
}