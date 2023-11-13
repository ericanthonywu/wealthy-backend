package subsriptions

import (
	"github.com/gin-gonic/gin"
	"github.com/semicolon-indonesia/wealthy-backend/utils/errorsinfo"
	"github.com/semicolon-indonesia/wealthy-backend/utils/response"
)

type (
	SubscriptionController struct {
		useCase ISubscriptionUseCase
	}

	ISubscriptionController interface {
		Plan(ctx *gin.Context)
	}
)

func NewSubscriptionController(useCase ISubscriptionUseCase) *SubscriptionController {
	return &SubscriptionController{useCase: useCase}
}

func (c *SubscriptionController) Plan(ctx *gin.Context) {
	data, httpCode, errInfo := c.useCase.Plan(ctx)

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	response.SendBack(ctx, data, errInfo, httpCode)
	return
}