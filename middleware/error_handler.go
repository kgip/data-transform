package middleware

import (
	"data-transform/ex"
	"data-transform/global"
	"data-transform/model/common"
	"github.com/gin-gonic/gin"
)

func ErrorHandler(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			if e, ok := err.(*ex.Exception); ok {
				global.LOG.Error(e.DetailInfo())
				common.FailWithMessageCode(e.Error(), e.Code, c)
			} else if e, ok := err.(error); ok {
				global.LOG.Error(e.Error())
				common.FailWithMessage(e.Error(), c)
			} else if msg, ok := err.(string); ok {
				global.LOG.Error(msg)
				common.FailWithMessage(msg, c)
			} else {
				global.LOG.Error(ex.InternalException.Error())
				common.FailWithMessage(ex.InternalException.Error(), c)
			}
			c.Abort()
		}
	}()
	c.Next()
}
