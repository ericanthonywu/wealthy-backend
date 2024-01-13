package wallets

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wealthy-app/wealthy-backend/api/v1/wallets/dtos"
	"github.com/wealthy-app/wealthy-backend/constants"
	"github.com/wealthy-app/wealthy-backend/utils/errorsinfo"
	"github.com/wealthy-app/wealthy-backend/utils/response"
)

type (
	WalletController struct {
		useCase IWalletUseCase
	}

	IWalletController interface {
		NewWallet(ctx *gin.Context)
		GetAllWallets(ctx *gin.Context)
		UpdateWallet(ctx *gin.Context)
	}
)

func NewWalletController(useCase IWalletUseCase) *WalletController {
	return &WalletController{
		useCase: useCase,
	}
}

func (c *WalletController) NewWallet(ctx *gin.Context) {
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

	if dtoRequest.IDMasterWallet == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id categories wallet  empty")
	}

	if dtoRequest.TotalAsset == 0 {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "total assets empty")
	}

	// check wallet id
	isValid := dtoRequest.IDMasterWallet == constants.IDCash ||
		dtoRequest.IDMasterWallet == constants.IDCreditCard ||
		dtoRequest.IDMasterWallet == constants.IDDebitCard ||
		dtoRequest.IDMasterWallet == constants.IDInvestment ||
		dtoRequest.IDMasterWallet == constants.IDSaving ||
		dtoRequest.IDMasterWallet == constants.IDEWallet

	if !isValid {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id categories wallet unregistered")
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

	dtoResponse, httpCode, errInfo := c.useCase.NewWallet(ctx, &dtoRequest)
	response.SendBack(ctx, dtoResponse, errInfo, httpCode)
	return

}

func (c *WalletController) GetAllWallets(ctx *gin.Context) {

}

func (c *WalletController) UpdateWallet(ctx *gin.Context) {

}