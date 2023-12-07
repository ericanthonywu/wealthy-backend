package wallets

import (
	"github.com/gin-gonic/gin"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/wallets/dtos"
	"github.com/semicolon-indonesia/wealthy-backend/constants"
	"github.com/semicolon-indonesia/wealthy-backend/utils/errorsinfo"
	"github.com/semicolon-indonesia/wealthy-backend/utils/response"
	"net/http"
	"strings"
)

type (
	WalletController struct {
		useCase IWalletUseCase
	}

	IWalletController interface {
		Add(ctx *gin.Context)
		List(ctx *gin.Context)
		UpdateAmount(ctx *gin.Context)
	}
)

func NewWalletController(useCase IWalletUseCase) *WalletController {
	return &WalletController{useCase: useCase}
}

func (c *WalletController) Add(ctx *gin.Context) {
	var (
		dtoRequest dtos.WalletAddRequest
		errInfo    []errorsinfo.Errors
	)

	// bind
	if err := ctx.ShouldBindJSON(&dtoRequest); err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "no body payload")
		response.SendBack(ctx, struct{}{}, errInfo, http.StatusBadRequest)
		return
	}

	// validate
	if dtoRequest.WalletName == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "wallet name empty")
	}

	if dtoRequest.WalletType == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "wallet type empty")
	}

	if dtoRequest.TotalAsset == 0 {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "total assets empty")
	}

	// check wallet types
	isTypeMatch := strings.ToUpper(dtoRequest.WalletType) == constants.Cash || strings.ToUpper(dtoRequest.WalletType) == constants.CreditCard ||
		strings.ToUpper(dtoRequest.WalletType) == constants.DebitCard || strings.ToUpper(dtoRequest.WalletType) == constants.Investment ||
		strings.ToUpper(dtoRequest.WalletType) == constants.Saving

	if !isTypeMatch {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "wallet type must contain one of values. [ CASH, CREDIT_CARD, DEBIT_CARD, INVESTMENT, SAVING ]")
		response.SendBack(ctx, struct{}{}, errInfo, http.StatusBadRequest)
		return
	}

	if dtoRequest.FeeInvestBuy == 0 {
		dtoRequest.FeeInvestSell = 0.15
	}

	if dtoRequest.FeeInvestSell == 0 {
		dtoRequest.FeeInvestSell = 0.25
	}

	if len(errInfo) > 0 {
		response.SendBack(ctx, struct{}{}, errInfo, http.StatusBadRequest)
		return
	}

	dtoResponse, httpCode, errInfo := c.useCase.Add(ctx, &dtoRequest)
	response.SendBack(ctx, dtoResponse, errInfo, httpCode)
	return
}

func (c *WalletController) List(ctx *gin.Context) {
	data, httpCode, errInfo := c.useCase.List(ctx)
	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *WalletController) UpdateAmount(ctx *gin.Context) {
	var (
		dtoRequest  dtos.WalletUpdateAmountRequest
		dtoResponse dtos.WalletUpdateAmountResponse
		httpCode    int
		errInfo     []errorsinfo.Errors
		data        interface{}
	)

	if err := ctx.ShouldBindJSON(&dtoRequest); err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "no body payload")
		response.SendBack(ctx, dtoResponse, errInfo, http.StatusBadRequest)
		return
	}

	walletID := ctx.Param("id-wallet")
	data, httpCode, errInfo = c.useCase.UpdateAmount(walletID, &dtoRequest)

	response.SendBack(ctx, data, errInfo, httpCode)
	return
}