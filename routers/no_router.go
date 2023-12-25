package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/wealthy-app/wealthy-backend/utils/errorsinfo"
	"github.com/wealthy-app/wealthy-backend/utils/response"
	"net/http"
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