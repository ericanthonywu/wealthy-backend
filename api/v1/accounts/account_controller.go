package accounts

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/wealthy-app/wealthy-backend/api/v1/accounts/dtos"
	"github.com/wealthy-app/wealthy-backend/constants"
	"github.com/wealthy-app/wealthy-backend/utils/errorsinfo"
	"github.com/wealthy-app/wealthy-backend/utils/response"
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
		ForgotPassword(ctx *gin.Context)
		ValidateRefCode(ctx *gin.Context)
		SetAvatar(ctx *gin.Context)
		RemoveAvatar(ctx *gin.Context)
		SearchAccount(ctx *gin.Context)
		InviteSharing(ctx *gin.Context)
		AcceptSharing(ctx *gin.Context)
		RejectSharing(ctx *gin.Context)
		RemoveSharing(ctx *gin.Context)
		VerifyOTP(ctx *gin.Context)
		ChangePasswordForgot(ctx *gin.Context)
		GroupSharingAccepted(ctx *gin.Context)
		GroupSharingPending(ctx *gin.Context)
		DeleteAccount(ctx *gin.Context)
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
		dtoRequest dtos.AccountSignInRequest
		errInfo    []errorsinfo.Errors
	)

	// bind
	if err := ctx.ShouldBindJSON(&dtoRequest); err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "no body payload")
		response.SendBack(ctx, struct{}{}, errInfo, http.StatusBadRequest)
		return
	}

	// validate
	if dtoRequest.Email == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "email account empty value")
	}

	if dtoRequest.Password == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "password empty value")
	}

	if len(errInfo) > 0 {
		response.SendBack(ctx, struct{}{}, errInfo, http.StatusBadRequest)
		return
	}

	dtoResponse, httpCode, errInfo := c.useCase.SignIn(&dtoRequest)
	response.SendBack(ctx, dtoResponse, errInfo, httpCode)
	return
}

func (c *AccountController) SignOut(ctx *gin.Context) {
	c.useCase.SignOut()
}

func (c *AccountController) GetProfile(ctx *gin.Context) {
	var errInfo []errorsinfo.Errors

	// check method
	if ctx.Request.Method != http.MethodGet {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "invalid method")
		response.SendBack(ctx, struct{}{}, errInfo, http.StatusMethodNotAllowed)
		return
	}

	data, httpCode, errInfo := c.useCase.GetProfile(ctx)
	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *AccountController) UpdateProfile(ctx *gin.Context) {
	var (
		dtoRequest map[string]interface{}
		errInfo    []errorsinfo.Errors
	)

	// bind
	if err := ctx.ShouldBindJSON(&dtoRequest); err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "body payload required")
		response.SendBack(ctx, dtos.AccountSetProfileRequest{}, errInfo, http.StatusBadRequest)
		return
	}

	dtoResponse, httpCode, errInfo := c.useCase.UpdateProfile(ctx, dtoRequest)
	response.SendBack(ctx, dtoResponse, errInfo, httpCode)
	return
}

func (c *AccountController) ChangePassword(ctx *gin.Context) {
	var (
		dtoRequest dtos.AccountChangePassword
		errInfo    []errorsinfo.Errors
		httpCode   int
	)

	// bind
	if err := ctx.ShouldBindJSON(&dtoRequest); err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "body payload required")
		response.SendBack(ctx, struct{}{}, errInfo, http.StatusBadRequest)
		return
	}

	// validation
	if dtoRequest.NewPassword == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "new password attribute empty value")
	}

	if dtoRequest.OldPassword == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "old password attribute empty value")
	}

	if len(errInfo) > 0 {
		response.SendBack(ctx, struct{}{}, errInfo, http.StatusBadRequest)
		return
	}

	dataResponse, httpCode, errInfo := c.useCase.ChangePassword(ctx, &dtoRequest)
	response.SendBack(ctx, dataResponse, errInfo, httpCode)
	return
}

func (c *AccountController) ForgotPassword(ctx *gin.Context) {
	var (
		dtoRequest dtos.AccountForgotPasswordRequest
		errInfo    []errorsinfo.Errors
	)

	// bind
	if err := ctx.ShouldBindJSON(&dtoRequest); err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "body payload required")
		response.SendBack(ctx, dtos.AccountRefCodeValidationResponse{}, errInfo, http.StatusBadRequest)
		return
	}

	// validation
	if dtoRequest.EmailAccount == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "email account empty value")
		response.SendBack(ctx, dtos.AccountRefCodeValidationResponse{}, errInfo, http.StatusBadRequest)
		return
	}

	dtoResponse, httpCode, errInfo := c.useCase.ForgotPassword(ctx, &dtoRequest)
	response.SendBack(ctx, dtoResponse, errInfo, httpCode)
	return
}

func (c *AccountController) ValidateRefCode(ctx *gin.Context) {
	var (
		dtoRequest dtos.AccountRefCodeValidationRequest
		errInfo    []errorsinfo.Errors
	)

	// bind
	if err := ctx.ShouldBindJSON(&dtoRequest); err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "body payload required")
		response.SendBack(ctx, struct{}{}, errInfo, http.StatusBadRequest)
		return
	}

	// validate
	if dtoRequest.RefCode == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "referral code empty value")
	}

	// err empty
	if len(errInfo) > 0 {
		response.SendBack(ctx, struct{}{}, errInfo, http.StatusBadRequest)
		return
	}

	dtoResponse, httpCode, errInfo := c.useCase.ValidateRefCode(&dtoRequest)
	response.SendBack(ctx, dtoResponse, errInfo, httpCode)
	return
}

