package masters

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/semicolon-indonesia/wealthy-backend/utils/errorsinfo"
	"github.com/semicolon-indonesia/wealthy-backend/utils/response"
	"github.com/sirupsen/logrus"
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
		WalletType(ctx *gin.Context)
		Invest(ctx *gin.Context)
		Broker(ctx *gin.Context)
		TransactionPriority(ctx *gin.Context)
		Gender(ctx gin.Context)
		SubExpenseCategories(ctx *gin.Context)
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

func (c *MasterController) WalletType(ctx *gin.Context) {
	data := c.useCase.WalletType()
	response.SendBack(ctx, data, []errorsinfo.Errors{}, http.StatusOK)
	return
}

func (c *MasterController) Invest(ctx *gin.Context) {
	data := c.useCase.InvestType()
	response.SendBack(ctx, data, []errorsinfo.Errors{}, http.StatusOK)
	return
}

func (c *MasterController) Broker(ctx *gin.Context) {
	data := c.useCase.Broker()
	response.SendBack(ctx, data, []errorsinfo.Errors{}, http.StatusOK)
	return
}

func (c *MasterController) TransactionPriority(ctx *gin.Context) {
	data := c.useCase.TransactionPriority()
	response.SendBack(ctx, data, []errorsinfo.Errors{}, http.StatusOK)
	return
}

func (c *MasterController) Gender(ctx *gin.Context) {
	data := c.useCase.Gender()
	response.SendBack(ctx, data, []errorsinfo.Errors{}, http.StatusOK)
	return
}

func (c *MasterController) SubExpenseCategories(ctx *gin.Context) {
	var (
		errInfo       []errorsinfo.Errors
		expenseIDUUID uuid.UUID
		err           error
	)

	expenseID := ctx.Param("expense-id")

	if expenseID == "" {
		errInfo := errorsinfo.ErrorWrapper(errInfo, "", "expense ID require in url parameter")
		response.SendBack(ctx, nil, errInfo, http.StatusBadRequest)
		return
	}

	expenseIDUUID, err = uuid.Parse(expenseID)
	if err != nil {
		logrus.Error(err.Error())
		errInfo := errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		response.SendBack(ctx, nil, errInfo, http.StatusBadRequest)
		return
	}

	data := c.useCase.SubExpenseCategories(expenseIDUUID)
	if data == nil {
		errInfo := errorsinfo.ErrorWrapper(errInfo, "", "expense ID maybe not registered")
		response.SendBack(ctx, data, errInfo, http.StatusBadRequest)
		return
	}

	response.SendBack(ctx, data, []errorsinfo.Errors{}, http.StatusOK)
	return
}