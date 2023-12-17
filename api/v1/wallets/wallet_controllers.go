package wallets

import (
	"github.com/gin-gonic/gin"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/wallets/dtos"
	"github.com/semicolon-indonesia/wealthy-backend/constants"
	"github.com/semicolon-indonesia/wealthy-backend/utils/errorsinfo"
	"github.com/semicolon-indonesia/wealthy-backend/utils/response"
	"net/http"
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

	if dtoRequest.WalletID == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "wallet id empty")
	}

	if dtoRequest.TotalAsset == 0 {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "total assets empty")
	}

	// check wallet id
	isValid := dtoRequest.WalletID == constants.IDCash ||
		dtoRequest.WalletID == constants.IDCreditCard ||
		dtoRequest.WalletID == constants.IDDebitCard ||
		dtoRequest.WalletID == constants.IDInvestment ||
		dtoRequest.WalletID == constants.IDSaving

	if !isValid {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id master wallet unregistered")
		response.SendBack(ctx, struct{}{}, errInfo, http.StatusBadRequest)
		return
	}

	//isTypeMatch := strings.ToUpper(dtoRequest.WalletType) == constants.Cash || strings.ToUpper(dtoRequest.WalletType) == constants.CreditCard ||
	//	strings.ToUpper(dtoRequest.WalletType) == constants.DebitCard || strings.ToUpper(dtoRequest.WalletType) == constants.Investment ||
	//	strings.ToUpper(dtoRequest.WalletType) == constants.Saving

	//if !isTypeMatch {
	//	var builder strings.Builder
	//	builder.WriteString("wallet type must contain one of values ")
	//	builder.WriteString("[")
	//	builder.WriteString(constants.Cash)
	//	builder.WriteString(",")
	//	builder.WriteString(constants.CreditCard)
	//	builder.WriteString(",")
	//	builder.WriteString(constants.DebitCard)
	//	builder.WriteString(",")
	//	builder.WriteString(constants.Investment)
	//	builder.WriteString(",")
	//	builder.WriteString(constants.Saving)
	//	builder.WriteString("]")
	//
	//	errInfo = errorsinfo.ErrorWrapper(errInfo, "", builder.String())
	//	response.SendBack(ctx, struct{}{}, errInfo, http.StatusBadRequest)
	//	return
	//}

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
		dtoRequest dtos.WalletUpdateAmountRequest
		httpCode   int
		errInfo    []errorsinfo.Errors
		data       interface{}
	)

	// binding
	if err := ctx.ShouldBindJSON(&dtoRequest); err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "no body payload")
		response.SendBack(ctx, struct{}{}, errInfo, http.StatusBadRequest)
		return
	}

	walletID := ctx.Param("id-wallet")

	if walletID == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id wallet required in query url")
	}

	if len(errInfo) > 0 {
		response.SendBack(ctx, struct{}{}, errInfo, http.StatusBadRequest)
		return
	}

	data, httpCode, errInfo = c.useCase.UpdateAmount(walletID, &dtoRequest)
	response.SendBack(ctx, data, errInfo, httpCode)
	return
}