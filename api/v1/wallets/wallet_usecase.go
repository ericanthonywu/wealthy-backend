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
	"strings"
)

type (
	WalletUseCase struct {
		repo IWalletRepository
	}

	IWalletUseCase interface {
		Add(ctx *gin.Context, request *dtos.WalletAddRequest) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		List(ctx *gin.Context) (data interface{}, httpCode int, errInfo []errorsinfo.Errors)
		UpdateAmount(IDWallet string, request *dtos.WalletUpdateAmountRequest) (data interface{}, httpCode int, errInfo []errorsinfo.Errors)
		writeInitialTransaction(request *dtos.WalletAddRequest, walletEntity *entities.WalletEntity, data *personalaccounts.PersonalAccountEntities) (err error)
	}
)

func NewWalletUseCase(repo IWalletRepository) *WalletUseCase {
	return &WalletUseCase{repo: repo}
}

func (s *WalletUseCase) Add(ctx *gin.Context, request *dtos.WalletAddRequest) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		accountUUID        uuid.UUID
		err                error
		IDMasterWalletType string
	)

	// get information from context
	accountType := fmt.Sprintf("%v", ctx.MustGet("accountType"))
	accountID := fmt.Sprintf("%v", ctx.MustGet("accountID"))

	if accountID != "" {
		accountUUID, err = uuid.Parse(accountID)
		if err != nil {
			logrus.Error(err.Error())
		}
	}

	// mapping wallet type with wallet id
	switch strings.ToUpper(request.WalletType) {
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
		ID:                 uuid.New(),
		Active:             true,
		WalletName:         request.WalletName,
		WalletType:         strings.ToUpper(request.WalletType),
		IDMasterWalletType: UUIDIDMasterWalletType,
		FeeInvestBuy:       request.FeeInvestBuy,
		FeeInvestSell:      request.FeeInvestSell,
		TotalAssets:        request.TotalAsset,
		IDAccount:          accountUUID,
	}

	// account type BASIC
	if accountType == constants.AccountBasic {
		// get total wallet
		totalWallet, err := s.repo.TotalWallet(accountUUID)
		if err != nil {
			logrus.Error(err)
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
			return struct{}{}, http.StatusInternalServerError, errInfo
		}

		// reach max
		if totalWallet == constants.MaxWalletBasic {
			resp := struct {
				Message string `json:"message"`
			}{
				Message: "can not add new wallet. please upgrade to PRO subscription",
			}
			return resp, http.StatusUnprocessableEntity, []errorsinfo.Errors{}
		}

		err = s.repo.Add(&walletEntity)
		if err != nil {
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
			return struct{}{}, http.StatusInternalServerError, errInfo
		}

		if strings.ToUpper(request.WalletType) != constants.Investment {
			// save initial transaction
			err = s.writeInitialTransaction(request, &walletEntity, accountUUID)
			if err != nil {
				logrus.Error(err.Error())
				errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
				return struct{}{}, http.StatusInternalServerError, errInfo
			}
		}

	}

	// account type PRO
	if accountType == constants.AccountPro {
		// add new wallets
		err = s.repo.Add(&walletEntity)
		if err != nil {
			logrus.Error(err.Error())
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
			return struct{}{}, http.StatusInternalServerError, errInfo
		}

		if strings.ToUpper(request.WalletType) != constants.Investment {
			// save initial transaction
			err = s.writeInitialTransaction(request, &walletEntity, accountUUID)
			if err != nil {
				logrus.Error(err.Error())
				errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
				return struct{}{}, http.StatusInternalServerError, errInfo
			}
		}
	}

	// if no error message
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
		//var totalAsset int64

		// if wallet type is not investments
		if strings.ToUpper(v.WalletType) != constants.Investment {

			// fetch data from transaction latest row to get balance information
			//dataTrx, err := s.repo.LatestAmountWalletInTransaction(v.ID)
			//if err != nil {
			//	logrus.Error(err.Error())
			//}

			//totalAsset = int64(dataTrx.Balance)
		}

		// if wallet type is investments
		if strings.ToUpper(v.WalletType) == constants.Investment {
			//totalAsset = v.TotalAssets
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
			TotalAssets:   v.TotalAssets,
		})
	}

	return dtoResponse, http.StatusOK, []errorsinfo.Errors{}
}

func (s *WalletUseCase) UpdateAmount(IDWallet string, request *dtos.WalletUpdateAmountRequest) (data interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		err        error
		UUIDWallet uuid.UUID
	)

	UUIDWallet, err = uuid.Parse(IDWallet)
	if err != nil {
		logrus.Error(err.Error())
	}

	// update amount
	if request.Amount != 0 {
		data, httpCode, err = s.repo.UpdateAmount(UUIDWallet, request.Amount)
		if err != nil {
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
			return data, http.StatusInternalServerError, errInfo
		}
	}

	// update wallet name
	if request.WalletName != "" {
		err = s.repo.UpdateWalletName(UUIDWallet, request.WalletName)
		if err != nil {
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
			return data, http.StatusInternalServerError, errInfo
		}
	}

	// if e
	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	resp := struct {
		Message string `json:"message"`
	}{
		Message: "update wallet success",
	}

	return resp, http.StatusOK, errInfo
}

func (s *WalletUseCase) writeInitialTransaction(request *dtos.WalletAddRequest, walletEntity *entities.WalletEntity, IDPersonalAccount uuid.UUID) (err error) {
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
		IDPersonalAccount:             IDPersonalAccount,
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

	// record on transaction as initial
	err = s.repo.InitTransaction(&trx, &trxDetail)
	if err != nil {
		return err
	}

	return nil
}