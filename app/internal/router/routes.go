package router

import (
	"app/internal/api"
	"app/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}

func (h *Handler) initializeRoutes(router *gin.Engine, handler *api.Handler) {
	basePath := "/value"
	route := router.Group(basePath)
	{
		route.POST("/set", handler.HandleSetValue)
		route.GET("/get", handler.HandleGetValue)
		route.POST("/sync", handler.HandleSyncValue)
		route.GET("/check", handler.HandleCheckValue)
	}
}
