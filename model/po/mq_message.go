package po

import (
	"gorm.io/plugin/soft_delete"
	"time"
)

var MqMessageModel = &MqMessage{}

type MqMessage struct {
	Id          int
	ServiceName string                `gorm:"type:varchar(50) not null;index:message_idx"`
	MessageId   string                `gorm:"type:varchar(100) not null;index:message_idx"`
	State       int                   `gorm:"type:tinyint(1) not null;default 1;index:message_idx"` // 1-未确认 2-已确认
	Retry       int                   `gorm:"type:int(11) not null;default 0"`
	Content     string                `gorm:"type:varchar(500) not null;default ''"`
	Created     *time.Time            `json:"created" gorm:"autoCreateTime"`
	Updated     *time.Time            `json:"updated" gorm:"autoUpdateTime"`
	Deleted     soft_delete.DeletedAt `gorm:"type:tinyint(1);softDelete:flag,default:0"`
}
