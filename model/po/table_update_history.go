package po

import "time"

// TableUpdateHistory 表更新记录
type TableUpdateHistory struct {
	ID           int
	TableId      string
	TableVersion int
	Contents     string //相较于上一个版本的更新内容，json格式，[{"type": "v", }]
	Created      *time.Time
	Updated      *time.Time
}
