package po

import "time"

type Task struct {
	ID           int
	TableId      string //上传表id
	TableVersion int    //上传表版本
	ProdId       string //产品id
	Status       int    //任务状态 0-失败 1-成功
	Stage        int    //任务阶段 0-上传中 1-导入中 2-完成
	tokenInfoId  int    //远程服务器token
	Created      *time.Time
	Updated      *time.Time
}
