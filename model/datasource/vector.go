package datasource

import "github.com/google/uuid"

type VectorFormat int8

const (
	Shapefile VectorFormat = iota + 1
	GeoJSON
)

func (vf VectorFormat) Ext() string {
	switch vf {
	case Shapefile:
		return "shp"
	case GeoJSON:
		return "geojson"
	default:
		return ""
	}
}

type VectorSource struct {
	uuid        uuid.UUID       `json:"uuid,omitempty"`
	name        string          `json:"name,omitempty"`
	format      VectorFormat    `json:"format,omitempty"`
	size        int64           `json:"size,omitempty"`
	compression CompressionType `json:"compression,omitempty"`
	filePath    FilePath        `json:"file_path"`
}

func NewVectorSource(name string, format VectorFormat, size int64, compression CompressionType, filePath FilePath) *VectorSource {
	return &VectorSource{
		uuid:        uuid.New(),
		name:        name,
		format:      format,
		size:        size,
		compression: compression,
		filePath:    filePath,
	}
}

func (v *VectorSource) GetUuid() string {
	return v.uuid.String()
}

func (v *VectorSource) GetType() DataSourceType {
	return VectorType
}

func (v *VectorSource) GetExt() string {
	return v.format.Ext()
}

func (v *VectorSource) GetSize() int64 {
	return v.size
}

func (v *VectorSource) GetName() string {
	return v.name
}

func (v *VectorSource) GetFilePath() string {
	return v.filePath.GetPath()
}

func (v *VectorSource) IsCompressed() bool {
	return !(v.compression == Uncompressed)
}

// GetCompressedExt TODO error handle
func (v *VectorSource) GetCompressedExt() string {
	if v.IsCompressed() {
		v.compression.Ext()
	}
	return ""
}
