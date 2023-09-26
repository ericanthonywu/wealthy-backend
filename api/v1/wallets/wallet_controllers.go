package wallets

import (
	"fmt"
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
	}
)

func NewWalletController(useCase IWalletUseCase) *WalletController {
	return &WalletController{useCase: useCase}
}

func (c *WalletController) Add(ctx *gin.Context) {
	var (
		dtoRequest  dtos.WalletAddRequest
		dtoResponse dtos.WalletAddResponse
		httpCode    int
		errInfo     []errorsinfo.Errors
	)

	if err := ctx.ShouldBindJSON(&dtoRequest); err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "no body payload")
		response.SendBack(ctx, dtoResponse, errInfo, http.StatusBadRequest)
		return
	}

	usrEmail := fmt.Sprintf("%v", ctx.MustGet("email"))
	dtoResponse, httpCode, errInfo = c.useCase.Add(&dtoRequest, usrEmail)

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	response.SendBack(ctx, dtoResponse, errInfo, httpCode)
	return
}

func (c *WalletController) List(ctx *gin.Context) {
	data, httpCode, errInfo := c.useCase.List(ctx)
	response.SendBack(ctx, data, errInfo, httpCode)
	return
}
