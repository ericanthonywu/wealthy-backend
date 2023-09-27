package wallets

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/wallets/dtos"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/wallets/entities"
	"github.com/semicolon-indonesia/wealthy-backend/utils/errorsinfo"
	"net/http"
)

type (
	WalletUseCase struct {
		repo IWalletRepository
	}

	IWalletUseCase interface {
		Add(request *dtos.WalletAddRequest, usrEmail string) (response dtos.WalletAddResponse, httpCode int, errInfo []errorsinfo.Errors)
		List(ctx *gin.Context) (data interface{}, httpCode int, errInfo []errorsinfo.Errors)
		UpdateAmount(IDWallet string, request *dtos.WalletUpdateAmountRequest) (data interface{}, httpCode int, errInfo []errorsinfo.Errors)
	}
)

func NewWalletUseCase(repo IWalletRepository) *WalletUseCase {
	return &WalletUseCase{repo: repo}
}

func (s *WalletUseCase) Add(request *dtos.WalletAddRequest, usrEmail string) (response dtos.WalletAddResponse, httpCode int, errInfo []errorsinfo.Errors) {
	var err error

	walletEntity := entities.WalletEntity{
		Active:        true,
		InvestType:    request.InvestType,
		InvestName:    request.InvestName,
		WalletType:    request.WalletType,
		FeeInvestBuy:  request.FeeInvestBuy,
		FeeInvestSell: request.FeeInvestBuy,
		Amount:        request.Amount,
	}

	data := s.repo.PersonalAccount(usrEmail)

	if data.ID == uuid.Nil {
		httpCode = http.StatusNotFound
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "not found")
		return response, httpCode, errInfo
	}

	walletEntity.ID = uuid.New()
	walletEntity.IDAccount = data.ID

	response.InvestType = request.InvestType
	response.InvestName = request.InvestName
	response.WalletType = request.WalletType
	response.Amount = request.Amount
	httpCode = http.StatusOK

	if data.TotalWallets < 2 && data.AccountTypes == "BASIC" {
		err = s.repo.Add(&walletEntity)
		if err != nil {
			httpCode = http.StatusInternalServerError
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
			return response, httpCode, errInfo
		}
	}

	if data.AccountTypes == "PRO" {
		err = s.repo.Add(&walletEntity)
		if err != nil {
			httpCode = http.StatusInternalServerError
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
			return response, httpCode, errInfo
		}
	}

	if data.TotalWallets == 2 && data.AccountTypes == "BASIC" {
		httpCode = http.StatusUnprocessableEntity
		response = dtos.WalletAddResponse{}
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "can not add new wallet. please upgrade to PRO subscription")
	}

	return response, httpCode, errInfo
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
