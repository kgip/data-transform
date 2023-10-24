package initialize

import (
	"crypto/tls"
	"data-transform/api/v1"
	"data-transform/http"
	"data-transform/http/impl"
	"data-transform/middleware"
	"data-transform/service"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	_http "net/http"
	"time"
)

var (
	DataExchangeService http.DataExchangeService
	BuilderService      http.BuilderService
	OssGateway          http.OssGateway

	TaskService *service.TaskService

	TaskController *v1.TaskController
)

func initService() {
	client := &_http.Client{
		Timeout: 5 * time.Second,
		Transport: &_http.Transport{
			MaxConnsPerHost:     100,
			MaxIdleConnsPerHost: 16,
			IdleConnTimeout:     30 * time.Second,
			DisableCompression:  true,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true}},
	}

	DataExchangeService = &impl.DataExchangeService{Client: client}
	BuilderService = &impl.BuilderService{Client: client}
	OssGateway = &impl.OssGateway{}

	TaskService = &service.TaskService{DataExchangeService: DataExchangeService, BuilderService: BuilderService, OssGateway: OssGateway, C: cron.New()}
}

func initController() {
	TaskController = &v1.TaskController{TaskService: TaskService}
}

func init() {
	initService()
	initController()
}

func Router() *gin.Engine {
	router := gin.Default()
	group := router.Group("transform")
	group.Use(middleware.ErrorHandler)
	g1 := group.Group("v1")
	{
		g1.POST("/upload", TaskController.Upload)
	}
	return router
}
