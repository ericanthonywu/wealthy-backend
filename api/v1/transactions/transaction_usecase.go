package transactions

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/transactions/dtos"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/transactions/entities"
	"github.com/semicolon-indonesia/wealthy-backend/utils/errorsinfo"
	"github.com/semicolon-indonesia/wealthy-backend/utils/personalaccounts"
	"github.com/semicolon-indonesia/wealthy-backend/utils/utilities"
	"github.com/sirupsen/logrus"
	"net/http"
)

type (
	TransactionUseCase struct {
		repo ITransactionRepository
	}

	ITransactionUseCase interface {
		Add(ctx *gin.Context, request *dtos.TransactionRequest) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		ExpenseTransactionHistory(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		IncomeTransactionHistory(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		TravelTransactionHistory(ctx *gin.Context, IDTravel uuid.UUID) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		TransferTransactionHistory(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		InvestTransactionHistory(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		IncomeSpending(ctx *gin.Context, month string, year string) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		Investment(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		ByNotes(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		Suggestion(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
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
	)

	accountUUID := ctx.MustGet("accountID").(uuid.UUID)

	// convert string to UUID
	IDWalletUUID, err := uuid.Parse(request.IDWallet)
	if err != nil {
		logrus.Error(err.Error())
	}

	IDMasterInvestUUID, err := uuid.Parse(request.IDMasterInvest)
	if err != nil {
		logrus.Error(err.Error())
	}

	IDMasterBrokerUUID, err := uuid.Parse(request.IDMasterBroker)
	if err != nil {
		logrus.Error(err.Error())
	}

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

	if request.IDTravel == "" {
		convertAmount = request.Amount
		IDTravelUUID = uuid.Nil

		// is wallet true exists
		if !s.repo.WalletExist(IDWalletUUID) {
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id wallet unregistered before")
			return struct{}{}, http.StatusBadRequest, errInfo
		}

		// check for wallet type
	}

	// translate string to uuid
	if request.IDMasterIncomeCategories != "" {
		IDMasterIncCatUUID, err = uuid.Parse(request.IDMasterIncomeCategories)
		if err != nil {
			logrus.Error(err.Error())
		}
	}

	if request.IDMasterExpenseCategories != "" {
		IDMasterExpCatUUID, err = uuid.Parse(request.IDMasterExpenseCategories)
		if err != nil {
			logrus.Error(err.Error())
		}
	}

	if request.IDMasterExpenseSubCategories != "" {
		IDMasterSubExpCatUUID, err = uuid.Parse(request.IDMasterExpenseSubCategories)
		if err != nil {
			logrus.Error(err.Error())
		}
	}

	if request.IDMasterTransactionPriorities != "" {
		IDMasterTransPriUUID, err = uuid.Parse(request.IDMasterTransactionPriorities)
		if err != nil {
			logrus.Error(err.Error())
		}
	}

	if request.IDMasterTransactionTypes != "" {
		IDMasterTransTypeUUID, err = uuid.Parse(request.IDMasterTransactionTypes)
		if err != nil {
			logrus.Error(err.Error())
		}
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
		IDMasterInvest:                IDMasterInvestUUID,
		IDMasterBroker:                IDMasterBrokerUUID,
		IDMasterTransactionPriorities: IDMasterTransPriUUID,
		IDMasterTransactionTypes:      IDMasterTransTypeUUID,
	}

	modelTransactionDetail := entities.TransactionDetailEntity{
		IDTransaction:     trxID,
		Repeat:            request.Repeat,
		Note:              request.Note,
		From:              request.TransferFrom,
		To:                request.TransferTo,
		MutualFundProduct: request.MutualFundProduct,
		StockCode:         request.StockCode,
		Lot:               request.Lot,
		SellBuy:           request.SellBuy,
		IDTravel:          IDTravelUUID,
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

func (s *TransactionUseCase) ExpenseTransactionHistory(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		dtoResponse                  dtos.TransactionHistoryForIncomeExpenses
		responseExpenseTotalHistory  entities.TransactionExpenseTotalHistory
		responseExpenseDetailHistory interface{}
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
		responseExpenseTotalHistory = s.repo.ExpenseTotalHistoryWithoutDate(personalAccount.ID)
		responseExpenseDetailHistory = s.repo.ExpenseDetailHistoryWithoutDate(personalAccount.ID)
	} else {
		responseExpenseTotalHistory = s.repo.ExpenseTotalHistoryWithDate(personalAccount.ID, startDate, endDate)
		responseExpenseDetailHistory = s.repo.ExpenseDetailHistoryWithDate(personalAccount.ID, startDate, endDate)
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

	usrEmail := ctx.MustGet("email").(string)
	personalAccount := personalaccounts.Informations(ctx, usrEmail)

	if personalAccount.ID == uuid.Nil {
		httpCode = http.StatusNotFound
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "token contains invalid information")
		return response, httpCode, errInfo
	}

	if startDate == "" || endDate == "" {
		responseTravelDetailHistory = s.repo.TravelDetailWithoutData(personalAccount.ID, IDTravel)
	} else {
		responseTravelDetailHistory = s.repo.TravelDetailWithData(personalAccount.ID, IDTravel, startDate, endDate)
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
				Category: v.Category,
				Note:     v.Note,
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

		isNotExist := responseIncomeSpendingDetailAnnually[0].NetIncome == 0 && responseIncomeSpendingDetailAnnually[0].TotalIncome == 0 && responseIncomeSpendingDetailAnnually[0].TotalSpending == 0
		if isNotExist {
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
		dtoResponse             dtos.TransactionIncomeSpendingInvestment
		responseInvestmentTotal interface{}
		//responseInvestmentDetail interface{}
	)

	month := ctx.Query("month")
	year := ctx.Query("year")

	usrEmail := ctx.MustGet("email").(string)
	personalAccount := personalaccounts.Informations(ctx, usrEmail)

	if personalAccount.ID == uuid.Nil {
		httpCode = http.StatusNotFound
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "token contains invalid information")
		return response, httpCode, errInfo
	}

	if month == "" && year == "" {
		httpCode = http.StatusBadRequest
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "need month and or year information")
		return response, httpCode, errInfo
	}

	if month != "" && year != "" {
		responseInvestmentTotal = s.repo.InvestMonthlyTotal(personalAccount.ID, month, year)
		//responseInvestmentDetail = s.repo.InvestMonthlyDetail(personalAccount.ID, month, year)
	}

	if month == "" && year != "" {
		responseInvestmentTotal = s.repo.InvestAnnuallyTotal(personalAccount.ID, year)
		//responseInvestmentDetail = s.repo.InvestAnnuallyDetail(personalAccount.ID, year)
	}

	dtoResponse.Summary = responseInvestmentTotal
	//dtoResponse.Detail = responseInvestmentDetail

	return dtoResponse, http.StatusOK, []errorsinfo.Errors{}
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