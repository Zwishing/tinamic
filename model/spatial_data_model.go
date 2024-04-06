package model

import (
	"github.com/gofrs/uuid"
	"os"
	"path/filepath"
	"strings"
	"time"
	"tinamic/util"
)

type SpatialData struct {
	Uid       uuid.UUID `json:"uid"`
	Name      string    `json:"name"`
	IsPublish bool      `json:"is_publish"`
	FileType  string    `json:"file_type"`
	Size      int64     `json:"size"`
	FilePath  string    `json:"file_path"`
	CreateAt  time.Time `json:"create_at"`
	UpdateAt  time.Time `json:"update_at"`
}

func NewSpatialData(fp string) *SpatialData {
	sd := new(SpatialData)
	uid, _ := uuid.NewV4()
	sd.Uid = uid
	sd.Name = strings.Split(util.GetFileName(fp), ".")[0]
	stat, err := os.Stat(fp)
	if err != nil {
		return nil
	}
	sd.Size = stat.Size()
	sd.FilePath, _ = filepath.Abs(fp)
	sd.CreateAt = stat.ModTime()
	sd.UpdateAt = time.Now()
	return sd
}
