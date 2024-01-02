package categories

import (
	"github.com/gin-gonic/gin"
	"github.com/wealthy-app/wealthy-backend/utils/response"
)

type (
	CategoryHandler struct {
		useCase ICategoryUseCase
	}

	ICategoryHandler interface {
		GetCategoriesExpenseList(ginContext *gin.Context)
		GetCategoriesIncomeList(ginContext *gin.Context)
	}
)

func NewCategoryHandler(useCase ICategoryUseCase) *CategoryHandler {
	return &CategoryHandler{useCase: useCase}
}

func (c *CategoryHandler) GetCategoriesExpenseList(ginContext *gin.Context) {
	data, httpCode, errInfo := c.useCase.GetCategoriesExpenseList(ginContext)
	response.SendBack(ginContext, data, errInfo, httpCode)
	return
}

func (c *CategoryHandler) GetCategoriesIncomeList(ginContext *gin.Context) {
	data, httpCode, errInfo := c.useCase.GetCategoriesIncomeList(ginContext)
	response.SendBack(ginContext, data, errInfo, httpCode)
	return
}