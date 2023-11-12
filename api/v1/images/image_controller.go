package images

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type (
	ShowImageController struct {
		useCase IShowImageUseCase
	}

	IShowImageController interface {
		Avatar(ctx *gin.Context)
	}
)

func NewShowImageController(useCase IShowImageUseCase) *ShowImageController {
	return &ShowImageController{useCase: useCase}
}

func (c *ShowImageController) Avatar(ctx *gin.Context) {
	fileName := ctx.Param("filename")
	imagePath := fmt.Sprintf("./assets/avatar/%s", fileName)
	ctx.File(imagePath)
}