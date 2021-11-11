package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Index 首页
func Index(ctx *gin.Context) {
	ctx.String(http.StatusOK, "ok")
}
