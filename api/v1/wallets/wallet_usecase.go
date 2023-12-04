package wallets

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/wallets/dtos"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/wallets/entities"
	"github.com/semicolon-indonesia/wealthy-backend/constants"
	"github.com/semicolon-indonesia/wealthy-backend/utils/errorsinfo"
	"github.com/semicolon-indonesia/wealthy-backend/utils/personalaccounts"
	"net/http"
)

type (
	WalletUseCase struct {
		repo IWalletRepository
	}

	IWalletUseCase interface {
		Add(ctx *gin.Context, request *dtos.WalletAddRequest) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		List(ctx *gin.Context) (data interface{}, httpCode int, errInfo []errorsinfo.Errors)
		UpdateAmount(IDWallet string, request *dtos.WalletUpdateAmountRequest) (data interface{}, httpCode int, errInfo []errorsinfo.Errors)
	}
)

func NewWalletUseCase(repo IWalletRepository) *WalletUseCase {
	return &WalletUseCase{repo: repo}
}

func (s *WalletUseCase) Add(ctx *gin.Context, request *dtos.WalletAddRequest) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		err         error
		dtoResponse dtos.WalletAddResponse
	)

	// setup wallet model
	walletEntity := entities.WalletEntity{
		Active:        true,
		InvestType:    request.InvestType,
		InvestName:    request.InvestName,
		WalletType:    request.WalletType,
		FeeInvestBuy:  request.FeeInvestBuy,
		FeeInvestSell: request.FeeInvestBuy,
		Amount:        request.Amount,
	}

	usrEmail := fmt.Sprintf("%v", ctx.MustGet("email"))
	data := personalaccounts.Informations(ctx, usrEmail)

	if data.ID == uuid.Nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", constants.TokenInvalidInformation)
		return struct{}{}, http.StatusNotFound, errInfo
	}

	walletEntity.ID = uuid.New()
	walletEntity.IDAccount = data.ID

	dtoResponse.InvestType = request.InvestType
	dtoResponse.InvestName = request.InvestName
	dtoResponse.WalletType = request.WalletType
	dtoResponse.Amount = request.Amount

	// account type BASIC
	if data.TotalWallets < 2 && data.AccountTypes == "BASIC" {

		// add new wallets
		err = s.repo.Add(&walletEntity)
		if err != nil {
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
			return struct{}{}, http.StatusInternalServerError, errInfo
		}

		// record on transaction as initial
	}

	// account type BASIC reach limit
	if data.TotalWallets == 2 && data.AccountTypes == "BASIC" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "can not add new wallet. please upgrade to PRO subscription")
		return struct{}{}, http.StatusUnprocessableEntity, errInfo
	}

	// account type PRO
	if data.AccountTypes == "PRO" {
		// add new wallets
		err = s.repo.Add(&walletEntity)
		if err != nil {
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
			return struct{}{}, http.StatusInternalServerError, errInfo
		}

		// record on transaction as initial
	}

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	return response, http.StatusOK, errInfo
}

func (s *WalletUseCase) List(ctx *gin.Context) (data interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		err   error
		email string
	)

	email = fmt.Sprintf("%v", ctx.MustGet("email"))
	data, httpCode, err = s.repo.List(email)
	if err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
	}

	if err == nil {
		errInfo = []errorsinfo.Errors{}
	}

	return data, httpCode, errInfo
}

func (s *WalletUseCase) UpdateAmount(IDWallet string, request *dtos.WalletUpdateAmountRequest) (data interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var err error

	data, httpCode, err = s.repo.UpdateAmount(IDWallet, request.Amount)
	if err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return data, httpCode, errInfo
	}

	errInfo = []errorsinfo.Errors{}
	return data, httpCode, errInfo
}