func (c *AccountController) SetAvatar(ctx *gin.Context) {
	var (
		dtoRequest dtos.AccountAvatarRequest
		errInfo    []errorsinfo.Errors
	)

	// check method
	if ctx.Request.Method != http.MethodPost {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "invalid method")
		response.SendBack(ctx, struct{}{}, errInfo, http.StatusMethodNotAllowed)
		return
	}

	// bind
	if err := ctx.ShouldBindJSON(&dtoRequest); err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "body payload required")
		response.SendBack(ctx, struct{}{}, errInfo, http.StatusBadRequest)
		return
	}

	// validate
	if dtoRequest.ImageBase64 == "" {
		resp := struct {
			Message string `json:"message,omitempty"`
		}{
			Message: "image base64 empty value",
		}
		response.SendBack(ctx, resp, errInfo, http.StatusBadRequest)
		return
	}

	dtoResponse, httpCode, errInfo := c.useCase.SetAvatar(ctx, &dtoRequest)
	response.SendBack(ctx, dtoResponse, errInfo, httpCode)
	return

}

func (c *AccountController) RemoveAvatar(ctx *gin.Context) {
	var (
		errInfo  []errorsinfo.Errors
		httpCode int
	)

	// check method
	if ctx.Request.Method != http.MethodDelete {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "invalid method")
		response.SendBack(ctx, struct{}{}, errInfo, http.StatusMethodNotAllowed)
		return
	}

	dataResponse, httpCode, errInfo := c.useCase.RemoveAvatar(ctx)
	response.SendBack(ctx, dataResponse, errInfo, httpCode)
	return
}

func (c *AccountController) SearchAccount(ctx *gin.Context) {
	var (
		dtoRequest dtos.AccountGroupSharing
		errInfo    []errorsinfo.Errors
		httpCode   int
	)

	// get account type
	accountType := ctx.MustGet("accountType").(string)

	// if basic account
	if accountType == constants.AccountBasic {
		resp := struct {
			Message string `json:"message"`
		}{
			Message: constants.ProPlan,
		}
		response.SendBack(ctx, resp, []errorsinfo.Errors{}, http.StatusUpgradeRequired)
		return
	}

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

	// get account type
	accountType := ctx.MustGet("accountType").(string)

	// if basic account
	if accountType == constants.AccountBasic {
		resp := struct {
			Message string `json:"message"`
		}{
			Message: constants.ProPlan,
		}
		response.SendBack(ctx, resp, []errorsinfo.Errors{}, http.StatusUpgradeRequired)
		return
	}

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
	if dtoRequest.IDGroupSharing == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id group sharing empty value")
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
	if dtoRequest.IDGroupSharing == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "id group sharing empty value")
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
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "email account empty value")
		response.SendBack(ctx, struct{}{}, errInfo, http.StatusBadRequest)
		return
	}

	data, httpCode, errInfo := c.useCase.RemoveSharing(ctx, &dtoRequest)
	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *AccountController) VerifyOTP(ctx *gin.Context) {
	var (
		dtoRequest dtos.AccountOTPVerify
		errInfo    []errorsinfo.Errors
	)

	// binding
	if err := ctx.ShouldBindJSON(&dtoRequest); err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "body payload required")
		response.SendBack(ctx, struct{}{}, errInfo, http.StatusBadRequest)
		return
	}

	// validate
	if dtoRequest.EmailAccount == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "email account empty value")
	}

	if dtoRequest.OTPCode == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "otp code empty value")
	}

	if len(errInfo) > 0 {
		response.SendBack(ctx, struct{}{}, errInfo, http.StatusBadRequest)
		return
	}

	data, httpCode, errInfo := c.useCase.VerifyOTP(ctx, &dtoRequest)
	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *AccountController) ChangePasswordForgot(ctx *gin.Context) {
	var (
		dtoRequest dtos.AccountChangeForgotPassword
		errInfo    []errorsinfo.Errors
	)

	// binding
	if err := ctx.ShouldBindJSON(&dtoRequest); err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "body payload required")
		response.SendBack(ctx, struct{}{}, errInfo, http.StatusBadRequest)
		return
	}

	// validate
	if dtoRequest.NewPassword == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "new password required")
	}

	// if any error
	if len(errInfo) > 0 {
		response.SendBack(ctx, struct{}{}, errInfo, http.StatusBadRequest)
		return
	}

	data, httpCode, errInfo := c.useCase.ChangePasswordForgot(ctx, &dtoRequest)
	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *AccountController) GroupSharingAccepted(ctx *gin.Context) {
	// get account type
	accountType := ctx.MustGet("accountType").(string)

	// if basic account
	if accountType == constants.AccountBasic {
		resp := struct {
			Message string `json:"message"`
		}{
			Message: constants.ProPlan,
		}
		response.SendBack(ctx, resp, []errorsinfo.Errors{}, http.StatusUpgradeRequired)
		return
	}

	data, httpCode, errInfo := c.useCase.GroupSharingAccepted(ctx)
	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *AccountController) GroupSharingPending(ctx *gin.Context) {
	// get account type
	accountType := ctx.MustGet("accountType").(string)

	// if basic account
	if accountType == constants.AccountBasic {
		resp := struct {
			Message string `json:"message"`
		}{
			Message: constants.ProPlan,
		}
		response.SendBack(ctx, resp, []errorsinfo.Errors{}, http.StatusUpgradeRequired)
		return
	}

	data, httpCode, errInfo := c.useCase.GroupSharingPending(ctx)
	response.SendBack(ctx, data, errInfo, httpCode)
	return
}

func (c *AccountController) DeleteAccount(ctx *gin.Context) {
	data, httpCode, errInfo := c.useCase.DeleteAccount(ctx)
	response.SendBack(ctx, data, errInfo, httpCode)
	return
}