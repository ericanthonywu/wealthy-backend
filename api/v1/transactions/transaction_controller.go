package transactions

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/wealthy-app/wealthy-backend/api/v1/transactions/dtos"
	"github.com/wealthy-app/wealthy-backend/constants"
	"github.com/wealthy-app/wealthy-backend/utils/datecustoms"
	"github.com/wealthy-app/wealthy-backend/utils/errorsinfo"
	"github.com/wealthy-app/wealthy-backend/utils/response"
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
		WalletNonInvestment(ctx *gin.Context)
		WalletInvestment(ctx *gin.Context)
		validateTravelTransactionPayload(request *dtos.TransactionRequest) (errInfo []errorsinfo.Errors)
		validateIncomeTransactionPayload(request *dtos.TransactionRequest) (errInfo []errorsinfo.Errors)
		validateExpenseTransactionPayload(request *dtos.TransactionRequest) (errInfo []errorsinfo.Errors)
		validateInvestTransactionPayload(request *dtos.TransactionRequestInvestment) (errInfo []errorsinfo.Errors)
		validateTransferTransactionPayload(request *dtos.TransactionRequest) (errInfo []errorsinfo.Errors)
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

	if dtoRequest.Date != "" {
		if !datecustoms.ValidDateFormat(dtoRequest.Date) {
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", "date transaction must format YYYY-MM-DD")
		}

		if datecustoms.TotalDaysBetweenDate(dtoRequest.Date) < 0 {
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", "date transaction not permit for future transaction record")
		}
	}

	if dtoRequest.IDMasterTransactionTypes == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id master transaction type empty value")
	}

	switch dtoRequest.IDMasterTransactionTypes {
	case constants.TravelTrx:
		{
			// get account type
			accountType := ctx.MustGet("accountType").(string)

			// check account type
			if accountType == constants.AccountBasic {
				resp := struct {
					Message string `json:"message"`
				}{
					Message: constants.ProPlan,
				}
				response.SendBack(ctx, resp, []errorsinfo.Errors{}, http.StatusUpgradeRequired)
				return
			}

			// validate for travel transaction
			errInfoValidate := c.validateTravelTransactionPayload(&dtoRequest)
			errInfo = append(errInfo, errInfoValidate...)
			break
		}
	case constants.IncomeTrx:
		{
			// validate for income transaction
			errInfoValidate := c.validateIncomeTransactionPayload(&dtoRequest)
			errInfo = append(errInfo, errInfoValidate...)
			break
		}
	case constants.ExpenseTrx:
		{
			// validate for expense transaction
			errInfoValidate := c.validateExpenseTransactionPayload(&dtoRequest)
			errInfo = append(errInfo, errInfoValidate...)
			break
		}
	case constants.TransferTrx:
		{
			// validate for transfer transaction
			errInfoValidate := c.validateTransferTransactionPayload(&dtoRequest)
			errInfo = append(errInfo, errInfoValidate...)
			break
		}
	default:
		{
			errInfoValidate := errorsinfo.ErrorWrapper(errInfo, "", "unknown transaction types")
			errInfo = append(errInfo, errInfoValidate...)
			break
		}
	}

	// response error
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

	if request.IDMasterTransactionPriorities != "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id master transaction priorities unnecessary for travel transaction")
	}

	if request.IDMasterExpenseSubCategories != "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id master expense sub categories unnecessary for travel transaction")
	}

	if request.TransferFrom != "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "transfer form unnecessary for travel transaction")
	}

	if request.TransferTo != "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "transfer to unnecessary for travel transaction")
	}

	if request.Fees < 0 || request.Fees > 0 {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "fee unnecessary for travel transaction")
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
	if request.IDMasterExpenseCategories != "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id master expense category unnecessary for income transaction")
	}

	if request.IDMasterExpenseSubCategories != "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id master expense sub category unnecessary for income transaction")
	}

	if request.IDTravel != "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id travel unnecessary for income transaction")
	}

	if request.IDMasterTransactionPriorities != "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id master transaction priorities unnecessary for income transaction")
	}

	if request.TransferFrom != "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "transfer form unnecessary for income transaction")
	}

	if request.TransferTo != "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "transfer to unnecessary for income transaction")
	}

	if request.Fees > 0 || request.Fees < 0 {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "fee unnecessary for income transaction")
	}

	return errInfo
}

func (c *TransactionController) validateExpenseTransactionPayload(request *dtos.TransactionRequest) (errInfo []errorsinfo.Errors) {
	// mandatory field
	if request.Amount == 0 {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "amount have greater than 0")
	}

	if request.IDWallet == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id wallet empty value")
	}

	if request.IDMasterExpenseCategories == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id master expense category empty value")
	}

	if request.IDMasterExpenseSubCategories == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id master expense sub category empty value")
	}

	if request.IDMasterTransactionPriorities == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id master transaction priorities empty value")
	}

	// unnecessary field
	if request.IDMasterIncomeCategories != "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id master income categories unnecessary for expense transaction")
	}

	if request.IDTravel != "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id travel unnecessary unnecessary for expense transaction")
	}

	if request.TransferFrom != "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "transfer form unnecessary for expense transaction")
	}

	if request.TransferTo != "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "transfer to unnecessary for expense transaction")
	}

	if request.Fees > 0 || request.Fees < 0 {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "fee unnecessary for expense transaction")
	}

	return errInfo
}

func (c *TransactionController) validateInvestTransactionPayload(request *dtos.TransactionRequestInvestment) (errInfo []errorsinfo.Errors) {

	if request.IDWallet == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id wallet empty value")
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

func (c *TransactionController) validateTransferTransactionPayload(request *dtos.TransactionRequest) (errInfo []errorsinfo.Errors) {
	// mandatory field
	if request.Amount == 0 {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "amount have greater than 0")
	}

	if request.TransferFrom == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "transfer form empty value")
	}

	if request.TransferTo == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "transfer to empty value")
	}

	if request.Fees <= 0 {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "fee must greater than 0")
	}

	// unnecessary field
	if request.IDMasterExpenseCategories != "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id master expense category unnecessary for transfer transaction")
	}

	if request.IDMasterExpenseSubCategories != "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id master expense sub category unnecessary for transfer transaction")
	}

	if request.IDMasterTransactionPriorities != "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id master transaction priorities unnecessary for transfer transaction")
	}

	if request.IDMasterIncomeCategories != "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id master income categories unnecessary for transfer transaction")
	}

	if request.IDTravel != "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id travel unnecessary unnecessary for transfer transaction")
	}

	if request.IDWallet != "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id wallet unnecessary unnecessary for transfer transaction")
	}

	return errInfo
}

func (c *TransactionController) WalletNonInvestment(ctx *gin.Context) {
	data, httpCode, errInfo := c.useCase.WalletNonInvestment(ctx)
	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *TransactionController) WalletInvestment(ctx *gin.Context) {
	data, httpCode, errInfo := c.useCase.WalletInvestment(ctx)
	response.SendBack(ctx, data, errInfo, httpCode)
	return
}