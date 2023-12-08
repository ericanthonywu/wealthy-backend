package investments

import (
	"github.com/gin-gonic/gin"
	"github.com/semicolon-indonesia/wealthy-backend/utils/response"
)

type (
	InvestmentController struct {
		useCase IInvestmentUseCase
	}

	IInvestmentController interface {
		Portfolio(ctx *gin.Context)
		GainLoss(ctx *gin.Context)
	}
)

func NewInvestmentController(useCase IInvestmentUseCase) *InvestmentController {
	return &InvestmentController{useCase: useCase}
}

func (c *InvestmentController) Portfolio(ctx *gin.Context) {
	data, httpCode, errInfo := c.useCase.Portfolio(ctx)
	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *InvestmentController) GainLoss(ctx *gin.Context) {
	data, httpCode, errInfo := c.useCase.GainLoss(ctx)
	response.SendBack(ctx, data, errInfo, httpCode)
	return
}