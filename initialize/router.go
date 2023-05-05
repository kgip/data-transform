package initialize

import (
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	router := gin.Default()
	group := router.Group("data-transform")

	v1 := group.Group("v1")
	{
		v1.GET("/")
	}
	return router
}
