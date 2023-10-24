package po

import (
	"gorm.io/plugin/soft_delete"
	"time"
)

var TokenInfoModel = &TokenInfo{}

type TokenInfo struct {
	ID      int
	Token   string                `gorm:"type:char(64) not null"`
	ProdId  string                `gorm:"type:char(64) not null"`
	Created *time.Time            `json:"created" gorm:"autoCreateTime"`
	Updated *time.Time            `json:"updated" gorm:"autoUpdateTime"`
	Deleted soft_delete.DeletedAt `gorm:"type:tinyint(1);softDelete:flag,default:0"`
}
