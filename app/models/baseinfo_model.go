package models

import (
	"github.com/jackc/pgtype"
)

// LayerBaseInfo 管理所有图层信息
type LayerBaseInfo struct {
	ID pgtype.UUID
	Schema pgtype.Text
	Name pgtype.Text
	Attr map[pgtype.Text]pgtype.Text
	LayerType pgtype.Int2
	Description pgtype.Text
	CreatedAt    pgtype.Timestamptz
	UpdatedAt    pgtype.Timestamptz
}

