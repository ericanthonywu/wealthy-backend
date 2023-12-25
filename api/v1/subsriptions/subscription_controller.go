package subsriptions

import (
	"github.com/gin-gonic/gin"
	"github.com/wealthy-app/wealthy-backend/utils/errorsinfo"
	"github.com/wealthy-app/wealthy-backend/utils/response"
)

type (
	SubscriptionController struct {
		useCase ISubscriptionUseCase
	}

	ISubscriptionController interface {
		Plan(ctx *gin.Context)
		FAQ(ctx *gin.Context)
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

func (c *SubscriptionController) FAQ(ctx *gin.Context) {
	data, httpCode, errInfo := c.useCase.FAQ(ctx)

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	response.SendBack(ctx, data, errInfo, httpCode)
	return
}