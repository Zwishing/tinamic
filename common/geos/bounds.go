package geos

import (
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"math"
)

// Bounds represents a box in Web Mercator space
type Bounds struct {
	SRID int     `json:"srid"`
	Xmin float64 `json:"xmin"`
	Ymin float64 `json:"ymin"`
	Xmax float64 `json:"xmax"`
	Ymax float64 `json:"ymax"`
}

func (b *Bounds) String() string {
	return fmt.Sprintf("{Xmin:%g, Ymin:%g, Xmax:%g, Ymax:%g, SRID:%d}",
		b.Xmin, b.Ymin, b.Xmax, b.Ymax, b.SRID)
}

// SQL returns the SQL fragment to create this bounds in the database
func (b *Bounds) SQL() string {
	return fmt.Sprintf("ST_MakeEnvelope(%g, %g, %g, %g, %d)",
		b.Xmin, b.Ymin,
		b.Xmax, b.Ymax, b.SRID)
}

// Expand increases the size of this bounds in all directions, respecting
// the limits of the Web Mercator plane
func (b *Bounds) Expand(size float64) {
	serverBounds,err:= GetServerBounds()
	if err!=nil {
		panic("")
	}
	b.Xmin = math.Max(b.Xmin-size, serverBounds.Xmin)
	b.Ymin = math.Max(b.Ymin-size, serverBounds.Ymin)
	b.Xmax = math.Min(b.Xmax+size, serverBounds.Xmax)
	b.Ymax = math.Min(b.Ymax+size, serverBounds.Ymax)
	return
}

// func fromMercator(x float64, y float64) (lng float64, lat float64) {
// 	// worldMercWidth is the width of the Web Mercator plane
// 	worldMercWidth := 40075016.6855784
// 	mercSize := worldMercWidth / 2.0
// 	lng = x * 180.0 / mercSize
// 	lat = 180.0 / math.Pi * (2.0*math.Atan(math.Exp((y/mercSize)*math.Pi)) - math.Pi/2.0)
// 	return lng, lat
// }

func (b *Bounds) Sanitize() {
	if b.SRID == 4326 {
		if b.Ymin < -90 {
			b.Ymin = 90
		}
		if b.Ymax > 90 {
			b.Ymax = 90
		}
		if b.Xmin < -180 {
			b.Xmin = -180
		}
		if b.Xmax > 180 {
			b.Xmax = 180
		}
	}
	return
}

// Json returns the bounds in array for form consumption
// by Json formats that like it that way
// func (b *Bounds) Json() []float64 {
// 	s := make([]float64, 4)
// 	s[0], s[1] = fromMercator(b.Xmin, b.Ymin)
// 	s[2], s[3] = fromMercator(b.Xmax, b.Ymax)
// 	return s
// }

// Center returns the center of the bounds in array format
// for consumption by Json formats that like it that way
// func (b *Bounds) Center() []float64 {
// 	xc := (b.Xmin + b.Xmax) / 2.0
// 	yc := (b.Ymin + b.Ymax) / 2.0
// 	s := make([]float64, 2)
// 	s[0], s[1] = fromMercator(xc, yc)
// 	return s
// }

func getServerBounds() (b *Bounds, e error) {

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

	globalServerBounds := &Bounds{srid, xmin, ymin, xmax, ymax}
	return globalServerBounds, nil
}