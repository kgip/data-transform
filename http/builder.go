package http

import "data-transform/model/to"

type BuilderService interface {
	GetKgInfo(kgId string) (*to.KgInfoTo, error)
	ExportKgMetaData(kgId string) (map[string]interface{}, error)
}
