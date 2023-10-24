package http

type DataExchangeService interface {
	CheckProduction(token, prodId string) (int, error)
	SyncKgMetaData(data map[string]interface{}) error
	Notice(token, prodId string, state int) error
	ImportKg(token, prodId string, startIndex, endIndex int) error
	UploadFile(file []byte, token, prodId, filename, hash string) error
}
