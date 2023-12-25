package payments

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/wealthy-app/wealthy-backend/api/v1/payments/dtos"
	"github.com/wealthy-app/wealthy-backend/utils/errorsinfo"
	"github.com/wealthy-app/wealthy-backend/utils/response"
	"github.com/wealthy-app/wealthy-backend/utils/utilities"
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

	// bind
	if err := ctx.ShouldBindJSON(&dtoRequest); err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "body payload required")
		response.SendBack(ctx, dtos.PaymentSubscription{}, errInfo, http.StatusBadRequest)
		return
	}

	// validate
	if dtoRequest.PackageID == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "package id empty value")
	}

	// show error
	if len(errInfo) > 0 {
		response.SendBack(ctx, struct{}{}, errInfo, http.StatusBadRequest)
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

	// bind
	if err := ctx.ShouldBindJSON(&dtoRequest); err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "body payload required")
		response.SendBack(ctx, dtos.PaymentSubscription{}, errInfo, http.StatusBadRequest)
		return
	}

	fmt.Println("%+v", dtoRequest)

	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	sha512, err := utilities.CalculateSHA512(dtoRequest.OrderId + dtoRequest.StatusCode + dtoRequest.GrossAmount + serverKey)
	if err != nil {
		logrus.Error(err.Error())
	}

	if sha512 != dtoRequest.SignatureKey {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "signature not match")
		response.SendBack(ctx, dtos.PaymentSubscription{}, errInfo, http.StatusBadRequest)
		return
	}

	data, httpCode, errInfo := c.useCase.MidtransWebhook(ctx, &dtoRequest)
	response.SendBack(ctx, data, errInfo, httpCode)
	return
}