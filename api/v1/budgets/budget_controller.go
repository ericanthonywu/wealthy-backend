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
		Set(ctx *gin.Context)
	}
)

func NewBudgetController(useCase IBudgetUseCase) *BudgetController {
	return &BudgetController{useCase: useCase}
}

func (c *BudgetController) AllCategories(ctx *gin.Context) {
	//data := c.useCase.AllCategories()
	response.SendBack(ctx, nil, []errorsinfo.Errors{}, http.StatusOK)
}

func (c *BudgetController) Set(ctx *gin.Context) {
	var ()
	c.useCase.Set()

}
