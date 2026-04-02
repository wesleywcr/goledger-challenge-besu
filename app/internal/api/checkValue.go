package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (handler *Handler) HandleCheckValue(ctx *gin.Context) {
	isEqual, blockchainValue, databaseValue, err := handler.service.CheckValue(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"equal":            isEqual,
		"blockchain_value": blockchainValue,
		"database_value":   databaseValue,
	})
}
