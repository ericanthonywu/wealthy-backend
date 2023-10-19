package statistics

import (
	"github.com/gin-gonic/gin"
	"github.com/semicolon-indonesia/wealthy-backend/utils/errorsinfo"
	"github.com/semicolon-indonesia/wealthy-backend/utils/response"
)

type (
	StatisticController struct {
		useCase IStatisticUseCase
	}

	IStatisticController interface {
		Statistic(ctx *gin.Context)
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

func (c *StatisticController) TransactionPriority(ctx *gin.Context) {
	data, httpCode, errInfo := c.useCase.TransactionPriority(ctx)

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *StatisticController) Trend(ctx *gin.Context) {
	data, httpCode, errInfo := c.useCase.Trend(ctx)

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
