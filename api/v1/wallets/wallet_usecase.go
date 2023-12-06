package wallets

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/wallets/dtos"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/wallets/entities"
	"github.com/semicolon-indonesia/wealthy-backend/constants"
	"github.com/semicolon-indonesia/wealthy-backend/utils/datecustoms"
	"github.com/semicolon-indonesia/wealthy-backend/utils/errorsinfo"
	"github.com/semicolon-indonesia/wealthy-backend/utils/personalaccounts"
	"github.com/sirupsen/logrus"
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
		err                error
		IDMasterWalletType string
	)

	switch request.WalletType {
	case constants.Cash:
		IDMasterWalletType = constants.IDCash
	case constants.DebitCard:
		IDMasterWalletType = constants.IDDebitCard
	case constants.CreditCard:
		IDMasterWalletType = constants.IDCreditCard
	case constants.Investment:
		IDMasterWalletType = constants.IDInvestment
	case constants.Saving:
		IDMasterWalletType = constants.IDSaving
	}

	UUIDIDMasterWalletType, err := uuid.Parse(IDMasterWalletType)
	if err != nil {
		logrus.Error(err.Error())
	}

	// setup wallet model
	walletEntity := entities.WalletEntity{
		Active:             true,
		WalletName:         request.WalletName,
		WalletType:         request.WalletType,
		IDMasterWalletType: UUIDIDMasterWalletType,
		FeeInvestBuy:       request.FeeInvestBuy,
		FeeInvestSell:      request.FeeInvestBuy,
		TotalAssets:        request.TotalAsset,
	}

	usrEmail := fmt.Sprintf("%v", ctx.MustGet("email"))
	data := personalaccounts.Informations(ctx, usrEmail)

	if data.ID == uuid.Nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", constants.TokenInvalidInformation)
		return struct{}{}, http.StatusUnauthorized, errInfo
	}

	walletEntity.ID = uuid.New()
	walletEntity.IDAccount = data.ID

	// setup id_master_income_category
	incomeCategoryUUID, err := uuid.Parse("13c4a525-2200-497b-af4b-ef8fa2fe93cc")
	if err != nil {
		logrus.Error(err.Error())
	}

	// setup id_master_transaction_priority
	trxPriorityUUID, err := uuid.Parse("9b96cdf8-8173-4d54-9142-e6ebd1f6aea3")
	if err != nil {
		logrus.Error(err.Error())
	}

	// setup id_master_transaction_type
	trxTypeUUID, err := uuid.Parse("c023a068-a239-42cd-b03a-70304f55d0d3")
	if err != nil {
		logrus.Error(err.Error())
	}

	// setup trx
	trx := entities.WalletInitTransaction{
		ID:                            uuid.New(),
		Date:                          datecustoms.NowTransaction(),
		Fees:                          0,
		Amount:                        float64(request.TotalAsset),
		IDPersonalAccount:             data.ID,
		IDWallet:                      walletEntity.ID,
		IDMasterIncomeCategories:      incomeCategoryUUID,
		IDMasterTransactionPriorities: trxPriorityUUID,
		IDMasterTransactionTypes:      trxTypeUUID,
		Credit:                        float64(request.TotalAsset),
		Debit:                         0,
		Balance:                       float64(request.TotalAsset),
	}

	// no setup trx detail
	trxDetail := entities.WalletInitTransactionDetail{}

	// account type BASIC
	if data.TotalWallets < 2 && data.AccountTypes == "BASIC" {
		// add new wallets
		err = s.repo.Add(&walletEntity)
		if err != nil {
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
			return struct{}{}, http.StatusInternalServerError, errInfo
		}

		// record on transaction as initial
		err = s.repo.InitTransaction(&trx, &trxDetail)
		if err != nil {
			logrus.Error(err.Error())
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
			return struct{}{}, http.StatusInternalServerError, errInfo
		}
	}

	// account type BASIC reach limit
	if data.TotalWallets == 2 && data.AccountTypes == "BASIC" {
		resp := struct {
			Message string `json:"message"`
		}{
			Message: "can not add new wallet. please upgrade to PRO subscription",
		}
		return resp, http.StatusUnprocessableEntity, []errorsinfo.Errors{}
	}

	// account type PRO
	if data.AccountTypes == "PRO" {
		// add new wallets
		err = s.repo.Add(&walletEntity)
		if err != nil {
			logrus.Error(err.Error())
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
			return struct{}{}, http.StatusInternalServerError, errInfo
		}

		// record on transaction as initial
		err = s.repo.InitTransaction(&trx, &trxDetail)
		if err != nil {
			logrus.Error(err.Error())
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
			return struct{}{}, http.StatusInternalServerError, errInfo
		}
	}

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	resp := struct {
		ID      uuid.UUID `json:"wallet_id"`
		Message string    `json:"message"`
	}{
		ID:      walletEntity.ID,
		Message: "success add new wallet",
	}

	return resp, http.StatusOK, errInfo
}

func (s *WalletUseCase) List(ctx *gin.Context) (data interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		err         error
		email       string
		dataList    []entities.WalletEntity
		dtoResponse []dtos.WalletListResponse
	)

	email = fmt.Sprintf("%v", ctx.MustGet("email"))
	personalData := personalaccounts.Informations(ctx, email)

	if personalData.ID == uuid.Nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", constants.TokenInvalidInformation)
		return struct{}{}, http.StatusUnauthorized, errInfo
	}

	// get data wallet
	dataList, err = s.repo.List(personalData.ID)
	if err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	// if not found
	if len(dataList) == 0 {
		resp := struct {
			Message string `json:"message"`
		}{
			Message: "no data wallet",
		}
		return resp, http.StatusNotFound, []errorsinfo.Errors{}
	}

	for _, v := range dataList {
		dataTrx, err := s.repo.LatestAmountWalletInTransaction(v.ID)
		if err != nil {
			logrus.Error(err.Error())
		}

		dtoResponse = append(dtoResponse, dtos.WalletListResponse{
			IDAccount: v.IDAccount,
			WalletDetails: dtos.WalletDetails{
				WalletID:           v.ID,
				WalletType:         v.WalletType,
				WalletName:         v.WalletName,
				IDMasterWalletType: v.IDMasterWalletType,
			},
			Active:        v.Active,
			FeeInvestBuy:  v.FeeInvestBuy,
			FeeInvestSell: v.FeeInvestSell,
			TotalAssets:   int64(dataTrx.Balance),
		})
	}

	return dtoResponse, http.StatusOK, []errorsinfo.Errors{}
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