package po

type TaskProgress struct {
	ID             int
	TaskId         string
	UploadContents string //要上传的文件列表，json格式 [{""}, {""}]
	UploadIndex    string //上传位置
}
