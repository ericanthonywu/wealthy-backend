package statistics

import (
	"github.com/gin-gonic/gin"
	"github.com/semicolon-indonesia/wealthy-backend/utils/errorsinfo"
	"github.com/semicolon-indonesia/wealthy-backend/utils/response"
	"net/http"
)

type (
	StatisticController struct {
		useCase IStatisticUseCase
	}

	IStatisticController interface {
		Statistic(ctx *gin.Context)
		Weekly(ctx *gin.Context)
		Summary(ctx *gin.Context)
		TransactionPriority(ctx *gin.Context)
		Trend(ctx *gin.Context)
		Category(ctx *gin.Context)
	}
)

func NewStatisticController(useCase IStatisticUseCase) *StatisticController {
	return &StatisticController{useCase: useCase}
}

func (c *StatisticController) Statistic(ctx *gin.Context) {
	data, httpCode, errInfo := c.useCase.Statistic(ctx)

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *StatisticController) Weekly(ctx *gin.Context) {
	var (
		errInfo    []errorsinfo.Errors
		data       interface{}
		statusCode int
	)
	month := ctx.Query("month")
	year := ctx.Query("year")

	if month == "" || year == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "month or year required in query url")
		response.SendBack(ctx, nil, errInfo, http.StatusBadRequest)
		return
	}

	data, statusCode, errInfo = c.useCase.Weekly(ctx, month, year)

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	response.SendBack(ctx, data, errInfo, statusCode)
	return
}

func (c *StatisticController) Summary(ctx *gin.Context) {
	var (
		errInfo    []errorsinfo.Errors
		data       interface{}
		statusCode int
	)
	month := ctx.Query("month")
	year := ctx.Query("year")

	if month == "" || year == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "month or year required in query url")
		response.SendBack(ctx, nil, errInfo, http.StatusBadRequest)
		return
	}

	data, statusCode, errInfo = c.useCase.Summary(ctx, month, year)

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	response.SendBack(ctx, data, errInfo, statusCode)
	return
}

func (c *StatisticController) TransactionPriority(ctx *gin.Context) {
	var (
		errInfo []errorsinfo.Errors
	)

	month := ctx.Query("month")
	year := ctx.Query("year")

	if month == "" || year == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "month or year required in query url")
		response.SendBack(ctx, nil, errInfo, http.StatusBadRequest)
		return
	}

	data, httpCode, errInfo := c.useCase.Priority(ctx, month, year)

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *StatisticController) Trend(ctx *gin.Context) {
	var (
		errInfo []errorsinfo.Errors
	)

	month := ctx.Query("month")
	year := ctx.Query("year")

	if month == "" || year == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "month or year required in query url")
		response.SendBack(ctx, nil, errInfo, http.StatusBadRequest)
		return
	}

	data, httpCode, errInfo := c.useCase.Trend(ctx, month, year)

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *StatisticController) Category(ctx *gin.Context) {
	data, httpCode, errInfo := c.useCase.Category(ctx)

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	response.SendBack(ctx, data, errInfo, httpCode)
	return
}