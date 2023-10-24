package v1

import (
	"data-transform/model/vo"
	"data-transform/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TaskController struct {
	TaskService *service.TaskService
}

func (controller *TaskController) Upload(c *gin.Context) {
	uploadVo := &vo.AddUploadTaskVo{}
	err := c.ShouldBind(uploadVo)
	if err != nil {
		panic(err)
	}
	controller.TaskService.Upload(uploadVo.KgId)
	c.JSON(http.StatusOK, "ok")
}

func (controller *TaskController) ImportKg(c *gin.Context) {
	importKgVo := &vo.ImportKgVo{}
	if err := c.ShouldBind(importKgVo); err != nil {
		panic(err)
	}
	controller.TaskService.ImportKg(importKgVo)
	c.JSON(http.StatusOK, "ok")
}
