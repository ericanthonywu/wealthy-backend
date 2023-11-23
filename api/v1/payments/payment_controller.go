package payments

import (
	"github.com/gin-gonic/gin"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/payments/dtos"
	"github.com/semicolon-indonesia/wealthy-backend/utils/errorsinfo"
	"github.com/semicolon-indonesia/wealthy-backend/utils/response"
	"net/http"
)

type (
	PaymentController struct {
		useCase IPaymentUseCase
	}

	IPaymentController interface {
		Subscriptions(ctx *gin.Context)
	}
)

func NewPaymentController(useCase IPaymentUseCase) *PaymentController {
	return &PaymentController{useCase: useCase}
}

func (c *PaymentController) Subscriptions(ctx *gin.Context) {
	var (
		dtoRequest dtos.PaymentSubscription
		errInfo    []errorsinfo.Errors
	)

	if err := ctx.ShouldBindJSON(&dtoRequest); err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "body payload required")
		response.SendBack(ctx, dtos.PaymentSubscription{}, errInfo, http.StatusBadRequest)
		return
	}

	data, httpCode, errInfo := c.useCase.Subscriptions(ctx, &dtoRequest)
	response.SendBack(ctx, data, errInfo, httpCode)
	return

}