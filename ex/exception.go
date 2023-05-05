package ex

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
)

func TryThrow(params ...interface{}) {
	if e, success := params[len(params)-1].(*Exception); success {
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

func HandleDbError(db *gorm.DB, exception ...*Exception) error {
	if db.RowsAffected > 0 {
		return nil
	} else if db.Error != nil {
		if len(exception) > 0 {
			return exception[0].SetCause(db.Error.Error())
		} else {
			return DBException.SetCause(db.Error.Error())
		}
	} else {
		if len(exception) > 0 {
			return exception[0]
		}
		return DBException.SetDetail(fmt.Sprintf("execute SQL '%s' affected row 0 ", db.Statement.SQL.String()))
	}
}

type Exception struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Cause  string `json:"cause"`
	Detail string `json:"detail"`
}

func (e *Exception) SetCause(cause string) *Exception {
	returnEx := *e
	returnEx.Cause = cause
	return &returnEx
}

func (e *Exception) SetDetail(detail string) *Exception {
	returnEx := *e
	returnEx.Detail = detail
	return &returnEx
}

func (e *Exception) Error() string {
	return e.Msg
}

func (e *Exception) DetailInfo() string {
	if bytes, err := json.Marshal(e); err != nil {
		panic(err)
	} else {
		return string(bytes)
	}
}

var (
	InternalException            = &Exception{50000, "Server Internal Error!", "", ""}
	DBException                  = &Exception{50001, "Database operation exception!", "", ""}
	RequestParamsException       = &Exception{50002, "Request Params Error!", "", ""}
	LoginException               = &Exception{50003, "User login info Error!", "", ""}
	InsufficientBalanceException = &Exception{50004, "Insufficient user balance!", "", ""}
	PacketIsEmptyException       = &Exception{50005, "The red envelope has been robbed!", "", ""}
	PacketIsExpireException      = &Exception{50006, "The red envelope has been expired!", "", ""}
	GenericRedPacketException    = &Exception{50007, "Failed to generate red envelope. Please Retry!", "", ""}
	RepeatGrabRedPacketException = &Exception{50008, "Repeated grab red envelopes!", "", ""}
)
