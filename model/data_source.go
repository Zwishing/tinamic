package model

type DataSourceType int8

const (
	VectorType DataSourceType = iota + 1
	ImageryType
	TerrainType
	PointCloudType
	PhotogrammetryType
	SurfaceModelType
	BIMType
)

func (dst DataSourceType) String() string {
	switch dst {
	case VectorType:
		return "vector"
	case ImageryType:
		return "imagery"
	case TerrainType:
		return "terrain"
	case PointCloudType:
		return "point cloud"
	case PhotogrammetryType:
		return "photogrammetry"
	case SurfaceModelType:
		return "surface model"
	case BIMType:
		return "building information model"
	default:
		return "unknown DataSourceType"
	}
}

type DataSource interface {
	GetUuid() string
	GetName() string
	GetType() DataSourceType
	GetExt() string
	GetSize() int64
	GetFilePath() string
}

type Extension interface {
	Ext() string
}

type CompressionType int8

const (
	Uncompressed CompressionType = iota
	Zip
)

func (ct CompressionType) String() string {
	switch ct {
	case Uncompressed:
		return "uncompressed"
	case Zip:
		return "zip"
	default:
		return "unknown CompressionType"
	}
}

func (ct CompressionType) Ext() string {
	switch ct {
	case Zip:
		return "zip"
	default:
		return ""
	}
}

type PathType int8

const (
	Nas PathType = iota
	Minio
)

func (pt PathType) String() string {
	switch pt {
	case Nas:
		return "NAS-local storage"
	case Minio:
		return "minio-object storage"
	default:
		return "unknown FilePathType"
	}
}

type FilePath struct {
	pathType PathType
	path     string
}

func (fp *FilePath) GetPath() string {
	return fp.path
}
