package transactions

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/wealthy-app/wealthy-backend/api/v1/transactions/dtos"
	"github.com/wealthy-app/wealthy-backend/api/v1/transactions/entities"
	"github.com/wealthy-app/wealthy-backend/utils/errorsinfo"
	"github.com/wealthy-app/wealthy-backend/utils/personalaccounts"
	"github.com/wealthy-app/wealthy-backend/utils/utilities"
	"net/http"
	"sort"
)

type (
	TransactionUseCase struct {
		repo ITransactionRepository
	}

	ITransactionUseCase interface {
		Add(ctx *gin.Context, request *dtos.TransactionRequest) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		AddInvestmentTransaction(ctx *gin.Context, request *dtos.TransactionRequestInvestment) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		ExpenseTransactionHistory(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		IncomeTransactionHistory(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		TravelTransactionHistory(ctx *gin.Context, IDTravel uuid.UUID) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		TransferTransactionHistory(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		InvestTransactionHistory(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		IncomeSpending(ctx *gin.Context, month string, year string) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		Investment(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		ByNotes(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		Suggestion(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		CashFlow(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		WalletNonInvestment(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		WalletInvestment(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		saveInvestTransaction(accountID uuid.UUID, request *dtos.TransactionRequestInvestment) (id uuid.UUID, err error)
		investmentCalculation(accountID uuid.UUID) (err error)
	}
)

func NewTransactionUseCase(repo ITransactionRepository) *TransactionUseCase {
	return &TransactionUseCase{repo: repo}
}

func (s *TransactionUseCase) Add(ctx *gin.Context, request *dtos.TransactionRequest) (data interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		trxID                 uuid.UUID
		err                   error
		convertAmount         int64
		IDTravelUUID          uuid.UUID
		IDMasterIncCatUUID    uuid.UUID
		IDMasterExpCatUUID    uuid.UUID
		IDMasterSubExpCatUUID uuid.UUID
		IDMasterTransPriUUID  uuid.UUID
		IDMasterTransTypeUUID uuid.UUID
		IDWalletUUID          uuid.UUID
		Credit                float64
		Debit                 float64
		Balance               float64
	)

	accountUUID := ctx.MustGet("accountID").(uuid.UUID)

	// if transaction is travel
	if request.IDTravel != "" {
		IDTravelUUID, err = uuid.Parse(request.IDTravel)
		if err != nil {
			logrus.Error(err.Error())
		}

		dataCurrency, err := s.repo.BudgetWithCurrency(IDTravelUUID)
		if err != nil {
			logrus.Error(err.Error())
		}
		convertAmount = dataCurrency.CurrencyValue * request.Amount
	}

	// if transaction is not travel
	if request.IDTravel == "" {
		convertAmount = request.Amount
		IDTravelUUID = uuid.Nil
	}

	// translate string to uuid
	if request.IDMasterIncomeCategories != "" {
		IDMasterIncCatUUID, err = uuid.Parse(request.IDMasterIncomeCategories)
		if err != nil {
			logrus.Error(err.Error())
		}

		// convert string to UUID
		IDWalletUUID, err = uuid.Parse(request.IDWallet)
		if err != nil {
			logrus.Error(err.Error())
		}

		// is wallet true exists
		if !s.repo.WalletExist(IDWalletUUID) {
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id wallet unregistered before")
			return struct{}{}, http.StatusBadRequest, errInfo
		}

	}

	if request.IDMasterExpenseCategories != "" {
		IDMasterExpCatUUID, err = uuid.Parse(request.IDMasterExpenseCategories)
		if err != nil {
			logrus.Error(err.Error())
		}

		// convert string to UUID
		IDWalletUUID, err = uuid.Parse(request.IDWallet)
		if err != nil {
			logrus.Error(err.Error())
		}

		if request.IDTravel == "" {
			// is wallet true exists
			if !s.repo.WalletExist(IDWalletUUID) {
				errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id wallet unregistered before")
				return struct{}{}, http.StatusBadRequest, errInfo
			}

			// check sub expense
			if request.IDMasterExpenseSubCategories != "" {
				if request.IDMasterExpenseCategories == "" {
					errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id master expense categories need along use id master expense subcategories")
					return struct{}{}, http.StatusBadRequest, errInfo
				}

				IDMasterSubExpCatUUID, err = uuid.Parse(request.IDMasterExpenseSubCategories)
				if err != nil {
					logrus.Error(err.Error())
				}
			}

			// deduction last balance
		}
	}

	if request.IDMasterIncomeCategories != "" || request.IDMasterExpenseCategories != "" {
		if request.IDMasterTransactionPriorities != "" {
			IDMasterTransPriUUID, err = uuid.Parse(request.IDMasterTransactionPriorities)
			if err != nil {
				logrus.Error(err.Error())
			}
		}
	}

	// all transactions using this
	if request.IDMasterTransactionTypes != "" {
		IDMasterTransTypeUUID, err = uuid.Parse(request.IDMasterTransactionTypes)
		if err != nil {
			logrus.Error(err.Error())
		}
	}

	// get last balance
	dataLastBalance, err := s.repo.LastBalance(accountUUID, IDWalletUUID)
	if err != nil {
		logrus.Error(err.Error())
	}

	if request.IDMasterIncomeCategories != "" {
		Balance = dataLastBalance.Balance + float64(request.Amount)
		Debit = 0
		Credit = float64(request.Amount)
	}

	if request.IDMasterExpenseCategories != "" {
		Balance = dataLastBalance.Balance - float64(request.Amount)
		Credit = 0
		Debit = float64(request.Amount)
	}

	trxID = uuid.New()
	modelTransaction := entities.TransactionEntity{
		ID:                            trxID,
		Date:                          request.Date,
		Fees:                          float64(request.Fees),
		Amount:                        float64(convertAmount),
		IDPersonalAccount:             accountUUID,
		IDWallet:                      IDWalletUUID,
		IDMasterIncomeCategories:      IDMasterIncCatUUID,
		IDMasterExpenseCategories:     IDMasterExpCatUUID,
		IDMasterExpenseSubCategories:  IDMasterSubExpCatUUID,
		IDMasterTransactionPriorities: IDMasterTransPriUUID,
		IDMasterTransactionTypes:      IDMasterTransTypeUUID,
		Balance:                       Balance,
		Debit:                         Debit,
		Credit:                        Credit,
	}

	modelTransactionDetail := entities.TransactionDetailEntity{
		IDTransaction: trxID,
		Repeat:        request.Repeat,
		Note:          request.Note,
		From:          request.TransferFrom,
		To:            request.TransferTo,
		IDTravel:      IDTravelUUID,
	}

	// save transaction
	err = s.repo.Add(&modelTransaction, &modelTransactionDetail)
	if err != nil {
		data := struct {
			IDTransaction uuid.UUID
		}{
			IDTransaction: uuid.Nil,
		}
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return data, http.StatusInternalServerError, errInfo
	}

	data = struct {
		IDTransaction uuid.UUID `json:"id_transaction"`
	}{
		IDTransaction: trxID,
	}

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	return data, http.StatusOK, errInfo
}

func (s *TransactionUseCase) AddInvestmentTransaction(ctx *gin.Context, request *dtos.TransactionRequestInvestment) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		trxID             uuid.UUID
		err               error
		collectionsWallet = make(map[string]bool)
	)

	// account uuid
	accountUUID := ctx.MustGet("accountID").(uuid.UUID)

	// validate wallet type
	dataWalletType, err := s.repo.WalletInvestment(accountUUID)
	if err != nil {
		logrus.Error()
	}

	if len(dataWalletType) == 0 {
		resp := struct {
			Message string `json:"message"`
		}{
			Message: "no wallet type for investment. please create new one",
		}
		return resp, http.StatusNotFound, []errorsinfo.Errors{}
	}

	// if there is data , store to map
	if len(dataWalletType) > 0 {
		for _, v := range dataWalletType {
			collectionsWallet[v.ID.String()] = true
		}
	}

	_, exists := collectionsWallet[request.IDWallet]
	if !exists {
		resp := struct {
			Message string `json:"message"`
		}{
			Message: "no wallet type for investment. please create new one",
		}
		return resp, http.StatusNotFound, []errorsinfo.Errors{}
	}

	// save investment transaction
	trxID, err = s.saveInvestTransaction(accountUUID, request)
	if err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	// do calculation. fetching data from transaction detail
	err = s.investmentCalculation(accountUUID)
	if err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	// prepare for response
	resp := struct {
		IDTransaction uuid.UUID `json:"id_transaction"`
	}{
		IDTransaction: trxID,
	}

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	return resp, http.StatusOK, errInfo
}

func (s *TransactionUseCase) ExpenseTransactionHistory(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		dtoResponse                  dtos.TransactionHistoryForIncomeExpenses
		responseExpenseTotalHistory  entities.TransactionExpenseTotalHistory
		responseExpenseDetailHistory interface{}
	)

	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")

	accountUUID := ctx.MustGet("accountID").(uuid.UUID)

	if startDate == "" || endDate == "" {
		responseExpenseTotalHistory = s.repo.ExpenseTotalHistoryWithoutDate(accountUUID)
		responseExpenseDetailHistory = s.repo.ExpenseDetailHistoryWithoutDate(accountUUID)
	} else {
		responseExpenseTotalHistory = s.repo.ExpenseTotalHistoryWithDate(accountUUID, startDate, endDate)
		responseExpenseDetailHistory = s.repo.ExpenseDetailHistoryWithDate(accountUUID, startDate, endDate)
	}

	if responseExpenseTotalHistory.TotalExpense == 0 || responseExpenseDetailHistory == nil {
		httpCode = http.StatusNotFound
		response := struct {
			Message string `json:"message"`
		}{
			Message: "there is not expense transaction between periods : " + startDate + " until " + endDate,
		}
		return response, httpCode, errInfo
	}

	dtoResponse.Total = responseExpenseTotalHistory.TotalExpense
	dtoResponse.Detail = responseExpenseDetailHistory

	return dtoResponse, http.StatusOK, []errorsinfo.Errors{}
}

func (s *TransactionUseCase) IncomeTransactionHistory(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		dtoResponse                 dtos.TransactionHistoryForIncomeExpenses
		responseIncomeTotalHistory  entities.TransactionIncomeTotalHistory
		responseIncomeDetailHistory interface{}
	)

	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")

	usrEmail := ctx.MustGet("email").(string)
	personalAccount := personalaccounts.Informations(ctx, usrEmail)

	if personalAccount.ID == uuid.Nil {
		httpCode = http.StatusNotFound
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "token contains invalid information")
		return response, httpCode, errInfo
	}

	if startDate == "" || endDate == "" {
		responseIncomeTotalHistory = s.repo.IncomeTotalHistoryWithoutDate(personalAccount.ID)
		responseIncomeDetailHistory = s.repo.IncomeDetailHistoryWithoutData(personalAccount.ID)
	} else {
		responseIncomeTotalHistory = s.repo.IncomeTotalHistoryWithData(personalAccount.ID, startDate, endDate)
		responseIncomeDetailHistory = s.repo.IncomeDetailHistoryWithData(personalAccount.ID, startDate, endDate)
	}

	if responseIncomeTotalHistory.TotalIncome == 0 || responseIncomeDetailHistory == nil {
		httpCode = http.StatusNotFound
		response := struct {
			Message string `json:"message"`
		}{
			Message: "there is not income transaction between periods : " + startDate + " until " + endDate,
		}
		return response, httpCode, errInfo
	}

	dtoResponse.Total = responseIncomeTotalHistory.TotalIncome
	dtoResponse.Detail = responseIncomeDetailHistory

	return dtoResponse, http.StatusOK, []errorsinfo.Errors{}
}

func (s *TransactionUseCase) TravelTransactionHistory(ctx *gin.Context, IDTravel uuid.UUID) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		dtoResponse                 dtos.TransactionHistoryForTravel
		details                     []dtos.TransactionHistoryForTravelDetail
		responseTravelDetailHistory []entities.TransactionDetailTravel
	)

	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")

	accountUUID := ctx.MustGet("accountID").(uuid.UUID)

	if startDate == "" || endDate == "" {
		responseTravelDetailHistory = s.repo.TravelDetailWithoutData(accountUUID, IDTravel)
	} else {
		responseTravelDetailHistory = s.repo.TravelDetailWithData(accountUUID, IDTravel, startDate, endDate)
	}

	if len(responseTravelDetailHistory) == 0 {
		resp := struct {
			Message string `json:"message"`
		}{
			Message: "no data for travel transaction",
		}
		return resp, http.StatusNotFound, []errorsinfo.Errors{}
	}

	if len(responseTravelDetailHistory) > 0 {
		for _, v := range responseTravelDetailHistory {
			details = append(details, dtos.TransactionHistoryForTravelDetail{
				DateTransaction: v.DateTransaction,
				IDTransaction:   v.IDTransaction,
				Amount: dtos.Amount{
					CurrencyCode: "IDR",
					Value:        float64(v.Amount),
				},
				Category:     v.Category,
				CategoryIcon: v.TransactionCatIcon,
				Note:         v.Note,
			})
		}
	}
	dtoResponse.Detail = details
	return dtoResponse, http.StatusOK, []errorsinfo.Errors{}
}

func (s *TransactionUseCase) TransferTransactionHistory(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		dtoResponse                   dtos.TransactionHistoryForTransfer
		responseTransferDetailHistory []entities.TransactionDetailTransfer
	)

	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")

	usrEmail := ctx.MustGet("email").(string)
	personalAccount := personalaccounts.Informations(ctx, usrEmail)

	if personalAccount.ID == uuid.Nil {
		httpCode = http.StatusNotFound
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "token contains invalid information")
		return response, httpCode, errInfo
	}

	if startDate == "" || endDate == "" {
		responseTransferDetailHistory = s.repo.TransferDetailWithoutData(personalAccount.ID)
	} else {
		responseTransferDetailHistory = s.repo.TransferDetailWithData(personalAccount.ID, startDate, endDate)
	}

	if len(responseTransferDetailHistory) == 0 {
		response := struct {
			Message string `json:"message"`
		}{
			Message: "there is not transfer transaction between periods : " + startDate + " until " + endDate,
		}
		return response, http.StatusNotFound, []errorsinfo.Errors{}
	}

	dtoResponse.Detail = responseTransferDetailHistory
	return dtoResponse, http.StatusOK, []errorsinfo.Errors{}
}

