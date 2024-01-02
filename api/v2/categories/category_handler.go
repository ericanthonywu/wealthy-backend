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
		GetCatagoriesExpenseList(ginContext *gin.Context)
	}
)

func NewCategoryHandler(useCase ICategoryUseCase) *CategoryHandler {
	return &CategoryHandler{useCase: useCase}
}

func (c *CategoryHandler) GetCatagoriesExpenseList(ginContext *gin.Context) {
	data, httpCode, errInfo := c.useCase.GetCatagoriesExpenseList(ginContext)
	response.SendBack(ginContext, data, errInfo, httpCode)
	return
}