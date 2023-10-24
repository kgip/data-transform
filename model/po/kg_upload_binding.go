package po

import (
	"gorm.io/plugin/soft_delete"
	"time"
)

var KgUploadBindingModel = &KgUploadBinding{}

type KgUploadBinding struct {
	ID          int
	KgId        string                `gorm:"type:char(64) not null;index:binding_idx"` //上传图谱id
	TokenInfoId int                   `gorm:"type:int(11) not null;index:binding_idx"`
	MaxIndex    int                   `gorm:"type:bigint(11) not null;default 0;index:binding_idx"` //当前上传的最大文件索引，下一次上传从下一个文件开始
	Created     *time.Time            `json:"created" gorm:"autoCreateTime"`
	Updated     *time.Time            `json:"updated" gorm:"autoUpdateTime"`
	Deleted     soft_delete.DeletedAt `gorm:"type:tinyint(1);softDelete:flag,default:0"`
}
