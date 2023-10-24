package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ErrorHandler(c *gin.Context) {
	defer func() {
		if e := recover(); e != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, e)
		}
	}()
	c.Next()
}
