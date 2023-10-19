package budgets

import (
	"github.com/gin-gonic/gin"
	"github.com/semicolon-indonesia/wealthy-backend/utils/errorsinfo"
	"github.com/semicolon-indonesia/wealthy-backend/utils/response"
	"net/http"
)

type (
	BudgetController struct {
		useCase IBudgetUseCase
	}

	IBudgetController interface {
		AllCategories(ctx *gin.Context)
		Overview(ctx *gin.Context)
		Category(ctx *gin.Context)
		LatestSixMonths(ctx *gin.Context)
		Set(ctx *gin.Context)
	}
)

func NewBudgetController(useCase IBudgetUseCase) *BudgetController {
	return &BudgetController{useCase: useCase}
}

func (c *BudgetController) AllCategories(ctx *gin.Context) {
	c.useCase.AllCategories(ctx)
	response.SendBack(ctx, nil, []errorsinfo.Errors{}, http.StatusOK)
	return
}

func (c *BudgetController) Overview(ctx *gin.Context) {
	data, httpCode, errInfo := c.useCase.Overview(ctx)

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

func (c *BudgetController) Set(ctx *gin.Context) {
	c.useCase.Set()
}