func (s *TransactionUseCase) InvestTransactionHistory(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		dtoResponse                 dtos.TransactionHistoryForInvest
		responseInvestDetailHistory []entities.TransactionDetailInvest
	)

	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")

	usrEmail := ctx.MustGet("email").(string)
	personalAccount := personalaccounts.Informations(ctx, usrEmail)

	if personalAccount.ID == uuid.Nil {
		httpCode = http.StatusNotFound
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "not found")
		return response, httpCode, errInfo
	}

	if startDate == "" || endDate == "" {
		responseInvestDetailHistory = s.repo.InvestDetailWithoutData(personalAccount.ID)
	} else {
		responseInvestDetailHistory = s.repo.InvestDetailWithData(personalAccount.ID, startDate, endDate)
	}

	if len(responseInvestDetailHistory) == 0 {
		response := struct {
			Message string `json:"message"`
		}{
			Message: "there is not invest transaction between periods : " + startDate + " until " + endDate,
		}
		return response, http.StatusNotFound, []errorsinfo.Errors{}
	}

	dtoResponse.Detail = responseInvestDetailHistory
	return dtoResponse, http.StatusOK, []errorsinfo.Errors{}
}

func (s *TransactionUseCase) IncomeSpending(ctx *gin.Context, month string, year string) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		dtoResponse                          dtos.TransactionIncomeSpendingInvestment
		dtoResponseAnnually                  dtos.TransactionIncomeSpendingInvestmentAnnually
		dtoAnnualyDetail                     []dtos.TransactionDetailAnnually
		responseIncomeSpendingTotal          interface{}
		responseIncomeSpendingDetailMonthly  []entities.TransactionIncomeSpendingDetailMonthly
		responseIncomeSpendingDetailAnnually []entities.TransactionIncomeSpendingDetailAnnually
		detailsMonthly                       []dtos.TransactionIncomeSpendingInvestmentDetail
		deepDetailsMonthly                   []dtos.TransactionDetails
	)

	accountUUID := ctx.MustGet("accountID").(uuid.UUID)

	// if month and year not empty
	if month != "" && year != "" {
		responseIncomeSpendingTotal = s.repo.IncomeSpendingMonthlyTotal(accountUUID, month, year)
		responseIncomeSpendingDetailMonthly = s.repo.IncomeSpendingMonthlyDetail(accountUUID, month, year)

		isNotExist := len(responseIncomeSpendingDetailMonthly) == 0
		if isNotExist {
			resp := struct {
				Message string `json:"message"`
			}{
				Message: "no data for income-spending",
			}
			return resp, http.StatusNotFound, []errorsinfo.Errors{}
		}
	}

	// if there is year only
	if month == "" && year != "" {
		responseIncomeSpendingTotal = s.repo.IncomeSpendingAnnuallyTotal(accountUUID, year)
		responseIncomeSpendingDetailAnnually = s.repo.IncomeSpendingAnnuallyDetail(accountUUID, year)

		//if len(responseIncomeSpendingDetailAnnually) > 0 {
		//	isNotExist := responseIncomeSpendingDetailAnnually[0].NetIncome == 0 && responseIncomeSpendingDetailAnnually[0].TotalIncome == 0 && responseIncomeSpendingDetailAnnually[0].TotalSpending == 0
		//	if isNotExist {
		//		resp := struct {
		//			Message string `json:"message"`
		//		}{
		//			Message: "no data for income-spending",
		//		}
		//		return resp, http.StatusNotFound, []errorsinfo.Errors{}
		//	}
		//}

		if len(responseIncomeSpendingDetailAnnually) == 0 {
			resp := struct {
				Message string `json:"message"`
			}{
				Message: "no data for income-spending",
			}
			return resp, http.StatusNotFound, []errorsinfo.Errors{}
		}
	}

	// todo : figure out the another way
	//if len(responseIncomeSpendingDetailMonthly) == 0 {
	//	return dtos.TransactionIncomeSpendingInvestment{}, http.StatusNotFound, errInfo
	//}
	//
	//if len(responseIncomeSpendingDetailAnnually) == 0 {
	//	return dtos.TransactionIncomeSpendingInvestmentAnnually{}, http.StatusNotFound, errInfo
	//}

	if len(responseIncomeSpendingDetailMonthly) > 0 {

		var dateTempPrev string
		length := len(responseIncomeSpendingDetailMonthly)

		for k, v := range responseIncomeSpendingDetailMonthly {

			// for first time
			if dateTempPrev == "" {
				dateTempPrev = v.TransactionDate

				deepDetailsMonthly = append(deepDetailsMonthly, dtos.TransactionDetails{
					TransactionCategory: v.TransactionCategory,
					TransactionType:     v.TransactionType,
					TransactionAmount: dtos.Amount{
						CurrencyCode: "IDR",
						Value:        float64(v.TransactionAmount),
					},
					TransactionNote: v.TransactionNote,
				})

				if k == length-1 {
					detailsMonthly = append(detailsMonthly, dtos.TransactionIncomeSpendingInvestmentDetail{
						TransactionDate:    dateTempPrev,
						TransactionDetails: deepDetailsMonthly,
					})

					dtoResponse.Detail = append(dtoResponse.Detail, detailsMonthly...)
				}
			}

			// if previous is different current
			if dateTempPrev != v.TransactionDate {
				deepDetailsMonthly = append(deepDetailsMonthly, dtos.TransactionDetails{
					TransactionCategory: v.TransactionCategory,
					TransactionType:     v.TransactionType,
					TransactionAmount: dtos.Amount{
						CurrencyCode: "IDR",
						Value:        float64(v.TransactionAmount),
					},
					TransactionNote: v.TransactionNote,
				})

				detailsMonthly = append(detailsMonthly, dtos.TransactionIncomeSpendingInvestmentDetail{
					TransactionDate:    dateTempPrev,
					TransactionDetails: deepDetailsMonthly,
				})

				dtoResponse.Detail = append(dtoResponse.Detail, detailsMonthly...)

				// clear data
				detailsMonthly = []dtos.TransactionIncomeSpendingInvestmentDetail{}
				deepDetailsMonthly = []dtos.TransactionDetails{}

				dateTempPrev = v.TransactionDate

				if k == (length - 1) {
					deepDetailsMonthly = append(deepDetailsMonthly, dtos.TransactionDetails{
						TransactionCategory: v.TransactionCategory,
						TransactionType:     v.TransactionType,
						TransactionAmount: dtos.Amount{
							CurrencyCode: "IDR",
							Value:        float64(v.TransactionAmount),
						},
						TransactionNote: v.TransactionNote,
					})

					detailsMonthly = append(detailsMonthly, dtos.TransactionIncomeSpendingInvestmentDetail{
						TransactionDate:    dateTempPrev,
						TransactionDetails: deepDetailsMonthly,
					})
				}
			}
		}
		dtoResponse.Summary = responseIncomeSpendingTotal
		return dtoResponse, http.StatusOK, []errorsinfo.Errors{}
	}

	if len(responseIncomeSpendingDetailAnnually) > 0 {

		for _, v := range responseIncomeSpendingDetailAnnually {
			dtoAnnualyDetail = append(dtoAnnualyDetail, dtos.TransactionDetailAnnually{
				LastDayInMonth:  utilities.GetLastDay(v.DateOrigin),
				MonthYear:       v.MonthYear,
				TotalDayInMonth: v.TotalDayInMonth,
				TotalIncome:     v.TotalIncome,
				TotalSpending:   v.TotalSpending,
				NetIncome:       v.NetIncome,
			})
		}

		dtoResponseAnnually.Summary = responseIncomeSpendingTotal
		dtoResponseAnnually.Detail = dtoAnnualyDetail
		return dtoResponseAnnually, http.StatusOK, []errorsinfo.Errors{}
	}

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	return response, http.StatusPreconditionFailed, errInfo
}

