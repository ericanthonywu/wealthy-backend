package routers

import (
	"github.com/semicolon-indonesia/wealthy-backend/utils/errorsinfo"
	"github.com/semicolon-indonesia/wealthy-backend/utils/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NoRoute(ctx *gin.Context) {
	var errInfo []errorsinfo.Errors

	resp := struct {
		Message string `json:"message,omitempty"`
	}{
		Message: "no routes or wrong method",
	}

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	response.SendBack(ctx, resp, errInfo, http.StatusNotFound)
	return
}