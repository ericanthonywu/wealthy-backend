package investments

import (
	"github.com/gin-gonic/gin"
	"github.com/wealthy-app/wealthy-backend/constants"
	"github.com/wealthy-app/wealthy-backend/utils/errorsinfo"
	"github.com/wealthy-app/wealthy-backend/utils/response"
	"net/http"
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
	// get account type
	accountType := ctx.MustGet("accountType").(string)

	// if basic account
	if accountType == constants.AccountBasic {
		resp := struct {
			Message string `json:"message"`
		}{
			Message: constants.ProPlan,
		}
		response.SendBack(ctx, resp, []errorsinfo.Errors{}, http.StatusUpgradeRequired)
		return
	}

	data, httpCode, errInfo := c.useCase.Portfolio(ctx)
	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *InvestmentController) GainLoss(ctx *gin.Context) {
	// get account type
	accountType := ctx.MustGet("accountType").(string)

	// if basic account
	if accountType == constants.AccountBasic {
		resp := struct {
			Message string `json:"message"`
		}{
			Message: constants.ProPlan,
		}
		response.SendBack(ctx, resp, []errorsinfo.Errors{}, http.StatusUpgradeRequired)
		return
	}

	data, httpCode, errInfo := c.useCase.GainLoss(ctx)
	response.SendBack(ctx, data, errInfo, httpCode)
	return
}