package wallets

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wealthy-app/wealthy-backend/api/v1/wallets/dtos"
	"github.com/wealthy-app/wealthy-backend/api/v1/wallets/entities"
	"github.com/wealthy-app/wealthy-backend/constants"
	"github.com/wealthy-app/wealthy-backend/utils/datecustoms"
	"github.com/wealthy-app/wealthy-backend/utils/errorsinfo"
	"github.com/wealthy-app/wealthy-backend/utils/personalaccounts"
)

type (
	WalletUseCase struct {
		repo IWalletRepository
	}

	IWalletUseCase interface {
		NewWallet(ctx *gin.Context, request *dtos.WalletAddRequest) (response interface{}, httpCode int, errInfo []string)
		addWalletProAccount(request *dtos.WalletAddRequest, walletType string, accountUUID, UUIDIDMasterWalletType uuid.UUID) (walletID uuid.UUID, err error)
		addWalletBasicAccount(request *dtos.WalletAddRequest, walletType string, accountUUID, UUIDIDMasterWalletType uuid.UUID) (walletID uuid.UUID, err error)
		addNewWallet(walletEntity *entities.WalletEntity) (err error)
		setWalletType(IDMasterWallet string) (walletType string)
		setInitialBalance(walletType string, amount int64, walletID, accountUUID uuid.UUID) (err error)
		setBalanceInvestment(amount int64, walletID, accountUUID uuid.UUID) (err error)
		setBalanceNonInvestment(amount int64, walletID, accountUUID uuid.UUID) (err error)
	}
)

func NewWalletUseCase(repo IWalletRepository) *WalletUseCase {
	return &WalletUseCase{
		repo: repo,
	}
}

func (s *WalletUseCase) NewWallet(ctx *gin.Context, request *dtos.WalletAddRequest) (response interface{}, httpCode int, errInfo []string) {
	var (
		idWallet   uuid.UUID
		err        error
		walletType string
	)

	// get information from context : account type and account uuid
	accountType, accountUUID := personalaccounts.AccountInformation(ctx)

	// determine wallet type from id master wallet
	walletType = s.setWalletType(request.IDMasterWallet)

	// change type from string into uuid type
	UUIDIDMasterWalletType, err := uuid.Parse(request.IDMasterWallet)
	if err != nil {
		errInfo = errorsinfo.ErrorWrapperArray(errInfo, err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	// save wallet based on account type profile
	switch accountType {
	case constants.AccountPro:
		idWallet, err = s.addWalletProAccount(request, walletType, accountUUID, UUIDIDMasterWalletType)
		break
	case constants.AccountBasic:
		idWallet, err = s.addWalletBasicAccount(request, walletType, accountUUID, UUIDIDMasterWalletType)
		break
	}

	if err != nil {
		errInfo = errorsinfo.ErrorWrapperArray(errInfo, err.Error())
		return struct{}{}, http.StatusUnprocessableEntity, errInfo
	}

	// mapping success response
	resp := struct {
		ID      uuid.UUID `json:"wallet_id"`
		Message string    `json:"message"`
	}{
		ID:      idWallet,
		Message: "success add new wallet",
	}

	return resp, http.StatusOK, []string{}
}

func (s *WalletUseCase) addWalletProAccount(request *dtos.WalletAddRequest, walletType string, accountUUID, UUIDIDMasterWalletType uuid.UUID) (walletID uuid.UUID, err error) {
	// determine new id wallet
	walletID = uuid.New()

	// wallet model
	walletEntity := entities.WalletEntity{
		ID:                 walletID,
		Active:             true,
		WalletName:         request.WalletName,
		WalletType:         walletType,
		IDMasterWalletType: UUIDIDMasterWalletType,
		FeeInvestBuy:       request.FeeInvestBuy,
		FeeInvestSell:      request.FeeInvestSell,
		TotalAssets:        request.TotalAsset,
		IDAccount:          accountUUID,
	}

	// add new wallet
	err = s.addNewWallet(&walletEntity)
	if err != nil {
		return uuid.Nil, err
	}

	return walletID, nil
}

func (s *WalletUseCase) addWalletBasicAccount(request *dtos.WalletAddRequest, walletType string, accountUUID, UUIDIDMasterWalletType uuid.UUID) (walletID uuid.UUID, err error) {
	// number of current wallet by account ID
	totalWallet, err := s.repo.NumberOfWalletsByID(accountUUID)
	if err != nil {
		return uuid.Nil, err
	}

	// maximum 2 wallets each account
	if totalWallet == constants.MaxWalletBasic {
		return uuid.Nil, errors.New(constants.ProPlan)
	}

	// determine new id wallet
	walletID = uuid.New()

	// wallet model
	walletEntity := entities.WalletEntity{
		ID:                 walletID,
		Active:             true,
		WalletName:         request.WalletName,
		WalletType:         walletType,
		IDMasterWalletType: UUIDIDMasterWalletType,
		FeeInvestBuy:       request.FeeInvestBuy,
		FeeInvestSell:      request.FeeInvestSell,
		TotalAssets:        request.TotalAsset,
		IDAccount:          accountUUID,
	}

	// add new wallet
	err = s.addNewWallet(&walletEntity)
	if err != nil {
		return uuid.Nil, err
	}

	// set initial balance
	err = s.setInitialBalance(walletType, request.TotalAsset, walletID, accountUUID)
	if err != nil {
		return uuid.Nil, err
	}

	return walletID, nil
}

func (s *WalletUseCase) addNewWallet(walletEntity *entities.WalletEntity) (err error) {
	return s.repo.NewWallet(walletEntity)
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

func (s *WalletUseCase) setInitialBalance(walletType string, amount int64, walletID, accountUUID uuid.UUID) (err error) {
	switch walletType {
	case constants.Investment:
		return s.setBalanceInvestment(amount, walletID, accountUUID)
	case constants.Cash, constants.CreditCard, constants.DebitCard, constants.EWallet, constants.Saving:
		return s.setBalanceNonInvestment(amount, walletID, accountUUID)
	}
	return nil
}

func (s *WalletUseCase) setBalanceInvestment(amount int64, walletID, accountUUID uuid.UUID) (err error) {
	// setup balance
	newTransactionAsBalance := entities.WalletInitTransactionInvestment{
		IDWallet:          walletID,
		Balance:           float64(amount),
		IDPersonalAccount: accountUUID,
	}

	// record as initialized balance
	return s.repo.SetBalanceInvestment(&newTransactionAsBalance)
}

func (s *WalletUseCase) setBalanceNonInvestment(amount int64, walletID, accountUUID uuid.UUID) (err error) {
	// setup balance
	newTransactionAsBalance := entities.WalletInitTransaction{
		ID:                uuid.New(),
		Date:              datecustoms.NowTransaction(),
		Fees:              0,
		Amount:            float64(amount),
		IDPersonalAccount: accountUUID,
		IDWallet:          walletID,
		Credit:            float64(amount),
		Debit:             0,
		Balance:           float64(amount),
	}

	// no setup trx detail
	trxDetail := entities.WalletInitTransactionDetail{}

	// record as initialized balance
	return s.repo.SetBalanceNonInvestment(&newTransactionAsBalance, &trxDetail)
}