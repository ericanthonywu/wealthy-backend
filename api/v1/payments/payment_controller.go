package payments

import (
	"github.com/gin-gonic/gin"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/payments/dtos"
	"github.com/semicolon-indonesia/wealthy-backend/utils/errorsinfo"
	"github.com/semicolon-indonesia/wealthy-backend/utils/response"
	"github.com/semicolon-indonesia/wealthy-backend/utils/utilities"
	"net/http"
	"os"
)

type (
	PaymentController struct {
		useCase IPaymentUseCase
	}

	IPaymentController interface {
		Subscriptions(ctx *gin.Context)
		MidtransWebhook(ctx *gin.Context)
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

func (c *PaymentController) MidtransWebhook(ctx *gin.Context) {
	var (
		dtoRequest dtos.MidTransWebhook
		errInfo    []errorsinfo.Errors
	)

	if err := ctx.ShouldBindJSON(&dtoRequest); err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "body payload required")
		response.SendBack(ctx, dtos.PaymentSubscription{}, errInfo, http.StatusBadRequest)
		return
	}

	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	sha512, _ := utilities.CalculateSHA512(dtoRequest.OrderId + dtoRequest.StatusCode + dtoRequest.GrossAmount + serverKey)

	if sha512 != dtoRequest.SignatureKey {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "signature not match")
		response.SendBack(ctx, dtos.PaymentSubscription{}, errInfo, http.StatusBadRequest)
		return
	}

	data, httpCode, errInfo := c.useCase.MidtransWebhook(ctx, &dtoRequest)
	response.SendBack(ctx, data, errInfo, httpCode)
	return
}