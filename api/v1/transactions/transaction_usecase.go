package transactions

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/transactions/dtos"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/transactions/entities"
	"github.com/semicolon-indonesia/wealthy-backend/utils/errorsinfo"
	"github.com/semicolon-indonesia/wealthy-backend/utils/personalaccounts"
	"github.com/semicolon-indonesia/wealthy-backend/utils/utilities"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

type (
	TransactionUseCase struct {
		repo ITransactionRepository
	}

	ITransactionUseCase interface {
		Add(ctx *gin.Context, request *dtos.TransactionRequest) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		ExpenseTransactionHistory(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		IncomeTransactionHistory(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		TravelTransactionHistory(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
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
		trxID      uuid.UUID
		err        error
		filename   string
		targetPath string
		imagePath  string
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

	if request.ImageBase64 != "" {
		imageData, err := base64.StdEncoding.DecodeString(request.ImageBase64)
		filename = fmt.Sprintf("%d", time.Now().Unix()) + ".png"
		targetPath = "assets/travel/" + filename
		imagePath = "images/travel/" + filename

		err = utilities.SaveImage(imageData, targetPath)
		if err != nil {
			logrus.Error(err.Error())
		}
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
		Departure:         request.Departure,
		Arrival:           request.Arrival,
		ImagePath:         imagePath,
		Filename:          filename,
		TravelStartDate:   request.TravelStartDate,
		TravelEndDate:     request.TravelEndDate,
		SellBuy:           request.SellBuy,
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

func (s *TransactionUseCase) TravelTransactionHistory(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
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
		responseTravelDetailHistory = s.repo.TravelDetailWithoutData(personalAccount.ID)
	} else {
		responseTravelDetailHistory = s.repo.TravelDetailWithData(personalAccount.ID, startDate, endDate)
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
				Departure: v.Departure,
				Arrival:   v.Arrival,
				Amount: dtos.Amount{
					CurrencyCode: "IDR",
					Value:        v.Amount,
				},
				ImagePath:       os.Getenv("APP_HOST") + "/v1/" + v.ImagePath,
				Filename:        v.Filename,
				TravelStartDate: v.TravelStartDate,
				TravelEndDate:   v.TravelEndDate,
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
		dtoResponse                  dtos.TransactionIncomeSpendingInvestment
		responseIncomeSpendingTotal  interface{}
		responseIncomeSpendingDetail interface{}
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
		httpCode = http.StatusBadGateway
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "need month information")
		return response, httpCode, errInfo
	}

	if month != "" && year != "" {
		responseIncomeSpendingTotal = s.repo.IncomeSpendingMonthlyTotal(personalAccount.ID, month, year)
		responseIncomeSpendingDetail = s.repo.IncomeSpendingMonthlyDetail(personalAccount.ID, month, year)
	}

	if month == "" && year != "" {
		responseIncomeSpendingTotal = s.repo.IncomeSpendingAnnuallyTotal(personalAccount.ID, year)
		responseIncomeSpendingDetail = s.repo.IncomeSpendingAnnuallyDetail(personalAccount.ID, year)
	}

	dtoResponse.Detail = responseIncomeSpendingDetail
	dtoResponse.Summary = responseIncomeSpendingTotal
	return dtoResponse, http.StatusOK, []errorsinfo.Errors{}
}

func (s *TransactionUseCase) Investment(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		dtoResponse              dtos.TransactionIncomeSpendingInvestment
		responseInvestmentTotal  interface{}
		responseInvestmentDetail interface{}
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
		responseInvestmentDetail = s.repo.InvestMonthlyDetail(personalAccount.ID, month, year)
	}

	if month == "" && year != "" {
		responseInvestmentTotal = s.repo.InvestAnnuallyTotal(personalAccount.ID, year)
		responseInvestmentDetail = s.repo.InvestAnnuallyDetail(personalAccount.ID, year)
	}

	dtoResponse.Summary = responseInvestmentTotal
	dtoResponse.Detail = responseInvestmentDetail

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