func (s *TransactionUseCase) Investment(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		dtoResponse              dtos.TransactionInvestment
		detail                   []dtos.TransactionInvestmentDetail
		trxDetails               []dtos.TransactionInvestDetails
		responseInvestmentTotal  interface{}
		responseInvestmentDetail []entities.TransactionInvestmentDetail
	)

	month := ctx.Query("month")
	year := ctx.Query("year")

	accountUUID := ctx.MustGet("accountID").(uuid.UUID)

	if month == "" && year == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "need month and or year information")
		return struct{}{}, http.StatusBadRequest, errInfo
	}

	if month != "" && year != "" {
		responseInvestmentTotal = s.repo.InvestMonthlyTotal(accountUUID, month, year)
		responseInvestmentDetail = s.repo.InvestMonthlyDetail(accountUUID, month, year)

	}

	if month == "" && year != "" {
		responseInvestmentTotal = s.repo.InvestAnnuallyTotal(accountUUID, year)
		responseInvestmentDetail = s.repo.InvestAnnuallyDetail(accountUUID, year)
	}

	switch v := responseInvestmentTotal.(type) {
	case entities.TransactionInvestmentTotals:
		if v.TotalBuy == 0 && v.TotalCurrentPortfolio == 0 && v.TotalSell == 0 {
			resp := struct {
				Message string `json:"message"`
			}{
				Message: "no data for investment transaction",
			}

			return resp, http.StatusNotFound, []errorsinfo.Errors{}
		}
	}

	dtoResponse.Summary = responseInvestmentTotal

	dateTemp := ""
	maxData := len(responseInvestmentDetail) - 1
	sellBuy := ""

	if len(responseInvestmentDetail) > 0 {
		for k, v := range responseInvestmentDetail {

			if v.SellBuy == 0 {
				sellBuy = "SELL"
			} else {
				sellBuy = "BUY"
			}

			// get trading
			dataTrading, err := s.repo.GetTradingInfo(v.StockCode)
			if err != nil {
				logrus.Error(err.Error())
			}

			// mapping process

			if dateTemp == v.Date {
				// populate
				trxDetails = append(trxDetails, dtos.TransactionInvestDetails{
					LotWithPrice: float64(int64(v.Lot) * v.Price),
					Name:         dataTrading.Name,
					Lot:          v.Lot,
					StockCode:    v.StockCode,
					Price:        v.Price,
					SellBuy:      sellBuy,
				})

				// reach end
				if k == maxData {
					// append to details
					detail = append(detail, dtos.TransactionInvestmentDetail{
						TransactionDate:    dateTemp,
						TransactionDetails: trxDetails,
					})
				}
			}

			// if first time
			if dateTemp == "" {
				dateTemp = v.Date

				// populate
				trxDetails = append(trxDetails, dtos.TransactionInvestDetails{
					LotWithPrice: float64(int64(v.Lot) * v.Price),
					Name:         dataTrading.Name,
					Lot:          v.Lot,
					StockCode:    v.StockCode,
					Price:        v.Price,
					SellBuy:      sellBuy,
				})

				// reach end
				if k == maxData {
					// append to details
					detail = append(detail, dtos.TransactionInvestmentDetail{
						TransactionDate:    dateTemp,
						TransactionDetails: trxDetails,
					})

					// clear previous
					trxDetails = nil
				}

			}

			if dateTemp != v.Date {
				// append to details previous
				detail = append(detail, dtos.TransactionInvestmentDetail{
					TransactionDate:    dateTemp,
					TransactionDetails: trxDetails,
				})

				// clear previous
				trxDetails = nil

				// renew
				dateTemp = v.Date

				// populate
				trxDetails = append(trxDetails, dtos.TransactionInvestDetails{
					LotWithPrice: float64(int64(v.Lot) * v.Price),
					Name:         dataTrading.Name,
					Lot:          v.Lot,
					StockCode:    v.StockCode,
					Price:        v.Price,
					SellBuy:      sellBuy,
				})

				// reach end
				if k == maxData {
					// append to details
					detail = append(detail, dtos.TransactionInvestmentDetail{
						TransactionDate:    dateTemp,
						TransactionDetails: trxDetails,
					})

					// clear previous
					trxDetails = nil
				}
			}
		}
	}

	dtoResponse.Detail = detail

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	return dtoResponse, http.StatusOK, errInfo
}

