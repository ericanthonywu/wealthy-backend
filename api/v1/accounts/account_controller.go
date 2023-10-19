package accounts

import (
	"github.com/gin-gonic/gin"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/accounts/dtos"
	"github.com/semicolon-indonesia/wealthy-backend/utils/errorsinfo"
	"github.com/semicolon-indonesia/wealthy-backend/utils/response"
	"net/http"
)

type (
	AccountController struct {
		useCase IAccountUseCase
	}

	IAccountController interface {
		SignUp(ctx *gin.Context)
		SignIn(ctx *gin.Context)
		SignOut(ctx *gin.Context)
		Profile(ctx *gin.Context)
	}
)

func NewAccountController(useCase IAccountUseCase) *AccountController {
	return &AccountController{useCase: useCase}
}

func (c *AccountController) SignUp(ctx *gin.Context) {
	var (
		dtoRequest  dtos.AccountSignUpRequest
		dtoResponse dtos.AccountSignUpResponse
		httpCode    int
		errInfo     []errorsinfo.Errors
	)

	if err := ctx.ShouldBindJSON(&dtoRequest); err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "no body payload")
		response.SendBack(ctx, dtos.AccountSignUpResponse{}, errInfo, http.StatusBadRequest)
		return
	}

	dtoResponse, httpCode, errInfo = c.useCase.SignUp(&dtoRequest)
	response.SendBack(ctx, dtoResponse, errInfo, httpCode)
	return

}

func (c *AccountController) SignIn(ctx *gin.Context) {
	var (
		dtoRequest  dtos.AccountSignInRequest
		dtoResponse dtos.AccountSignInResponse
		httpCode    int
		errInfo     []errorsinfo.Errors
	)

	if err := ctx.ShouldBindJSON(&dtoRequest); err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "no body payload")
		response.SendBack(ctx, dtos.AccountSignUpResponse{}, errInfo, http.StatusBadRequest)
		return
	}

	dtoResponse, httpCode, errInfo = c.useCase.SignIn(&dtoRequest)
	response.SendBack(ctx, dtoResponse, errInfo, httpCode)
	return
}

func (c *AccountController) SignOut(ctx *gin.Context) {
	c.useCase.SignOut()
}

func (c *AccountController) Profile(ctx *gin.Context) {
	data, httpCode, errInfo := c.useCase.Profile(ctx)

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	response.SendBack(ctx, data, errInfo, httpCode)
	return
}
