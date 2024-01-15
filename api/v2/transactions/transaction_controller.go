package transactions

import (
	"github.com/gin-gonic/gin"
	"github.com/wealthy-app/wealthy-backend/api/v1/transactions/dtos"
	"github.com/wealthy-app/wealthy-backend/constants"
	"github.com/wealthy-app/wealthy-backend/utils/errorsinfo"
	"github.com/wealthy-app/wealthy-backend/utils/personalaccounts"
	"github.com/wealthy-app/wealthy-backend/utils/response"
	"net/http"
)

type (
	TransactionController struct {
		useCase ITransactionUseCase
	}

	ITransactionController interface {
		InvestmentRecords(ctx *gin.Context)
		validateInvestmentRecords(request *dtos.TransactionRequestInvestment) (errInfo []string)
	}
)

func NewTransactionController(useCase ITransactionUseCase) *TransactionController {
	return &TransactionController{useCase: useCase}
}

func (c *TransactionController) InvestmentRecords(ctx *gin.Context) {
	var (
		request dtos.TransactionRequestInvestment
		errInfo []string
	)

	if personalaccounts.PremiumFeature(ctx) {
		resp := struct {
			Message string `json:"message"`
		}{
			Message: constants.ProPlan,
		}
		response.SendBack(ctx, resp, []errorsinfo.Errors{}, http.StatusUpgradeRequired)
		return
	}

	// bind
	if err := ctx.ShouldBindJSON(&request); err != nil {
		errInfo = errorsinfo.ErrorWrapperArray(errInfo, "no body payload")
		response.SendBack(ctx, struct{}{}, errInfo, http.StatusBadRequest)
		return
	}

	errInfo = c.validateInvestmentRecords(&request)

	// if error occurs
	if len(errInfo) > 0 {
		response.SendBack(ctx, struct{}{}, errInfo, http.StatusBadRequest)
		return
	}

	data, httpCode, errInfo := c.useCase.InvestmentRecords(ctx, &request)
	response.SendBack(ctx, data, errInfo, httpCode)
	return

}

func (c *TransactionController) validateInvestmentRecords(request *dtos.TransactionRequestInvestment) (errInfo []string) {

	if request.IDWallet == "" {
		errInfo = errorsinfo.ErrorWrapperArray(errInfo, "id wallet empty value")
	}

	if request.IDMasterInvest == "" {
		errInfo = errorsinfo.ErrorWrapperArray(errInfo, "id categories invest empty value")
	}

	if request.IDMasterTransactionTypes == "" {
		errInfo = errorsinfo.ErrorWrapperArray(errInfo, "id categories transaction type empty value")
	}

	if request.Lot <= 0 {
		errInfo = errorsinfo.ErrorWrapperArray(errInfo, "lot must greater than 0")
	}

	if request.Price <= 0 {
		errInfo = errorsinfo.ErrorWrapperArray(errInfo, "price must greater than 0")
	}

	if request.SellBuy < 0 || request.SellBuy > 2 {
		errInfo = errorsinfo.ErrorWrapperArray(errInfo, "sell buy must one of two values : 0 [ sell ] or 1 [ buy ] ")
	}

	if request.StockCode == "" {
		errInfo = errorsinfo.ErrorWrapperArray(errInfo, "stock code empty value")
	}

	return errInfo
}
