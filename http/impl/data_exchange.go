package impl

import "net/http"

type DataExchangeService struct {
	Client *http.Client
}

func (*DataExchangeService) CheckProduction(token, prodId string) (int, error) {
	return 1, nil
}

func (*DataExchangeService) SyncKgMetaData(data map[string]interface{}) error {
	return nil
}

func (*DataExchangeService) Notice(token, prodId string, state int) error {
	return nil
}

func (*DataExchangeService) ImportKg(token, prodId string, startIndex, endIndex int) error {
	return nil
}

func (*DataExchangeService) UploadFile(file []byte, token, prodId, filename, hash string) error {
	return nil
}
