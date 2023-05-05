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

/**
更新行为：
  索引更新：
    全文索引：添加或删除索引列
    原生唯一索引：修改唯一索引列
  点或边更新
    添加点类
    添加边类
    删除点类
    删除边类
    修改点类或边类的显示名
添加数据：
  点数据
  边数据
*/
