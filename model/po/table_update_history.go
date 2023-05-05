package po

import "time"

// TableUpdateHistory 表更新记录
type TableUpdateHistory struct {
	ID           int
	TableId      string
	TableVersion int
	Contents     string //相较于上一个版本的更新内容，json格式，[{"type": "v", "action": "add", ""}]
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

/*
[
{"type": "v", "action": "add", ""}


]
*/

/*
tables:
  tableId1:
    metadata: 保存创建语句
    data:
      vertex1: 点1
        metadata: 保存创建语句
        index: ["原索引列", "新索引列"]
        fulltext_index: ["索引列1", "索引列2", "索引列3"]
		data:
          vertex1_0
          vertex1_1
          ...
          vertex1_n
      vertex2: 点2
      ...
      edge1: 边1
      edge2: 边2
      ...
*/
