package controllers

import (
	"fmt"
	"testing"
)

func TestRecordFile(t *testing.T) {

}

func TestIsSupport(t *testing.T) {
	fmt.Println(IsSupport("/a.json",".geojson","json"))
}