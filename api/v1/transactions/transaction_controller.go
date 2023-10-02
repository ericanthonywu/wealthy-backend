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
		Add(ctx *gin.Context)
	}
)

func NewTransactionController(useCase ITransactionUseCase) *TransactionController {
	return &TransactionController{useCase: useCase}
}

func (c *TransactionController) Add(ctx *gin.Context) {
	var (
		dtoRequest dtos.TransactionRequest
		errInfo    []errorsinfo.Errors
		httpCode   int
		data       interface{}
	)

	if err := ctx.ShouldBindJSON(&dtoRequest); err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "no body payload")
		response.SendBack(ctx, dtos.TransactionRequest{}, errInfo, http.StatusBadRequest)
		return
	}

	data, httpCode, errInfo = c.useCase.Add(&dtoRequest)
	response.SendBack(ctx, data, []errorsinfo.Errors{}, httpCode)
	return
}
