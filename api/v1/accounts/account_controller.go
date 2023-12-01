package accounts

import (
	"errors"
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
		SearchAccount(ctx *gin.Context)
		InviteSharing(ctx *gin.Context)
		AcceptSharing(ctx *gin.Context)
		RejectSharing(ctx *gin.Context)
		RemoveSharing(ctx *gin.Context)
	}
)

func NewAccountController(useCase IAccountUseCase) *AccountController {
	return &AccountController{useCase: useCase}
}

func (c *AccountController) SignUp(ctx *gin.Context) {
	var (
		dtoRequest dtos.AccountSignUpRequest
		httpCode   int
		errInfo    []errorsinfo.Errors
	)

	// binding
	if err := ctx.ShouldBindJSON(&dtoRequest); err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		response.SendBack(ctx, dtos.AccountSignUpResponse{}, errInfo, http.StatusBadRequest)
		return
	}

	// validate
	if dtoRequest.Password == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", errors.New("password empty in payload").Error())
	}

	if dtoRequest.Username == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", errors.New("username empty in payload").Error())
	}

	if dtoRequest.Name == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", errors.New("name empty in payload").Error())
	}

	if dtoRequest.Email == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", errors.New("email empty in payload").Error())
	}

	if dtoRequest.RefCode == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", errors.New("referral code empty in payload").Error())
	}

	if dtoRequest.RefCode == dtoRequest.RefCodeReference {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", errors.New("referral code and referral code reference identically value in payload").Error())
	}

	if len(errInfo) > 0 {
		resp := struct {
			Message string `json:"message,omitempty"`
		}{}
		response.SendBack(ctx, resp, errInfo, http.StatusBadRequest)
		return
	}

	dataResponse, httpCode, errInfo := c.useCase.SignUp(&dtoRequest)
	response.SendBack(ctx, dataResponse, errInfo, httpCode)
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

func (c *AccountController) SearchAccount(ctx *gin.Context) {
	var (
		dtoRequest dtos.AccountGroupSharing
		errInfo    []errorsinfo.Errors
		httpCode   int
	)

	// binding
	if err := ctx.ShouldBindJSON(&dtoRequest); err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "body payload required")
		response.SendBack(ctx, struct{}{}, errInfo, http.StatusBadRequest)
		return
	}

	// validate
	if dtoRequest.EmailAccount == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "email account value empty")
		response.SendBack(ctx, struct{}{}, errInfo, http.StatusBadRequest)
		return
	}

	data, httpCode, errInfo := c.useCase.SearchAccount(ctx, &dtoRequest)
	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *AccountController) InviteSharing(ctx *gin.Context) {
	var (
		dtoRequest dtos.AccountGroupSharing
		errInfo    []errorsinfo.Errors
		httpCode   int
	)

	// binding
	if err := ctx.ShouldBindJSON(&dtoRequest); err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "body payload required")
		response.SendBack(ctx, struct{}{}, errInfo, http.StatusBadRequest)
		return
	}

	// validate
	if dtoRequest.EmailAccount == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "email account value empty")
		response.SendBack(ctx, struct{}{}, errInfo, http.StatusBadRequest)
		return
	}

	data, httpCode, errInfo := c.useCase.InviteSharing(ctx, &dtoRequest)
	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *AccountController) AcceptSharing(ctx *gin.Context) {
	var (
		dtoRequest dtos.AccountGroupSharingAccept
		errInfo    []errorsinfo.Errors
		httpCode   int
	)

	// binding
	if err := ctx.ShouldBindJSON(&dtoRequest); err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "body payload required")
		response.SendBack(ctx, struct{}{}, errInfo, http.StatusBadRequest)
		return
	}

	// validate process
	if dtoRequest.IDRecipient == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "ID Receipt can not be empty")
	}

	if dtoRequest.IDSender == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "ID Sender can not be empty")
	}

	if dtoRequest.IDSender == dtoRequest.IDRecipient {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "ID Sender and ID Receipt cannot be identically")
	}

	if len(errInfo) > 0 {
		response.SendBack(ctx, struct{}{}, errInfo, http.StatusBadRequest)
		return
	}

	data, httpCode, errInfo := c.useCase.AcceptSharing(ctx, &dtoRequest)
	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *AccountController) RejectSharing(ctx *gin.Context) {
	var (
		dtoRequest dtos.AccountGroupSharingAccept
		errInfo    []errorsinfo.Errors
		httpCode   int
	)

	// binding
	if err := ctx.ShouldBindJSON(&dtoRequest); err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "body payload required")
		response.SendBack(ctx, struct{}{}, errInfo, http.StatusBadRequest)
		return
	}

	// validate process
	if dtoRequest.IDRecipient == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "ID Receipt can not be empty")
	}

	if dtoRequest.IDSender == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "ID Sender can not be empty")
	}

	if dtoRequest.IDSender == dtoRequest.IDRecipient {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "ID Sender and ID Receipt cannot be identically")
	}

	if len(errInfo) > 0 {
		response.SendBack(ctx, struct{}{}, errInfo, http.StatusBadRequest)
		return
	}

	data, httpCode, errInfo := c.useCase.RejectSharing(ctx, &dtoRequest)
	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *AccountController) RemoveSharing(ctx *gin.Context) {
	var (
		dtoRequest dtos.AccountGroupSharingRemove
		errInfo    []errorsinfo.Errors
	)

	// binding
	if err := ctx.ShouldBindJSON(&dtoRequest); err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "body payload required")
		response.SendBack(ctx, struct{}{}, errInfo, http.StatusBadRequest)
		return
	}

	if dtoRequest.EmailAccount == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "body payload required")
		response.SendBack(ctx, struct{}{}, errInfo, http.StatusBadRequest)
		return
	}

	data, httpCode, errInfo := c.useCase.RemoveSharing(ctx, &dtoRequest)
	response.SendBack(ctx, data, errInfo, httpCode)
	return
}