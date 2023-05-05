package middleware

import (
	"data-transform/ex"
	"github.com/gin-gonic/gin"
)

func LoginHandler(c *gin.Context) {
	id := c.GetHeader("userId")
	if id != "" {
		c.Set("userId", id)
	} else {
		panic(ex.LoginException)
	}
	c.Next()
}
