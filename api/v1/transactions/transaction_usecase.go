package transactions

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/transactions/dtos"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/transactions/entities"
	"github.com/semicolon-indonesia/wealthy-backend/utils/errorsinfo"
	"github.com/semicolon-indonesia/wealthy-backend/utils/personalaccounts"
	"github.com/semicolon-indonesia/wealthy-backend/utils/utilities"
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
		IncomeSpending(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		Investment(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		ByNotes(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
	}
)

func NewTransactionUseCase(repo ITransactionRepository) *TransactionUseCase {
	return &TransactionUseCase{repo: repo}
}

func (s *TransactionUseCase) Add(ctx *gin.Context, request *dtos.TransactionRequest) (data interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		trxID uuid.UUID
		err   error
	)

	usrEmail := ctx.MustGet("email").(string)
	personalAccount := personalaccounts.Informations(ctx, usrEmail)

	if personalAccount.ID == uuid.Nil {
		httpCode = http.StatusUnauthorized
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "token contains invalid information")
		return nil, httpCode, errInfo
	}

	trxID, err = uuid.NewUUID()
	if err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return nil, http.StatusUnprocessableEntity, errInfo
	}

	modelTransaction := entities.TransactionEntity{
		ID:                            trxID,
		Date:                          request.Date,
		Fees:                          float64(request.Fees),
		Amount:                        float64(request.Amount),
		IDPersonalAccount:             personalAccount.ID,
		IDWallet:                      request.IDWallet,
		IDMasterIncomeCategories:      request.IDMasterIncomeCategories,
		IDMasterExpenseCategories:     request.IDMasterExpenseCategories,
		IDMasterExpenseSubCategories:  request.IDMasterExpenseSubCategories,
		IDMasterInvest:                request.IDMasterInvest,
		IDMasterBroker:                request.IDMasterBroker,
		IDMasterReksanadaTypes:        request.IDMasterReksanadaTypes,
		IDMasterTransactionPriorities: request.IDMasterTransactionPriorities,
		IDMasterTransactionTypes:      request.IDMasterTransactionTypes,
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
		IDTravel:          request.IDTravel,
	}

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
	return data, http.StatusOK, []errorsinfo.Errors{}
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

	if len(responseTravelDetailHistory) == 0 {
		response := struct {
			Message string `json:"message"`
		}{
			Message: "there is not travel transaction between periods : " + startDate + " until " + endDate,
		}
		return response, http.StatusNotFound, []errorsinfo.Errors{}
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

func (s *TransactionUseCase) IncomeSpending(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
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
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "need month information")
		return response, http.StatusBadGateway, errInfo
	}

	if month != "" && year != "" {
		responseIncomeSpendingTotal = s.repo.IncomeSpendingMonthlyTotal(personalAccount.ID, month, year)
		responseIncomeSpendingDetailMonthly = s.repo.IncomeSpendingMonthlyDetail(personalAccount.ID, month, year)
	}

	if month == "" && year != "" {
		responseIncomeSpendingTotal = s.repo.IncomeSpendingAnnuallyTotal(personalAccount.ID, year)
		responseIncomeSpendingDetailAnnually = s.repo.IncomeSpendingAnnuallyDetail(personalAccount.ID, year)
	}

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
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

	return response, http.StatusPreconditionFailed, []errorsinfo.Errors{}
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
		httpCode = http.StatusBadGateway
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "need month and or year information")
		return response, httpCode, errInfo
	}

	if month != "" && year != "" {
		response = s.repo.ByNote(personalAccount.ID, month, year)
	}

	return response, http.StatusOK, []errorsinfo.Errors{}
}