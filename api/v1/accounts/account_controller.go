package accounts

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/accounts/dtos"
	"github.com/semicolon-indonesia/wealthy-backend/utils/errorsinfo"
	"github.com/semicolon-indonesia/wealthy-backend/utils/response"
	"github.com/sirupsen/logrus"
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
		GetProfile(ctx *gin.Context)
		UpdateProfile(ctx *gin.Context)
		ChangePassword(ctx *gin.Context)
		ValidateRefCode(ctx *gin.Context)
		SetAvatar(ctx *gin.Context)
		RemoveAvatar(ctx *gin.Context)
		Sharing(ctx *gin.Context)
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

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}
	response.SendBack(ctx, dtoResponse, errInfo, httpCode)
	return
}

func (c *AccountController) SignOut(ctx *gin.Context) {
	c.useCase.SignOut()
}

func (c *AccountController) GetProfile(ctx *gin.Context) {
	data, httpCode, errInfo := c.useCase.GetProfile(ctx)

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *AccountController) UpdateProfile(ctx *gin.Context) {
	var (
		dtoRequest  map[string]interface{}
		dtoResponse map[string]bool
		errInfo     []errorsinfo.Errors
		httpCode    int
	)

	customerID := ctx.Param("id")
	if customerID == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id in url address is required")
		response.SendBack(ctx, dtos.AccountSetProfileRequest{}, errInfo, http.StatusBadRequest)
		return
	}

	if err := ctx.ShouldBindJSON(&dtoRequest); err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "body payload required")
		response.SendBack(ctx, dtos.AccountSetProfileRequest{}, errInfo, http.StatusBadRequest)
		return
	}

	customerIDUUID, err := uuid.Parse(customerID)
	if err != nil {
		logrus.Error(err.Error())
	}

	dtoResponse, httpCode, errInfo = c.useCase.UpdateProfile(ctx, customerIDUUID, dtoRequest)
	response.SendBack(ctx, dtoResponse, errInfo, httpCode)
	return
}

func (c *AccountController) ChangePassword(ctx *gin.Context) {
	var (
		dtoRequest   dtos.AccountChangePassword
		errInfo      []errorsinfo.Errors
		custUUID     uuid.UUID
		err          error
		httpCode     int
		dataResponse = make(map[string]bool)
	)

	dataResponse["success"] = false

	customerID := ctx.Param("id")
	if customerID == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "customer id required in url parameter")
		response.SendBack(ctx, dataResponse, errInfo, http.StatusBadRequest)
		return
	}

	if err := ctx.ShouldBindJSON(&dtoRequest); err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "body payload required")
		response.SendBack(ctx, dataResponse, errInfo, http.StatusBadRequest)
		return
	}

	custUUID, err = uuid.Parse(customerID)
	if err != nil {
		logrus.Error(err.Error())
	}

	dataResponse, httpCode, errInfo = c.useCase.ChangePassword(ctx, custUUID, &dtoRequest)
	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	response.SendBack(ctx, dataResponse, errInfo, httpCode)
	return
}

func (c *AccountController) ValidateRefCode(ctx *gin.Context) {
	var (
		dtoRequest  dtos.AccountRefCodeValidationRequest
		dtoResponse dtos.AccountRefCodeValidationResponse
		errInfo     []errorsinfo.Errors
		httpCode    int
	)

	if err := ctx.ShouldBindJSON(&dtoRequest); err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "body payload required")
		response.SendBack(ctx, dtos.AccountRefCodeValidationResponse{}, errInfo, http.StatusBadRequest)
		return
	}

	dtoResponse, httpCode, errInfo = c.useCase.ValidateRefCode(&dtoRequest)
	response.SendBack(ctx, dtoResponse, errInfo, httpCode)
	return
}

func (c *AccountController) SetAvatar(ctx *gin.Context) {
	var (
		dtoRequest  dtos.AccountAvatarRequest
		dtoResponse dtos.AccountAvatarResponse
		errInfo     []errorsinfo.Errors
		httpCode    int
	)

	if err := ctx.ShouldBindJSON(&dtoRequest); err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "body payload required")
		response.SendBack(ctx, dtos.AccountAvatarResponse{}, errInfo, http.StatusBadRequest)
		return
	}

	if dtoRequest.ImageBase64 == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "image_base64 attribute can not be empty")
		response.SendBack(ctx, dtos.AccountAvatarResponse{}, errInfo, http.StatusBadRequest)
		return
	}

	dtoResponse, httpCode, errInfo = c.useCase.SetAvatar(ctx, &dtoRequest)
	response.SendBack(ctx, dtoResponse, errInfo, httpCode)
	return

}

func (c *AccountController) RemoveAvatar(ctx *gin.Context) {
	var (
		dtoResponse dtos.AccountAvatarResponse
		errInfo     []errorsinfo.Errors
		httpCode    int
	)

	customerID := ctx.Param("customer-id")
	if customerID == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "customer ID required in URL parameter")
		response.SendBack(ctx, dtos.AccountAvatarResponse{}, errInfo, http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(customerID)
	if err != nil {
		logrus.Error(err.Error())
	}

	c.useCase.RemoveAvatar(ctx, id)

	response.SendBack(ctx, dtoResponse, errInfo, httpCode)
	return
}

func (c *AccountController) Sharing(ctx *gin.Context) {
	var (
		dtoRequest dtos.AccountGroupSharing
		errInfo    []errorsinfo.Errors
		httpCode   int
	)

	// binding
	if err := ctx.ShouldBindJSON(&dtoRequest); err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "body payload required")
		response.SendBack(ctx, dtos.AccountAvatarResponse{}, errInfo, http.StatusBadRequest)
		return
	}

	response.SendBack(ctx, nil, errInfo, httpCode)
	return
}