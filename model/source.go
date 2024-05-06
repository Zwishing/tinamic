package model

import "github.com/google/uuid"

type VectorType = int

const (
	ShapeFile VectorType = iota + 1
	GeoJson
)

type VectorSource struct {
	Uuid     uuid.UUID  `json:"uuid,omitempty"`
	Name     string     `json:"name,omitempty"`
	DataType VectorType `json:"data_type,omitempty"`
	Size     int64      `json:"size,omitempty"`
	Layers   []string   `json:"layers,omitempty"`
	FilePath string     `json:"file_path,omitempty"`
}

func NewVectorSource(name string, dataType VectorType, size int64, layers []string, filePath string) *VectorSource {
	uid, err := uuid.NewV7()
	if err != nil {
		uid = uuid.New()
	}
	return &VectorSource{Uuid: uid, Name: name, DataType: dataType, Size: size, Layers: layers, FilePath: filePath}
}
