package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (handler *Handler) HandleGetValue(ctx *gin.Context) {
	value, err := handler.service.GetValue(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"value": value})
}
