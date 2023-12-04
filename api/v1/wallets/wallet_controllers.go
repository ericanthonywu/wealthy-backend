package wallets

import (
	"github.com/gin-gonic/gin"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/wallets/dtos"
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
	//InvestType    string       `json:"invest_type"`
	//InvestName    string       `json:"invest_name"`
	//WalletType    string       `json:"id_master_wallet_type"`
	//WalletAmount  WalletAmount `json:"wallet_amount"`
	//FeeInvestBuy  int64        `json:"fee_invest_buy"`
	//FeeInvestSell int64        `json:"fee_invest_sell"`
	//Amount        int64        `json:"amount"

	if dtoRequest.InvestType == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "no body payload")
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