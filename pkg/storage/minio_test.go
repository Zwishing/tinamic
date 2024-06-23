package storage

import (
	"testing"
)

func TestStorage_GetFiles(t *testing.T) {
	s := GetMinioInstance()
	s.GetFileByFolder("raster", "")

}
