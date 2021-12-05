package utils

import (
	"strconv"
	"strings"
)

// 'BOX(111.964430789849 30.9798174714992,116.727509619973 34.8957377544938)' to [4]float
func BoxStringToArray(box2d string)(bounds [4]float64){

	result:=strings.FieldsFunc(box2d, func(r rune) bool {
		if r=='B'||r=='O'||r=='X'||r=='('|| r==' '||r==','||r==')'{
			return true
		}else {
			return false
		}
	})
	bounds[0], _ =strconv.ParseFloat(result[0],64)
	bounds[1], _ =strconv.ParseFloat(result[1],64)
	bounds[2], _ =strconv.ParseFloat(result[2],64)
	bounds[3], _ =strconv.ParseFloat(result[3],64)
	return bounds
}
