package model

type ImageryFormat int8

const (
	COG ImageryFormat = iota + 1
	GTiff
)

func (i ImageryFormat) Ext() string {
	switch i {
	case COG:
		return "tif"
	case GTiff:
		return "tif"
	default:
		return ""
	}
}

type Imagery struct {
}

func (i *Imagery) GetType() DataSourceType {
	return ImageryType
}
