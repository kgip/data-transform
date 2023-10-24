package vo

type AddUploadTaskVo struct {
	KgId string `json:"kgId" binding:"required"`
}

type ImportKgVo struct {
	ProdId     string `json:"prodId" binding:"required"`
	StartIndex int    `json:"startIndex" binding:"gte=0"`
	EndIndex   int    `json:"endIndex" binding:"gte=0"`
}
