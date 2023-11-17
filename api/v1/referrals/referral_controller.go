package referrals

import (
	"github.com/gin-gonic/gin"
	"github.com/semicolon-indonesia/wealthy-backend/utils/response"
)

type (
	ReferralController struct {
		useCase IReferralUseCase
	}

	IReferralController interface {
		Statistic(ctx *gin.Context)
		List(ctx *gin.Context)
	}
)

func NewReferralController(useCase IReferralUseCase) *ReferralController {
	return &ReferralController{useCase: useCase}
}

func (c *ReferralController) Statistic(ctx *gin.Context) {
	data, httpCode, err := c.useCase.Statistic(ctx)
	response.SendBack(ctx, data, err, httpCode)
	return
}

func (c *ReferralController) List(ctx *gin.Context) {
	data, httpCode, err := c.useCase.List(ctx)
	response.SendBack(ctx, data, err, httpCode)
	return
}