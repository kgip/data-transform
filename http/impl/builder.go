package impl

import (
	"data-transform/model/to"
	"net/http"
)

type BuilderService struct {
	Client *http.Client
}

func (service *BuilderService) GetKgInfo(kgId string) (*to.KgInfoTo, error) {
	return &to.KgInfoTo{State: 2, MaxIndex: 111}, nil
}

func (service *BuilderService) ExportKgMetaData(kgId string) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}
