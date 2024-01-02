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
		GetCatagoriesList(ginContext *gin.Context)
	}
)

func NewCategoryHandler(useCase ICategoryUseCase) *CategoryHandler {
	return &CategoryHandler{useCase: useCase}
}

func (c *CategoryHandler) GetCatagoriesList(ginContext *gin.Context) {
	data, httpCode, errInfo := c.useCase.GetCatagoriesList(ginContext)
	response.SendBack(ginContext, data, errInfo, httpCode)
	return
}