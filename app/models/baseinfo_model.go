/*
 * @Author: your name
 * @Date: 2021-11-28 10:17:19
 * @LastEditTime: 2021-11-28 13:23:26
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \tinamic\app\models\baseinfo_model.go
 */
package models

import (
	"time"

	"github.com/gofrs/uuid"
)

// LayerBaseInfo 管理所有图层信息
type LayerBaseInfo struct {
	UID         uuid.UUID
	Schema      string
	Name        string
	Attr        map[string]string
	LayerType   int8
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
