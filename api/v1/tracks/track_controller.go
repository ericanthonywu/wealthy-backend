package tracks

import (
	"github.com/gin-gonic/gin"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/tracks/dtos"
	"github.com/semicolon-indonesia/wealthy-backend/utils/errorsinfo"
	"github.com/semicolon-indonesia/wealthy-backend/utils/response"
	"net/http"
)

type (
	TrackController struct {
		useCase ITrackUseCase
	}

	ITrackController interface {
		ScreenTime(ctx *gin.Context)
	}
)

func NewTrackController(useCase ITrackUseCase) *TrackController {
	return &TrackController{useCase: useCase}
}

func (c *TrackController) ScreenTime(ctx *gin.Context) {
	var (
		dtoRequest dtos.ScreenTimeRequest
		errInfo    []errorsinfo.Errors
	)

	// bind
	if err := ctx.ShouldBindJSON(&dtoRequest); err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "body payload required")
		response.SendBack(ctx, dtos.ScreenTimeRequest{}, errInfo, http.StatusBadRequest)
		return
	}

	if dtoRequest.StartTime == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "start_time attribute needed in body payload")
		response.SendBack(ctx, dtos.ScreenTimeRequest{}, errInfo, http.StatusBadRequest)
		return
	}

	if dtoRequest.EndTime == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "end_time attribute needed in body payload")
		response.SendBack(ctx, dtos.ScreenTimeRequest{}, errInfo, http.StatusBadRequest)
		return
	}

	data, httpCode, errInfo := c.useCase.ScreenTime(ctx, &dtoRequest)
	response.SendBack(ctx, data, errInfo, httpCode)
	return
}