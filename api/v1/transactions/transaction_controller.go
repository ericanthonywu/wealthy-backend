package transactions

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/transactions/dtos"
	"github.com/semicolon-indonesia/wealthy-backend/constants"
	"github.com/semicolon-indonesia/wealthy-backend/utils/errorsinfo"
	"github.com/semicolon-indonesia/wealthy-backend/utils/response"
	"github.com/sirupsen/logrus"
	"net/http"
)

type (
	TransactionController struct {
		useCase ITransactionUseCase
	}

	ITransactionController interface {
		Add(ctx *gin.Context)
		AddInvestmentTransaction(ctx *gin.Context)
		ExpenseTransactionHistory(ctx *gin.Context)
		IncomeTransactionHistory(ctx *gin.Context)
		TransferTransactionHistory(ctx *gin.Context)
		IncomeSpending(ctx *gin.Context)
		Investment(ctx *gin.Context)
		ByNotes(ctx *gin.Context)
		TravelTransactionHistory(ctx *gin.Context)
		Suggestion(ctx *gin.Context)
		CashFlow(ctx *gin.Context)
		validateTravelTransactionPayload(request *dtos.TransactionRequest) (errInfo []errorsinfo.Errors)
		validateIncomeTransactionPayload(request *dtos.TransactionRequest) (errInfo []errorsinfo.Errors)
		validateExpenseTransactionPayload(request *dtos.TransactionRequest) (errInfo []errorsinfo.Errors)
		validateInvestTransactionPayload(request *dtos.TransactionRequest) (errInfo []errorsinfo.Errors)
	}
)

func NewTransactionController(useCase ITransactionUseCase) *TransactionController {
	return &TransactionController{useCase: useCase}
}

func (c *TransactionController) Add(ctx *gin.Context) {
	var (
		dtoRequest dtos.TransactionRequest
		errInfo    []errorsinfo.Errors
	)

	// bind
	if err := ctx.ShouldBindJSON(&dtoRequest); err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "no body payload")
		response.SendBack(ctx, struct{}{}, errInfo, http.StatusBadRequest)
		return
	}

	// validation
	if dtoRequest.Date == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "date transaction empty value")
	}

	if dtoRequest.IDMasterTransactionTypes == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id master transaction type empty value")
	}

	// for travel transaction
	if dtoRequest.IDMasterTransactionTypes == constants.TravelTrx {
		// get account type
		accountType := ctx.MustGet("accountType").(string)

		// if basic account
		if accountType == constants.AccountBasic {
			resp := struct {
				Message string `json:"message"`
			}{
				Message: constants.ProPlan,
			}
			response.SendBack(ctx, resp, []errorsinfo.Errors{}, http.StatusUpgradeRequired)
			return
		}

		errInfo = c.validateTravelTransactionPayload(&dtoRequest)

		// response error
		if len(errInfo) > 0 {
			response.SendBack(ctx, struct{}{}, errInfo, http.StatusBadRequest)
			return
		}
	}

	// for income transaction
	if dtoRequest.IDMasterIncomeCategories == constants.IncomeTrx {
		errInfo = c.validateIncomeTransactionPayload(&dtoRequest)

		// response error
		if len(errInfo) > 0 {
			response.SendBack(ctx, struct{}{}, errInfo, http.StatusBadRequest)
			return
		}
	}

	// for expense transaction
	if dtoRequest.IDMasterExpenseCategories == constants.ExpenseTrx {
		errInfo = c.validateExpenseTransactionPayload(&dtoRequest)

		// response error
		if len(errInfo) > 0 {
			response.SendBack(ctx, struct{}{}, errInfo, http.StatusBadRequest)
			return
		}
	}

	// send back with err information
	if len(errInfo) > 0 {
		response.SendBack(ctx, struct{}{}, errInfo, http.StatusBadRequest)
		return
	}

	data, httpCode, errInfo := c.useCase.Add(ctx, &dtoRequest)
	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *TransactionController) AddInvestmentTransaction(ctx *gin.Context) {
	var (
		request dtos.TransactionRequestInvestment
		errInfo []errorsinfo.Errors
	)

	// get account type
	accountType := ctx.MustGet("accountType").(string)

	// if basic account
	if accountType == constants.AccountBasic {
		resp := struct {
			Message string `json:"message"`
		}{
			Message: constants.ProPlan,
		}
		response.SendBack(ctx, resp, []errorsinfo.Errors{}, http.StatusUpgradeRequired)
		return
	}

	// bind
	if err := ctx.ShouldBindJSON(&request); err != nil {
		errInfo := errorsinfo.ErrorWrapper(errInfo, "", "no body payload")
		response.SendBack(ctx, struct{}{}, errInfo, http.StatusBadRequest)
		return
	}

	// validate first
	errInfo = c.validateInvestTransactionPayload(&request)

	// if error occurs
	if len(errInfo) > 0 {
		response.SendBack(ctx, struct{}{}, errInfo, http.StatusBadRequest)
		return
	}

	// do logic and return
	data, httpCode, errInfo := c.useCase.AddInvestmentTransaction(ctx, &request)
	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *TransactionController) ExpenseTransactionHistory(ctx *gin.Context) {
	data, httpCode, errInfo := c.useCase.ExpenseTransactionHistory(ctx)

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *TransactionController) IncomeTransactionHistory(ctx *gin.Context) {
	data, httpCode, errInfo := c.useCase.IncomeTransactionHistory(ctx)

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *TransactionController) TransferTransactionHistory(ctx *gin.Context) {
	data, httpCode, errInfo := c.useCase.TransferTransactionHistory(ctx)

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *TransactionController) InvestTransactionHistory(ctx *gin.Context) {
	data, httpCode, errInfo := c.useCase.InvestTransactionHistory(ctx)

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *TransactionController) IncomeSpending(ctx *gin.Context) {
	var errInfo []errorsinfo.Errors

	month := ctx.Query("month")
	year := ctx.Query("year")

	if year == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "year query param needed in url address")
	}

	if len(errInfo) > 0 {
		response.SendBack(ctx, struct{}{}, errInfo, http.StatusBadRequest)
		return
	}

	data, httpCode, errInfo := c.useCase.IncomeSpending(ctx, month, year)
	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *TransactionController) Investment(ctx *gin.Context) {
	data, httpCode, errInfo := c.useCase.Investment(ctx)
	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *TransactionController) ByNotes(ctx *gin.Context) {
	data, httpCode, errInfo := c.useCase.ByNotes(ctx)
	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *TransactionController) TravelTransactionHistory(ctx *gin.Context) {
	var errInfo []errorsinfo.Errors

	// get account type
	accountType := ctx.MustGet("accountType").(string)

	// if basic account
	if accountType == constants.AccountBasic {
		resp := struct {
			Message string `json:"message"`
		}{
			Message: constants.ProPlan,
		}
		response.SendBack(ctx, resp, []errorsinfo.Errors{}, http.StatusUpgradeRequired)
		return
	}

	idTravel := ctx.Query("idTravel")

	if idTravel == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "idTravel query param needed in url address")
	}

	IDTravelUUID, err := uuid.Parse(idTravel)
	if err != nil {
		logrus.Error(err.Error())
	}

	if len(errInfo) > 0 {
		response.SendBack(ctx, struct{}{}, errInfo, http.StatusBadRequest)
		return
	}

	data, httpCode, errInfo := c.useCase.TravelTransactionHistory(ctx, IDTravelUUID)
	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *TransactionController) Suggestion(ctx *gin.Context) {
	data, httpCode, errInfo := c.useCase.Suggestion(ctx)
	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *TransactionController) CashFlow(ctx *gin.Context) {
	data, httCode, errInfo := c.useCase.CashFlow(ctx)
	response.SendBack(ctx, data, errInfo, httCode)
	return
}

