package internals

import (
	"github.com/gin-gonic/gin"
	"github.com/semicolon-indonesia/wealthy-backend/utils/errorsinfo"
	"github.com/semicolon-indonesia/wealthy-backend/utils/response"
)

type (
	InternalController struct {
		useCase IInternalUseCase
	}

	IInternalController interface {
		TransactionNotes(ctx *gin.Context)
	}
)

func NewInternalController(useCase IInternalUseCase) *InternalController {
	return &InternalController{useCase: useCase}
}

func (c *InternalController) TransactionNotes(ctx *gin.Context) {
	data, httpCode, errInfo := c.useCase.TransactionNotes(ctx)

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	response.SendBack(ctx, data, errInfo, httpCode)
	return
}