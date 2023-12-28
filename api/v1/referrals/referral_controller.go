package referrals

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/wealthy-app/wealthy-backend/api/v1/referrals/dtos"
	"github.com/wealthy-app/wealthy-backend/utils/errorsinfo"
	"github.com/wealthy-app/wealthy-backend/utils/response"
	"net/http"
)

type (
	ReferralController struct {
		useCase IReferralUseCase
	}

	IReferralController interface {
		Statistic(ctx *gin.Context)
		List(ctx *gin.Context)
		Earn(ctx *gin.Context)
		Withdraw(ctx *gin.Context)
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

func (c *ReferralController) Earn(ctx *gin.Context) {
	data, httpCode, err := c.useCase.Earn(ctx)
	response.SendBack(ctx, data, err, httpCode)
	return
}

func (c *ReferralController) Withdraw(ctx *gin.Context) {
	var (
		request dtos.WithdrawRequest
		errInfo []errorsinfo.Errors
	)

	// bind
	if err := ctx.ShouldBindJSON(&request); err != nil {
		logrus.Error(err.Error())
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "body payload required")
		response.SendBack(ctx, struct{}{}, errInfo, http.StatusBadRequest)
		return
	}

	// validation
	if request.BankIssue == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "bank issue empty value")
	}

	if request.AccountName == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "account name empty value")
	}

	if request.AccountNumber == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "account number must greater than 0")
	}

	if request.WithdrawAmount == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "withdraw amount must greater than 0")
	}

	// send back error
	if len(errInfo) > 0 {
		response.SendBack(ctx, struct{}{}, errInfo, http.StatusBadRequest)
		return
	}

	data, httpCode, err := c.useCase.Withdraw(ctx, request)
	response.SendBack(ctx, data, err, httpCode)
	return
}