func (s *TransactionUseCase) ByNotes(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		dtoResponse      dtos.TransactionNotes
		detailNotes      dtos.TransactionNotesDetail
		deepDetailsNotes []dtos.TransactionNotesDeepDetail
		dataNotes        []entities.TransactionByNotes
	)

	month := ctx.Query("month")
	year := ctx.Query("year")

	usrEmail := ctx.MustGet("email").(string)
	personalAccount := personalaccounts.Informations(ctx, usrEmail)

	if personalAccount.ID == uuid.Nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "token contains invalid information")
		return struct{}{}, http.StatusBadRequest, errInfo
	}

	if month == "" && year == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "both month, year need in query url")
		return struct{}{}, http.StatusBadRequest, errInfo
	}

	if month != "" && year != "" {
		dataNotes = s.repo.ByNote(personalAccount.ID, month, year)
	}

	lengthData := len(dataNotes)

	if len(dataNotes) == 0 {
		resp := struct {
			Message string `json:"message,omitempty"`
		}{
			Message: "no data transaction by notes",
		}
		return resp, http.StatusNotFound, []errorsinfo.Errors{}
	}

	if len(dataNotes) > 0 {
		var catPrev string

		for k, v := range dataNotes {

			if catPrev == v.TransactionCategory {
				if k == (lengthData - 1) {
					//detailNotes.TransactionNotesDeepDetail = append(detailNotes.TransactionNotesDeepDetail, deepDetailsNotes...)
					//detailNotes.TransactionCategory = catPrev

					deepDetailsNotes = append(deepDetailsNotes, dtos.TransactionNotesDeepDetail{
						TransactionNote: v.TransactionNote,
						TransactionAmount: dtos.Amount{
							CurrencyCode: "IDR",
							Value:        v.Amount,
						},
						TransactionLimit: dtos.Amount{
							CurrencyCode: "IDR",
							Value:        v.Budget,
						},
					})

					detailNotes.TransactionCategory = v.TransactionCategory
					detailNotes.TransactionNotesDeepDetail = append(detailNotes.TransactionNotesDeepDetail, deepDetailsNotes...)

					dtoResponse.TransactionNotesDetail = append(dtoResponse.TransactionNotesDetail, detailNotes)

					// clear
					deepDetailsNotes = []dtos.TransactionNotesDeepDetail{}
					detailNotes = dtos.TransactionNotesDetail{}
				} else {
					deepDetailsNotes = append(deepDetailsNotes, dtos.TransactionNotesDeepDetail{
						TransactionNote: v.TransactionNote,
						TransactionAmount: dtos.Amount{
							CurrencyCode: "IDR",
							Value:        v.Amount,
						},
						TransactionLimit: dtos.Amount{
							CurrencyCode: "IDR",
							Value:        v.Budget,
						},
					})
				}
			}

			if catPrev == "" {
				catPrev = v.TransactionCategory
				deepDetailsNotes = append(deepDetailsNotes, dtos.TransactionNotesDeepDetail{
					TransactionNote: v.TransactionNote,
					TransactionAmount: dtos.Amount{
						CurrencyCode: "IDR",
						Value:        v.Amount,
					},
					TransactionLimit: dtos.Amount{
						CurrencyCode: "IDR",
						Value:        v.Budget,
					},
				})

				if k == lengthData-1 {
					// save previous
					detailNotes.TransactionCategory = catPrev
					detailNotes.TransactionNotesDeepDetail = append(detailNotes.TransactionNotesDeepDetail, deepDetailsNotes...)

					dtoResponse.TransactionNotesDetail = append(dtoResponse.TransactionNotesDetail, detailNotes)
				}
			}

			if catPrev != v.TransactionCategory {

				if k == lengthData-1 {

					// save previous
					detailNotes.TransactionCategory = catPrev
					detailNotes.TransactionNotesDeepDetail = append(detailNotes.TransactionNotesDeepDetail, deepDetailsNotes...)

					dtoResponse.TransactionNotesDetail = append(dtoResponse.TransactionNotesDetail, detailNotes)

					// clear
					deepDetailsNotes = []dtos.TransactionNotesDeepDetail{}
					detailNotes = dtos.TransactionNotesDetail{}

					deepDetailsNotes = append(deepDetailsNotes, dtos.TransactionNotesDeepDetail{
						TransactionNote: v.TransactionNote,
						TransactionAmount: dtos.Amount{
							CurrencyCode: "IDR",
							Value:        v.Amount,
						},
						TransactionLimit: dtos.Amount{
							CurrencyCode: "IDR",
							Value:        v.Budget,
						},
					})

					detailNotes.TransactionCategory = v.TransactionCategory
					detailNotes.TransactionNotesDeepDetail = append(detailNotes.TransactionNotesDeepDetail, deepDetailsNotes...)
					dtoResponse.TransactionNotesDetail = append(dtoResponse.TransactionNotesDetail, detailNotes)

				} else {
					// save previous
					detailNotes.TransactionCategory = catPrev
					detailNotes.TransactionNotesDeepDetail = append(detailNotes.TransactionNotesDeepDetail, deepDetailsNotes...)

					dtoResponse.TransactionNotesDetail = append(dtoResponse.TransactionNotesDetail, detailNotes)

					// clear
					deepDetailsNotes = []dtos.TransactionNotesDeepDetail{}
					detailNotes = dtos.TransactionNotesDetail{}

					// save new
					deepDetailsNotes = append(deepDetailsNotes, dtos.TransactionNotesDeepDetail{
						TransactionNote: v.TransactionNote,
						TransactionAmount: dtos.Amount{
							CurrencyCode: "IDR",
							Value:        v.Amount,
						},
						TransactionLimit: dtos.Amount{
							CurrencyCode: "IDR",
							Value:        v.Budget,
						},
					})

					catPrev = v.TransactionCategory
				}
			}
		}

		dtoResponse.TransactionDate = month + "-" + year

	}

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	return dtoResponse, http.StatusOK, errInfo
}

