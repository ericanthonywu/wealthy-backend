package transactions

import (
	"github.com/gin-gonic/gin"
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
		Expense(ctx *gin.Context)
		Income(ctx *gin.Context)
		Transfer(ctx *gin.Context)
		Invest(ctx *gin.Context)
	}
)

func NewTransactionController(useCase ITransactionUseCase) *TransactionController {
	return &TransactionController{useCase: useCase}
}

func (c *TransactionController) Expense(ctx *gin.Context) {
	var (
		dtoRequest dtos.TransactionExpenseRequest
		errInfo    []errorsinfo.Errors
	)

	if err := ctx.ShouldBindJSON(&dtoRequest); err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "no body payload")
		response.SendBack(ctx, dtos.TransactionExpenseRequest{}, errInfo, http.StatusBadRequest)
		return
	}

	c.useCase.Expense(&dtoRequest)
	return
}

func (c *TransactionController) Income(ctx *gin.Context) {

	return
}

func (c *TransactionController) Transfer(ctx *gin.Context) {
	return
}

func (c *TransactionController) Invest(ctx *gin.Context) {
	return
}
