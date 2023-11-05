package budgets

import (
	"github.com/gin-gonic/gin"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/budgets/dtos"
	"github.com/semicolon-indonesia/wealthy-backend/utils/errorsinfo"
	"github.com/semicolon-indonesia/wealthy-backend/utils/response"
	"net/http"
)

type (
	BudgetController struct {
		useCase IBudgetUseCase
	}

	IBudgetController interface {
		AllLimit(ctx *gin.Context)
		Overview(ctx *gin.Context)
		Category(ctx *gin.Context)
		LatestSixMonths(ctx *gin.Context)
		Limit(ctx *gin.Context)
	}
)

func NewBudgetController(useCase IBudgetUseCase) *BudgetController {
	return &BudgetController{useCase: useCase}
}

func (c *BudgetController) AllLimit(ctx *gin.Context) {
	data, httpCode, errInfo := c.useCase.AllLimit(ctx)

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}
	response.SendBack(ctx, data, []errorsinfo.Errors{}, httpCode)
	return
}

func (c *BudgetController) Overview(ctx *gin.Context) {

	var errInfo []errorsinfo.Errors

	month := ctx.Query("month")
	year := ctx.Query("year")

	if month == "" || year == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "month or year required in query url")
		response.SendBack(ctx, nil, errInfo, http.StatusBadRequest)
		return
	}

	data, httpCode, errInfo := c.useCase.Overview(ctx, month, year)

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *BudgetController) Category(ctx *gin.Context) {
	data, httpCode, errInfo := c.useCase.Category(ctx)

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *BudgetController) LatestSixMonths(ctx *gin.Context) {
	data, httpCode, errInfo := c.useCase.LatestSixMonths(ctx)

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *BudgetController) Limit(ctx *gin.Context) {
	var (
		dtoRequest  dtos.BudgetSetRequest
		dtoResponse interface{}
		errInfo     []errorsinfo.Errors
		httpCode    int
	)

	if err := ctx.ShouldBindJSON(&dtoRequest); err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "body payload required")
		response.SendBack(ctx, dtos.BudgetSetRequest{}, errInfo, http.StatusBadRequest)
		return
	}

	dtoResponse, httpCode, errInfo = c.useCase.Limit(ctx, &dtoRequest)
	response.SendBack(ctx, dtoResponse, errInfo, httpCode)
}