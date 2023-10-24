package http

type OssGateway interface {
	GetOsBucketInfos() ([]map[string]interface{}, error)
	SingleUpload(content []byte, key string) error
	GetSingleUploadInfo(key string) (map[string]interface{}, error)
	GetDownloadInfo(filename string) (string, error)
	DownloadFile(url string) ([]byte, string, error)
	DeleteFile(filename string) error
}
