package wallets

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
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
		GetAllWallets(ctx *gin.Context) (response interface{}, httpCode int, errInfo []string)
		UpdateAmount(ctx *gin.Context, walletID string, request map[string]interface{}) (response interface{}, httpCode int, errInfo []string)
		addWalletProAccount(request *dtos.WalletAddRequest, walletType string, accountUUID, UUIDIDMasterWalletType uuid.UUID) (walletID uuid.UUID, err error)
		addWalletBasicAccount(request *dtos.WalletAddRequest, walletType string, accountUUID, UUIDIDMasterWalletType uuid.UUID) (walletID uuid.UUID, err error)
		addNewWallet(walletEntity *entities.WalletEntity) (err error)
		setWalletType(IDMasterWallet string) (walletType string)
		setInitialBalance(walletType string, amount int64, walletID, accountUUID uuid.UUID) (err error)
		setBalanceInvestment(amount int64, walletID, accountUUID uuid.UUID) (err error)
		setBalanceNonInvestment(amount int64, walletID, accountUUID uuid.UUID) (err error)
		latestBalance(walletType string, walletID uuid.UUID) (amount int64)
		getBalanceInvestment(walletID uuid.UUID) (amount int64)
		getBalanceNonInvestment(walletID uuid.UUID) (amount int64)
		getValueFromRequest(request map[string]interface{}) (walletName string, amount float64, err error)
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

func (s *WalletUseCase) GetAllWallets(ctx *gin.Context) (response interface{}, httpCode int, errInfo []string) {
	var dtoResponse []dtos.WalletListResponse

	// get information from context : account type and account uuid
	_, accountUUID := personalaccounts.AccountInformation(ctx)

	// get all wallet by account ID
	walletCollection, err := s.repo.GetAllWallets(accountUUID)
	if err != nil {
		errInfo = errorsinfo.ErrorWrapperArray(errInfo, err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	// if no collection of wallet data
	if len(walletCollection) == 0 {
		errInfo = errorsinfo.ErrorWrapperArray(errInfo, errors.New("this account has no wallet").Error())
		return struct{}{}, http.StatusNotFound, errInfo
	}

	for _, v := range walletCollection {
		// get latest balance
		balance := s.latestBalance(v.WalletType, v.ID)

		// mapping response
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
			TotalAssets:   balance,
		})
	}

	return dtoResponse, http.StatusOK, []string{}

}

func (s *WalletUseCase) UpdateAmount(ctx *gin.Context, walletID string, request map[string]interface{}) (response interface{}, httpCode int, errInfo []string) {
	var (
		err error
	)

	// get information from context : account type and account uuid
	_, accountUUID := personalaccounts.AccountInformation(ctx)

	// convert wallet id type into uuid.UUID
	walletUUID, err := uuid.Parse(walletID)
	if err != nil {
		logrus.Error(err.Error())
	}

	// get wallet type
	dataWallets, err := s.repo.GetWalletType(walletUUID)
	if err != nil {
		errInfo = errorsinfo.ErrorWrapperArray(errInfo, err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	// get value of request
	_, amount, err := s.getValueFromRequest(request)
	if err != nil {
		errInfo = errorsinfo.ErrorWrapperArray(errInfo, err.Error())
		return struct{}{}, http.StatusUnprocessableEntity, errInfo
	}

	// set for latest balance
	if amount > 0 {
		err = s.setInitialBalance(dataWallets.WalletType, int64(amount), walletUUID, accountUUID)
	}

	// update wallet information
	err = s.repo.UpdateWalletInformation(walletUUID, request)
	if err != nil {
		errInfo = errorsinfo.ErrorWrapperArray(errInfo, err.Error())
		return struct{}{}, http.StatusUnprocessableEntity, errInfo
	}

	resp := struct {
		Message string `json:"message"`
	}{
		Message: "update wallet success",
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

	// set initial balance
	err = s.setInitialBalance(walletType, request.TotalAsset, walletID, accountUUID)
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
	return walletType
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

func (s *WalletUseCase) latestBalance(walletType string, walletID uuid.UUID) (amount int64) {
	switch walletType {
	case constants.Investment:
		return s.getBalanceInvestment(walletID)
	case constants.Cash, constants.CreditCard, constants.DebitCard, constants.EWallet, constants.Saving:
		return s.getBalanceNonInvestment(walletID)
	}
	return 0
}

func (s *WalletUseCase) getBalanceInvestment(walletID uuid.UUID) (amount int64) {
	data, err := s.repo.GetBalanceInvestment(walletID)
	if err != nil {
		return 0
	}
	return int64(data.Balance)
}

func (s *WalletUseCase) getBalanceNonInvestment(walletID uuid.UUID) (amount int64) {
	data, err := s.repo.GetBalanceNonInvestment(walletID)
	if err != nil {
		return 0
	}
	return int64(data.Balance)
}

func (s *WalletUseCase) getValueFromRequest(request map[string]interface{}) (walletName string, amount float64, err error) {
	// check wallet name exist from payload
	value, exists := request["wallet_name"]
	if exists {
		walletName = fmt.Sprintf("%v", value)
		if walletName == "" {
			return "", 0, errors.New("wallet name empty value")
		}
	}

	// check amount exist from payload
	value, exists = request["amount"]
	if exists {
		amount = value.(float64)
		if amount <= 0 {
			return "", 0, errors.New("amount must greater than 0")
		}
	}

	return walletName, amount, nil
}