func (s *TransactionUseCase) Suggestion(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		dataResponse         []entities.TransactionSuggestionNotes
		suggestionCollection []string
		err                  error
	)

	usrEmail := ctx.MustGet("email").(string)
	personalAccount := personalaccounts.Informations(ctx, usrEmail)

	if personalAccount.ID == uuid.Nil {
		httpCode = http.StatusNotFound
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "token contains invalid information")
		return []string{}, httpCode, errInfo
	}

	dataResponse, err = s.repo.Suggestion(personalAccount.ID)
	if err != nil {
		logrus.Error(err.Error())
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return []string{}, http.StatusInternalServerError, errInfo
	}

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	if len(dataResponse) == 0 {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "data not found")
		return []string{}, http.StatusBadRequest, errInfo
	}

	if len(dataResponse) > 0 {
		for _, v := range dataResponse {
			suggestionCollection = append(suggestionCollection, v.Note)
		}
	}

	return suggestionCollection, http.StatusOK, errInfo
}

func (s *TransactionUseCase) CashFlow(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {

	var dtoResponse dtos.CashFlowResponse

	// get account uuid
	accountUUID := ctx.MustGet("accountID").(uuid.UUID)

	dataIncomeEachDay, err := s.repo.AverageIncomeEachDay(accountUUID)
	if err != nil {
		logrus.Error(err.Error())
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	dataExpenseEachDay, err := s.repo.AverageExpenseEachDay(accountUUID)
	if err != nil {
		logrus.Error(err.Error())
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	dataIncomeMonthly, err := s.repo.AverageIncomeMonthly(accountUUID)
	if err != nil {
		logrus.Error(err.Error())
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	dataExpenseMonthly, err := s.repo.AverageExpenseMonthly(accountUUID)
	if err != nil {
		logrus.Error(err.Error())
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	dataCountIncome, err := s.repo.CountIncomeTransaction(accountUUID)
	if err != nil {
		logrus.Error(err.Error())
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	dataCountExpense, err := s.repo.CountExpenseTransaction(accountUUID)
	if err != nil {
		logrus.Error(err.Error())
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	dataTotalIncome, err := s.repo.DataTotalIncome(accountUUID)
	if err != nil {
		logrus.Error(err.Error())
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	dataTotalExpense, err := s.repo.DataTotalExpense(accountUUID)
	if err != nil {
		logrus.Error(err.Error())
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	dtoResponse.CountIncome = dataCountExpense.CountExpense
	dtoResponse.CountExpense = dataCountIncome.CountIncome
	dtoResponse.AverageDay.Income = dataIncomeEachDay.IncomeAverage
	dtoResponse.AverageDay.Expense = dataExpenseEachDay.ExpenseAverage
	dtoResponse.AverageMonth.Income = dataIncomeMonthly.IncomeAverage
	dtoResponse.AverageMonth.Expense = dataExpenseMonthly.ExpenseAverage
	dtoResponse.TotalAverageIncome = dataTotalIncome.TotalIncome
	dtoResponse.TotalAverageExpense = dataTotalExpense.TotalExpense
	dtoResponse.CashFlow = dtoResponse.TotalAverageIncome - dtoResponse.TotalAverageExpense

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	return dtoResponse, http.StatusOK, errInfo
}

func (s *TransactionUseCase) saveInvestTransaction(accountID uuid.UUID, request *dtos.TransactionRequestInvestment) (ID uuid.UUID, err error) {

	// convert process from string type to uuid type
	IDWalletUUID, err := uuid.Parse(request.IDWallet)
	if err != nil {
		logrus.Error(err.Error())
	}

	IDMasterBrokerUUID, err := uuid.Parse(request.IDMasterBroker)
	if err != nil {
		logrus.Error(err.Error())
	}

	IDMasterInvestUUID, err := uuid.Parse(request.IDMasterInvest)
	if err != nil {
		logrus.Error(err.Error())
	}

	IDMasterTrxTypesUUID, err := uuid.Parse(request.IDMasterTransactionTypes)
	if err != nil {
		logrus.Error(err.Error())
	}

	// mapping to entities process
	trxID := uuid.New()
	modelTransaction := entities.TransactionEntity{
		ID:                       trxID,
		Date:                     request.Date,
		Amount:                   float64(request.Price),
		IDPersonalAccount:        accountID,
		IDWallet:                 IDWalletUUID,
		IDMasterBroker:           IDMasterBrokerUUID,
		IDMasterInvest:           IDMasterInvestUUID,
		IDMasterTransactionTypes: IDMasterTrxTypesUUID,
	}

	modelTransactionDetail := entities.TransactionDetailEntity{
		IDTransaction: trxID,
		StockCode:     request.StockCode,
		Lot:           request.Lot,
		SellBuy:       request.SellBuy,
	}

	// save transaction into transaction table, and transaction detail table, return err if any and id
	return trxID, s.repo.Add(&modelTransaction, &modelTransactionDetail)
}

func (s *TransactionUseCase) investmentCalculation(accountID uuid.UUID) (err error) {
	var (
		lotBuy            int64
		lotSell           int64
		initialInvestment float64
		buy               float64
		sell              float64
		percentageReturn  float64
		sellCollection    []float64
		feeSell           float64
		feeBuy            float64
		netBuy            float64
		averageBuy        float64
		gainloss          float64
		potentialReturn   float64
		brokerName        string
		stockCode         string
		maxData           int
		IDMasterBroker    uuid.UUID
	)

	// get transaction detail table where id_personal_accounts
	dataTrxDetail, err := s.repo.AllInvestmentsTrx(accountID)
	if err != nil {
		logrus.Error(err.Error())
		return err
	}

	if dataTrxDetail == nil {
		logrus.Error(err.Error())
		return errors.New("system can not calculate investment process")
	}

	// sorting
	sort.Sort(dataTrxDetail)

	// length of data
	maxData = len(dataTrxDetail) - 1

	// do complex calculation
	for k, v := range dataTrxDetail {

		// get broker info
		brokerInfo, err := s.repo.GetBrokerInfo(v.IDMasterBroker)
		if err != nil {
			logrus.Error(err.Error())
		}

		tradingInfo, err := s.repo.GetTradingInfo(v.StockCode)
		{
			if err != nil {
				logrus.Error(err.Error())
			}
		}

		// if broker name is same as previous
		if brokerName == brokerInfo.Name {

			// if stock code same as previous
			if stockCode == v.StockCode {

				// if buy transaction
				if v.SellBuy == 1 {
					buy = float64(v.Lot) * v.Amount * 100
					initialInvestment += buy
					lotBuy += v.Lot
				}

				// if sell transaction
				if v.SellBuy == 0 {
					// sell calculation
					sell = float64(v.Lot) * v.Amount * 100
					sellCollection = append(sellCollection, sell)
					lotSell += v.Lot

				}

				// if latest data
				if k == maxData {
					// average buy
					averageBuy = (initialInvestment / float64(lotBuy)) * 100

					// fee buy calculation
					feeBuy = v.FeeBuy * buy

					// net buy
					netBuy = buy - feeBuy

					// potential return and percentage
					potentialReturn += (float64(tradingInfo.Close) - float64(averageBuy)) * float64(lotBuy) * 100
					percentageReturn = potentialReturn / initialInvestment

					if len(sellCollection) > 0 {
						// buy calculation
						buy = averageBuy * float64(lotSell) * 100
						feeBuy = v.FeeBuy * buy
						netBuy = buy - feeBuy

						// sell calculation
						totalSell := 0.0
						for _, v := range sellCollection {
							totalSell += v
						}
						feeSell = totalSell * v.FeeSell
						netSell := totalSell - feeSell

						gainloss = netSell - netBuy
						percentageReturn = gainloss / buy
					}

					// mapping
					trxInvestment := entities.TransactionInvestmentEntity{
						StockCode:         stockCode,
						TotalLot:          lotBuy - lotSell,
						ValueBuy:          buy,
						FeeBuy:            feeBuy,
						NetBuy:            netBuy,
						AverageBuy:        averageBuy,
						InitialInvestment: initialInvestment,
						IDPersonalAccount: accountID,
						IDMasterBroker:    IDMasterBroker,
						GainLoss:          gainloss,
						PotentialReturn:   potentialReturn,
						PercentageReturn:  percentageReturn,
					}

					// save data into investment table
					err = s.repo.RecordInvestTrx(&trxInvestment)
					if err != nil {
						logrus.Error(err.Error())
					}
				}
			}

			// if stock code different than previous
			if stockCode != v.StockCode {

				// average buy
				averageBuy = (initialInvestment / float64(lotBuy)) * 100

				// fee buy calculation
				feeBuy = v.FeeBuy * buy

				// net buy
				netBuy = buy - feeBuy

				// potential return and percentage
				potentialReturn += (float64(tradingInfo.Close) - averageBuy) * float64(lotBuy) * 100
				percentageReturn = potentialReturn / initialInvestment

				if len(sellCollection) > 0 {
					// buy calculation
					buy = averageBuy * float64(lotSell) * 100
					feeBuy = v.FeeBuy * buy
					netBuy = buy - feeBuy

					// sell calculation
					totalSell := 0.0
					for _, v := range sellCollection {
						totalSell += v
					}
					feeSell = totalSell * v.FeeSell
					netSell := totalSell - feeSell

					gainloss = netSell - netBuy
					percentageReturn = gainloss / buy
				}

				// mapping
				trxInvestment := entities.TransactionInvestmentEntity{
					StockCode:         stockCode,
					TotalLot:          lotBuy - lotSell,
					ValueBuy:          buy,
					FeeBuy:            feeBuy,
					NetBuy:            netBuy,
					AverageBuy:        averageBuy,
					InitialInvestment: initialInvestment,
					IDPersonalAccount: accountID,
					IDMasterBroker:    IDMasterBroker,
					GainLoss:          0,
					PotentialReturn:   potentialReturn,
					PercentageReturn:  percentageReturn,
				}

				// save previous data into investment table
				err = s.repo.RecordInvestTrx(&trxInvestment)
				if err != nil {
					logrus.Error(err.Error())
				}

				// clear
				initialInvestment = 0
				lotBuy = 0
				lotSell = 0
				potentialReturn = 0

				// set
				stockCode = v.StockCode

				// if buy transaction
				if v.SellBuy == 1 {
					buy = float64(v.Lot) * v.Amount * 100
					initialInvestment += buy
					lotBuy += v.Lot
				}

				// if latest data
				if k == maxData {
					// average buy
					averageBuy = (initialInvestment / float64(lotBuy)) * 100

					// fee buy calculation
					feeBuy = v.FeeBuy * buy

					// net buy
					netBuy = buy - feeBuy

					// potential return and percentage
					potentialReturn += (float64(tradingInfo.Close) - averageBuy) * float64(lotBuy) * 100
					percentageReturn = potentialReturn / initialInvestment

					// mapping
					trxInvestment = entities.TransactionInvestmentEntity{
						StockCode:         stockCode,
						TotalLot:          lotBuy - lotSell,
						ValueBuy:          buy,
						AverageBuy:        averageBuy,
						FeeBuy:            feeBuy,
						NetBuy:            netBuy,
						InitialInvestment: initialInvestment,
						IDPersonalAccount: accountID,
						IDMasterBroker:    IDMasterBroker,
						GainLoss:          0,
						PotentialReturn:   potentialReturn,
						PercentageReturn:  percentageReturn,
					}

					// save data into investment table
					err = s.repo.RecordInvestTrx(&trxInvestment)
					if err != nil {
						logrus.Error(err.Error())
					}
				}
			}
		}

		// if first time, set broker name, set stock code
		if brokerName == "" {
			// set
			brokerName = brokerInfo.Name
			stockCode = v.StockCode
			lotBuy += v.Lot
			IDMasterBroker = v.IDMasterBroker

			// if buy transaction
			if v.SellBuy == 1 {
				buy = float64(v.Lot) * v.Amount * 100
				initialInvestment += buy
			}

			// if latest data
			if k == maxData {
				// average buy
				averageBuy = (initialInvestment / float64(lotBuy)) * 100

				// fee buy calculation
				feeBuy = v.FeeBuy * buy

				// net buy
				netBuy = buy - feeBuy

				// potential return and percentage
				potentialReturn += (float64(tradingInfo.Close) - averageBuy) * float64(lotBuy) * 100
				percentageReturn = potentialReturn / initialInvestment

				// mapping
				trxInvestment := entities.TransactionInvestmentEntity{
					StockCode:         stockCode,
					TotalLot:          lotBuy - lotSell,
					ValueBuy:          buy,
					AverageBuy:        averageBuy,
					FeeBuy:            feeBuy,
					NetBuy:            netBuy,
					InitialInvestment: initialInvestment,
					IDPersonalAccount: accountID,
					IDMasterBroker:    IDMasterBroker,
					GainLoss:          0,
					PotentialReturn:   potentialReturn,
					PercentageReturn:  percentageReturn,
				}

				// save data into investment table
				err = s.repo.RecordInvestTrx(&trxInvestment)
				if err != nil {
					logrus.Error(err.Error())
				}
			}
		}

		// if broker name different than previous
		if brokerName != brokerInfo.Name {
			// average buy
			averageBuy = (initialInvestment / float64(lotBuy)) * 100

			// fee buy calculation
			feeBuy = v.FeeBuy * buy

			// net buy
			netBuy = buy - feeBuy

			// potential return and percentage
			potentialReturn += (float64(tradingInfo.Close) - averageBuy) * float64(lotBuy) * 100
			percentageReturn = potentialReturn / initialInvestment

			// mapping
			trxInvestment := entities.TransactionInvestmentEntity{
				StockCode:         stockCode,
				TotalLot:          lotBuy - lotSell,
				ValueBuy:          buy,
				AverageBuy:        averageBuy,
				FeeBuy:            feeBuy,
				NetBuy:            netBuy,
				InitialInvestment: initialInvestment,
				IDPersonalAccount: accountID,
				IDMasterBroker:    IDMasterBroker,
				GainLoss:          0,
				PotentialReturn:   potentialReturn,
				PercentageReturn:  percentageReturn,
			}

			// save previous data into investment table
			err = s.repo.RecordInvestTrx(&trxInvestment)
			if err != nil {
				logrus.Error(err.Error())
			}

			// clear
			initialInvestment = 0
			lotBuy = 0
			lotSell = 0
			potentialReturn = 0

			// renew value
			brokerName = brokerInfo.Name
			stockCode = v.StockCode
			IDMasterBroker = v.IDMasterBroker

			// if buy transaction
			if v.SellBuy == 1 {
				buy = float64(v.Lot) * v.Amount * 100
				initialInvestment += buy
				lotBuy += v.Lot
			}

			// if latest data
			if k == maxData {
				// average buy
				averageBuy = (initialInvestment / float64(lotBuy)) * 100

				// fee buy calculation
				feeBuy = v.FeeBuy * buy

				// net buy
				netBuy = buy - feeBuy

				// potential return and percentage
				potentialReturn += (float64(tradingInfo.Close) - averageBuy) * float64(lotBuy) * 100
				percentageReturn = potentialReturn / initialInvestment

				// mapping
				trxInvestment := entities.TransactionInvestmentEntity{
					StockCode:         stockCode,
					TotalLot:          lotBuy - lotSell,
					ValueBuy:          buy,
					AverageBuy:        averageBuy,
					FeeBuy:            feeBuy,
					NetBuy:            netBuy,
					InitialInvestment: initialInvestment,
					IDPersonalAccount: accountID,
					IDMasterBroker:    IDMasterBroker,
					GainLoss:          0,
					PotentialReturn:   potentialReturn,
					PercentageReturn:  percentageReturn,
				}

				// save data into investment table
				err = s.repo.RecordInvestTrx(&trxInvestment)
				if err != nil {
					logrus.Error(err.Error())
				}
			}
		}
	}

	return nil
}

func (s *TransactionUseCase) WalletNonInvestment(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var dtoResponse []dtos.WalletListResponse

	// accountUUID
	accountUUID := ctx.MustGet("accountID").(uuid.UUID)

	// fetch data wallet
	dataWallet, err := s.repo.WalletNonInvestment(accountUUID)
	if err != nil {
		logrus.Error()
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	// if data wallet
	if len(dataWallet) == 0 {
		resp := struct {
			Message string `json:"message"`
		}{
			Message: "no data for wallet non investment",
		}
		return resp, http.StatusNotFound, []errorsinfo.Errors{}
	}

	// if data wallet exist
	if len(dataWallet) > 0 {
		for _, v := range dataWallet {
			dtoResponse = append(dtoResponse, dtos.WalletListResponse{
				IDAccount: accountUUID,
				WalletDetails: dtos.WalletDetails{
					WalletID:           v.ID,
					WalletType:         v.WalletType,
					WalletName:         v.WalletName,
					IDMasterWalletType: v.IDMasterWalletType,
				},
			})
		}
	}

	return dtoResponse, http.StatusOK, errInfo
}

func (s *TransactionUseCase) WalletInvestment(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var dtoResponse []dtos.WalletListResponse

	// accountUUID
	accountUUID := ctx.MustGet("accountID").(uuid.UUID)

	dataWallet, err := s.repo.WalletInvestment(accountUUID)
	if err != nil {
		logrus.Error(err)
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return struct{}{}, http.StatusInternalServerError, errInfo
	}

	if len(dataWallet) == 0 {
		resp := struct {
			Message string `json:"message"`
		}{
			Message: "no data for wallet investment",
		}
		return resp, http.StatusNotFound, []errorsinfo.Errors{}
	}

	if len(dataWallet) > 0 {
		for _, v := range dataWallet {
			dtoResponse = append(dtoResponse, dtos.WalletListResponse{
				IDAccount: accountUUID,
				WalletDetails: dtos.WalletDetails{
					WalletID:           v.ID,
					WalletType:         v.WalletType,
					WalletName:         v.WalletName,
					IDMasterWalletType: v.IDMasterWalletType,
				},
			})
		}
	}

	return dtoResponse, http.StatusOK, errInfo
}