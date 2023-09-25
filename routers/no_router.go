package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NoRoute(ctx *gin.Context) {
	ctx.JSON(http.StatusNotFound, "no route found")
}
