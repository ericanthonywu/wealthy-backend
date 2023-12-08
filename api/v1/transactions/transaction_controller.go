package transactions

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/transactions/dtos"
	"github.com/semicolon-indonesia/wealthy-backend/utils/errorsinfo"
	"github.com/semicolon-indonesia/wealthy-backend/utils/response"
	"net/http"
)

type (
	TransactionController struct {
		useCase ITransactionUseCase
	}

	ITransactionController interface {
		Add(ctx *gin.Context)
		ExpenseTransactionHistory(ctx *gin.Context)
		IncomeTransactionHistory(ctx *gin.Context)
		TransferTransactionHistory(ctx *gin.Context)
		IncomeSpending(ctx *gin.Context)
		Investment(ctx *gin.Context)
		ByNotes(ctx *gin.Context)
		TravelTransactionHistory(ctx *gin.Context)
		Suggestion(ctx *gin.Context)
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

	if dtoRequest.IDWallet == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id wallet empty value")
	}

	// transaction for investments
	if dtoRequest.IDMasterInvest != "" || dtoRequest.IDMasterBroker != "" {
		if dtoRequest.StockCode == "" {
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", "stock code empty value")
		}

		if dtoRequest.Lot == 0 {
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", "lot must greater than 0 value")
		}

		if dtoRequest.IDMasterInvest == "" {
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id master invest empty value")
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
	data, httpCode, errInfo := c.useCase.IncomeSpending(ctx)

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *TransactionController) Investment(ctx *gin.Context) {
	data, httpCode, errInfo := c.useCase.Investment(ctx)

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *TransactionController) ByNotes(ctx *gin.Context) {
	data, httpCode, errInfo := c.useCase.ByNotes(ctx)

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *TransactionController) TravelTransactionHistory(ctx *gin.Context) {
	var errInfo []errorsinfo.Errors

	idTravel := ctx.Query("idTravel")

	if idTravel == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "idTravel query param needed in url address")
		response.SendBack(ctx, interface{}(nil), errInfo, http.StatusBadRequest)
		return
	}

	IDTravelUUID, _ := uuid.Parse(idTravel)
	data, httpCode, errInfo := c.useCase.TravelTransactionHistory(ctx, IDTravelUUID)

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *TransactionController) Suggestion(ctx *gin.Context) {
	data, httpCode, errInfo := c.useCase.Suggestion(ctx)
	response.SendBack(ctx, data, errInfo, httpCode)
	return
}