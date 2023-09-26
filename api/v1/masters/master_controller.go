package masters

import (
	"github.com/gin-gonic/gin"
	"github.com/semicolon-indonesia/wealthy-backend/utils/errorsinfo"
	"github.com/semicolon-indonesia/wealthy-backend/utils/response"
	"net/http"
)

type (
	MasterController struct {
		useCase IMasterUseCase
	}

	IMasterController interface {
		TransactionType(ctx *gin.Context)
		IncomeType(ctx *gin.Context)
		ExpenseType(ctx *gin.Context)
		ReksadanaType(ctx *gin.Context)
	}
)

func NewMasterController(useCase IMasterUseCase) *MasterController {
	return &MasterController{useCase: useCase}
}

func (c *MasterController) TransactionType(ctx *gin.Context) {
	data := c.useCase.TransactionType()
	response.SendBack(ctx, data, []errorsinfo.Errors{}, http.StatusOK)
	return
}

func (c *MasterController) IncomeType(ctx *gin.Context) {
	data := c.useCase.IncomeType()
	response.SendBack(ctx, data, []errorsinfo.Errors{}, http.StatusOK)
	return
}

func (c *MasterController) ExpenseType(ctx *gin.Context) {
	data := c.useCase.ExpenseType()
	response.SendBack(ctx, data, []errorsinfo.Errors{}, http.StatusOK)
	return
}

func (c *MasterController) ReksadanaType(ctx *gin.Context) {
	data := c.useCase.ReksadanaType()
	response.SendBack(ctx, data, []errorsinfo.Errors{}, http.StatusOK)
	return
}
