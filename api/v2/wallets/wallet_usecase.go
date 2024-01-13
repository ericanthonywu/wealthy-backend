package wallets

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wealthy-app/wealthy-backend/api/v1/wallets/dtos"
	"github.com/wealthy-app/wealthy-backend/api/v1/wallets/entities"
	"github.com/wealthy-app/wealthy-backend/constants"
	"github.com/wealthy-app/wealthy-backend/utils/errorsinfo"
)

type (
	WalletUseCase struct {
		repo IWalletRepository
	}

	IWalletUseCase interface {
		NewWallet(ctx *gin.Context, request *dtos.WalletAddRequest) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		addNewWallet(accountType string, accountUUID uuid.UUID, walletEntity *entities.WalletEntity) (httpCode int, err error)
		setWalletType(IDMasterWallet string) (walletType string)
		setInitialBalance(request *dtos.WalletAddRequest, walletEntity *entities.WalletEntity, accountUUID uuid.UUID, WalletType string) (err error)
	}
)

func NewWalletUseCase(repo IWalletRepository) *WalletUseCase {
	return &WalletUseCase{
		repo: repo,
	}
}

func (s *WalletUseCase) NewWallet(ctx *gin.Context, request *dtos.WalletAddRequest) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		err        error
		WalletType string
	)

	// get information from context
	accountType := fmt.Sprintf("%v", ctx.MustGet("accountType"))
	accountUUID := ctx.MustGet("accountUUID").(uuid.UUID)

	// determine wallet type from id master wallet
	WalletType = s.setWalletType(request.IDMasterWallet)

	UUIDIDMasterWalletType, err := uuid.Parse(request.IDMasterWallet)
	if err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	// setup wallet model
	walletEntity := entities.WalletEntity{
		ID:                 uuid.New(),
		Active:             true,
		WalletName:         request.WalletName,
		WalletType:         WalletType,
		IDMasterWalletType: UUIDIDMasterWalletType,
		FeeInvestBuy:       request.FeeInvestBuy,
		FeeInvestSell:      request.FeeInvestSell,
		TotalAssets:        request.TotalAsset,
		IDAccount:          accountUUID,
	}

	// add new wallet [ orchestrator function ]
	httpCode, err = s.addNewWallet(accountType, accountUUID, &walletEntity)
	if err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, httpCode, errInfo
	}

	// set initial balance [ orchestrator function ]
	err = s.setInitialBalance(request, &walletEntity, accountUUID, WalletType)
	if err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	resp := struct {
		ID      uuid.UUID `json:"wallet_id"`
		Message string    `json:"message"`
	}{
		ID:      walletEntity.ID,
		Message: "success add new wallet",
	}

	return resp, http.StatusOK, []errorsinfo.Errors{}
}

func (s *WalletUseCase) addNewWallet(accountType string, accountUUID uuid.UUID, walletEntity *entities.WalletEntity) (httpCode int, err error) {
	// account type BASIC
	if accountType == constants.AccountBasic {

		// get total wallet
		totalWallet, httpCode, err := s.repo.NumberOfWalletsByID(accountUUID)
		if err != nil {
			return httpCode, err
		}

		// reach 2 wallets. maximum 2 wallets each account
		if totalWallet == constants.MaxWalletBasic {
			return http.StatusBadRequest, errors.New(constants.ProPlan)
		}

		// add new wallet
		httpCode, err = s.repo.NewWallet(walletEntity)
		if err != nil {
			return httpCode, err
		}
	}

	// account type PRO
	if accountType == constants.AccountPro {
		// add new wallets
		httpCode, err = s.repo.NewWallet(walletEntity)
		if err != nil {
			return httpCode, err
		}
	}

	return http.StatusOK, nil
}

func (s *WalletUseCase) setWalletType(IDMasterWallet string) (walletType string) {
	switch IDMasterWallet {
	case constants.IDCash:
		walletType = constants.Cash
	case constants.IDDebitCard:
		walletType = constants.DebitCard
	case constants.IDCreditCard:
		walletType = constants.CreditCard
	case constants.IDInvestment:
		walletType = constants.Investment
	case constants.IDSaving:
		walletType = constants.Saving
	case constants.IDEWallet:
		walletType = constants.EWallet
	}
	return constants.Cash
}

func (s *WalletUseCase) setInitialBalance(request *dtos.WalletAddRequest, walletEntity *entities.WalletEntity, accountUUID uuid.UUID, WalletType string) (err error) {
	return
}