func (c *TransactionController) validateTravelTransactionPayload(request *dtos.TransactionRequest) (errInfo []errorsinfo.Errors) {
	// mandatory field
	if request.Amount == 0 {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "amount have greater than 0")
	}

	if request.IDTravel == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id travel empty value. get one from budget API")
	}

	if request.IDMasterExpenseCategories == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id master expense category empty value")
	}

	// unnecessary field
	if request.IDMasterIncomeCategories != "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id master income categories unnecessary for travel transaction")
	}

	if request.IDWallet != "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id wallet unnecessary for travel transaction")
	}

	if request.IDMasterExpenseSubCategories != "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id master expense sub categories unnecessary for travel transaction")
	}

	return errInfo
}

func (c *TransactionController) validateIncomeTransactionPayload(request *dtos.TransactionRequest) (errInfo []errorsinfo.Errors) {
	// mandatory field

	if request.Amount == 0 {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "amount have greater than 0")
	}

	if request.IDWallet == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id wallet empty value")
	}

	if request.IDMasterIncomeCategories == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id master income categories empty value")
	}

	// unnecessary field
	if request.IDMasterExpenseCategories == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id master expense category unnecessary for income transaction")
	}

	if request.IDMasterExpenseSubCategories == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id master expense sub category unnecessary for income transaction")
	}

	if request.IDTravel == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id travel unnecessary for income transaction")
	}

	return errInfo
}

func (c *TransactionController) validateExpenseTransactionPayload(request *dtos.TransactionRequest) (errInfo []errorsinfo.Errors) {
	return
}
func (c *TransactionController) validateInvestTransactionPayload(request *dtos.TransactionRequestInvestment) (errInfo []errorsinfo.Errors) {

	if request.IDWallet == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id wallet empty value")
	}

	if request.IDMasterBroker == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id master broker empty value")
	}

	if request.IDMasterInvest == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id master invest empty value")
	}

	if request.IDMasterTransactionTypes == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id master transaction type empty value")
	}

	if request.Lot <= 0 {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "lot must greater than 0")
	}

	if request.Price <= 0 {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "price must greater than 0")
	}

	if request.SellBuy < 0 || request.SellBuy > 2 {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "sell buy must one of two values : 0 [ sell ] or 1 [ buy ] ")
	}

	if request.StockCode == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "stock code empty value")
	}

	return errInfo
}