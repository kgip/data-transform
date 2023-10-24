package error

import (
	"encoding/json"
)

func TryThrow(params ...interface{}) {
	if e, success := params[len(params)-1].(*CustomError); success {
		for _, param := range params[:len(params)-1] {
			if param == nil {
				continue
			}
			if err, ok := param.(error); ok {
				e.Detail = err.Error()
				panic(e)
			}
		}
		if len(params) <= 1 {
			panic(params[len(params)-1])
		}
	} else {
		for _, param := range params {
			if param == nil {
				continue
			}
			if err, ok := param.(error); ok {
				panic(err)
			}
		}
	}
}

type CustomError struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Cause  string `json:"cause"`
	Detail string `json:"detail"`
}

func (e *CustomError) SetCause(cause string) *CustomError {
	returnEx := *e
	returnEx.Cause = cause
	return &returnEx
}

func (e *CustomError) SetDetail(detail string) *CustomError {
	returnEx := *e
	returnEx.Detail = detail
	return &returnEx
}

func (e *CustomError) Error() string {
	return e.Msg
}

func (e *CustomError) DetailInfo() string {
	if bytes, err := json.Marshal(e); err != nil {
		panic(err)
	} else {
		return string(bytes)
	}
}

var (
	ErrInternal  = &CustomError{50000, "Server internal error!", "", ""}
	ErrDB        = &CustomError{50001, "Database operation error!", "", ""}
	ErrParameter = &CustomError{50002, "Request params error!", "", ""}
)
