package to

type KgInfoTo struct {
	State        int //图谱状态
	IsMetaUpdate int //元数据信息是否更新
	MaxIndex     int //当前redo log最大文件索引
}
