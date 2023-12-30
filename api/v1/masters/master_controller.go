package masters

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/wealthy-app/wealthy-backend/api/v1/masters/dtos"
	"github.com/wealthy-app/wealthy-backend/constants"
	"github.com/wealthy-app/wealthy-backend/utils/errorsinfo"
	"github.com/wealthy-app/wealthy-backend/utils/response"
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
		Gender(ctx *gin.Context)
		SubExpenseCategories(ctx *gin.Context)
		Exchange(ctx *gin.Context)
		PersonalIncomeCategory(ctx *gin.Context)
		PersonalExpenseCategory(ctx *gin.Context)
		PersonalExpenseSubCategory(ctx *gin.Context)
		RenameIncomeCategory(ctx *gin.Context)
		RenameExpenseCategory(ctx *gin.Context)
		RenameSubExpenseCategory(ctx *gin.Context)
		AddIncomeCategory(ctx *gin.Context)
		AddExpenseCategory(ctx *gin.Context)
		AddSubExpenseCategory(ctx *gin.Context)
		Price(ctx *gin.Context)
		StockCode(ctx *gin.Context)
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

func (c *MasterController) Exchange(ctx *gin.Context) {
	data, httpCode, errInfo := c.useCase.Exchange()
	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *MasterController) PersonalIncomeCategory(ctx *gin.Context) {
	data, httpCode, errInfo := c.useCase.PersonalIncomeCategory(ctx)
	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *MasterController) PersonalExpenseCategory(ctx *gin.Context) {
	data, httpCode, errInfo := c.useCase.PersonalExpenseCategory(ctx)
	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *MasterController) PersonalExpenseSubCategory(ctx *gin.Context) {
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

	data, httpCode, errInfo := c.useCase.PersonalExpenseSubCategory(ctx, expenseIDUUID)
	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *MasterController) RenameIncomeCategory(ctx *gin.Context) {
	var (
		dtoRequest dtos.RenameCatRequest
		errInfo    []errorsinfo.Errors
		idUUID     uuid.UUID
		err        error
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
	if err := ctx.ShouldBindJSON(&dtoRequest); err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "body payload required")
		response.SendBack(ctx, dtos.RenameCatRequest{}, errInfo, http.StatusBadRequest)
		return
	}

	id := ctx.Param("id")

	if id == "" {
		errInfo := errorsinfo.ErrorWrapper(errInfo, "", "expense ID require in url parameter")
		response.SendBack(ctx, nil, errInfo, http.StatusBadRequest)
		return
	}

	idUUID, err = uuid.Parse(id)
	if err != nil {
		logrus.Error(err.Error())
		errInfo := errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		response.SendBack(ctx, nil, errInfo, http.StatusBadRequest)
		return
	}

	data, httpCode, errInfo := c.useCase.RenameIncomeCategory(ctx, idUUID, &dtoRequest)
	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *MasterController) RenameExpenseCategory(ctx *gin.Context) {
	var (
		dtoRequest dtos.RenameCatRequest
		errInfo    []errorsinfo.Errors
		idUUID     uuid.UUID
		err        error
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
	if err := ctx.ShouldBindJSON(&dtoRequest); err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "body payload required")
		response.SendBack(ctx, dtos.RenameCatRequest{}, errInfo, http.StatusBadRequest)
		return
	}

	id := ctx.Param("id")

	if id == "" {
		errInfo := errorsinfo.ErrorWrapper(errInfo, "", "expense ID require in url parameter")
		response.SendBack(ctx, nil, errInfo, http.StatusBadRequest)
		return
	}

	idUUID, err = uuid.Parse(id)
	if err != nil {
		logrus.Error(err.Error())
		errInfo := errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		response.SendBack(ctx, nil, errInfo, http.StatusBadRequest)
		return
	}

	data, httpCode, errInfo := c.useCase.RenameExpenseCategory(ctx, idUUID, &dtoRequest)
	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *MasterController) RenameSubExpenseCategory(ctx *gin.Context) {
	var (
		dtoRequest dtos.RenameCatRequest
		errInfo    []errorsinfo.Errors
		idUUID     uuid.UUID
		err        error
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
	if err := ctx.ShouldBindJSON(&dtoRequest); err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "body payload required")
		response.SendBack(ctx, dtos.RenameCatRequest{}, errInfo, http.StatusBadRequest)
		return
	}

	id := ctx.Param("id")

	if id == "" {
		errInfo := errorsinfo.ErrorWrapper(errInfo, "", "expense ID require in url parameter")
		response.SendBack(ctx, nil, errInfo, http.StatusBadRequest)
		return
	}

	idUUID, err = uuid.Parse(id)
	if err != nil {
		logrus.Error(err.Error())
		errInfo := errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		response.SendBack(ctx, nil, errInfo, http.StatusBadRequest)
		return
	}

	data, httpCode, errInfo := c.useCase.RenameSubExpenseCategory(ctx, idUUID, &dtoRequest)
	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *MasterController) AddIncomeCategory(ctx *gin.Context) {
	var (
		dtoRequest dtos.AddCategory
		errInfo    []errorsinfo.Errors
		err        error
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
	if err = ctx.ShouldBindJSON(&dtoRequest); err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "body payload required")
		response.SendBack(ctx, dtos.RenameCatRequest{}, errInfo, http.StatusBadRequest)
		return
	}

	data, httpCode, errInfo := c.useCase.AddIncomeCategory(ctx, &dtoRequest)
	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *MasterController) AddExpenseCategory(ctx *gin.Context) {
	var (
		dtoRequest dtos.AddCategory
		errInfo    []errorsinfo.Errors
		err        error
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
	if err = ctx.ShouldBindJSON(&dtoRequest); err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "body payload required")
		response.SendBack(ctx, dtos.RenameCatRequest{}, errInfo, http.StatusBadRequest)
		return
	}

	data, httpCode, errInfo := c.useCase.AddExpenseCategory(ctx, &dtoRequest)
	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *MasterController) AddSubExpenseCategory(ctx *gin.Context) {
	var (
		dtoRequest dtos.AddCategory
		errInfo    []errorsinfo.Errors
		err        error
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
	if err = ctx.ShouldBindJSON(&dtoRequest); err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "body payload required")
		response.SendBack(ctx, dtos.RenameCatRequest{}, errInfo, http.StatusBadRequest)
		return
	}

	data, httpCode, errInfo := c.useCase.AddSubExpenseCategory(ctx, &dtoRequest)
	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *MasterController) Price(ctx *gin.Context) {
	data, httpCode, errInfo := c.useCase.Price(ctx)
	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *MasterController) StockCode(ctx *gin.Context) {
	data := c.useCase.StockCode()
	response.SendBack(ctx, data, []errorsinfo.Errors{}, http.StatusOK)
	return
}