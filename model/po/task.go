package po

import (
	"gorm.io/plugin/soft_delete"
	"time"
)

var TaskModel = &Task{}

type Task struct {
	ID          int
	KgId        string                `gorm:"type:char(64) not null;index:task_idx"` //上传图谱id
	ProdId      string                `gorm:"type:char(64) not null;index:task_idx"` //产品id
	Status      int                   `gorm:"type:tinyint(1) not null"`              //任务状态 1-上传中 2-失败 3-成功
	State       int                   `gorm:"type:tinyint(1) not null"`              //任务阶段 1-准备中 2-上传中 3-导入中 4-完成
	Token       string                `gorm:"type:char(64) not null;default ''"`     //远程服务器token
	StartIndex  int                   `gorm:"type:bigint(11) not null"`
	LastIndex   int                   `gorm:"type:bigint(11) not null;index:task_idx"` //上传的最大文件索引
	ErrorDetail string                `gorm:"varchar(500) not null;default ''"`
	Created     *time.Time            `json:"created" gorm:"autoCreateTime"`
	Updated     *time.Time            `json:"updated" gorm:"autoUpdateTime"`
	Deleted     soft_delete.DeletedAt `gorm:"type:tinyint(1);softDelete:flag,default:0"`
}
