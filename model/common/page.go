package common

type Page struct {
	Page  int         `form:"page";json:"page"`
	Limit int         `form:"limit";json:"limit"`
	Total int         `form:"total";json:"total"`
	Data  interface{} `form:"data";json:"data"`
}
