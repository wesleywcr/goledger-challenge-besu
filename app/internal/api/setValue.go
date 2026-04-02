package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (handler *Handler) HandleSetValue(ctx *gin.Context) {
	var request setValueRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := handler.service.SetValue(ctx.Request.Context(), request.Value); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "ok", "value": request.Value})
}
