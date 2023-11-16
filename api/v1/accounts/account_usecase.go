package accounts

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/accounts/dtos"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/accounts/entities"
	"github.com/semicolon-indonesia/wealthy-backend/utils/errorsinfo"
	"github.com/semicolon-indonesia/wealthy-backend/utils/password"
	"github.com/semicolon-indonesia/wealthy-backend/utils/personalaccounts"
	"github.com/semicolon-indonesia/wealthy-backend/utils/token"
	"github.com/semicolon-indonesia/wealthy-backend/utils/utilities"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

type (
	AccountUseCase struct {
		repo IAccountRepository
	}

	IAccountUseCase interface {
		SignIn(request *dtos.AccountSignInRequest) (response dtos.AccountSignInResponse, httpCode int, errInfo []errorsinfo.Errors)
		SignUp(request *dtos.AccountSignUpRequest) (response dtos.AccountSignUpResponse, httpCode int, errInfo []errorsinfo.Errors)
		SignOut()
		GetProfile(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		UpdateProfile(ctx *gin.Context, customerID uuid.UUID, request map[string]interface{}) (response map[string]bool, httpCode int, errInfo []errorsinfo.Errors)
		ChangePassword(ctx *gin.Context, customerID uuid.UUID, request *dtos.AccountChangePassword) (response map[string]bool, httpCode int, errInfo []errorsinfo.Errors)
		ValidateRefCode(request *dtos.AccountRefCodeValidationRequest) (response dtos.AccountRefCodeValidationResponse, httpCode int, errInfo []errorsinfo.Errors)
		SetAvatar(ctx *gin.Context, request *dtos.AccountAvatarRequest) (response dtos.AccountAvatarResponse, httpCode int, errInfo []errorsinfo.Errors)
		RemoveAvatar(ctx *gin.Context, customerID uuid.UUID) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
	}
)

func NewAccountUseCase(repo IAccountRepository) *AccountUseCase {
	return &AccountUseCase{repo: repo}
}

func (s *AccountUseCase) SignUp(request *dtos.AccountSignUpRequest) (response dtos.AccountSignUpResponse, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		role              string
		idPersonalAccount uuid.UUID
		level             int
		err               error
	)

	personalAccountEntity := entities.AccountSignUpPersonalAccountEntity{
		Username:  request.Username,
		Name:      request.Name,
		Email:     request.Email,
		ReferCode: request.RefCode,
	}

	authAccountActivity := entities.AccountSignUpAuthenticationsEntity{
		Password: password.Generate(request.Password),
		Active:   true,
	}

	role, idPersonalAccount, httpCode, errInfo = s.repo.SignUp(&personalAccountEntity, &authAccountActivity)

	if len(errInfo) > 0 {
		return response, httpCode, errInfo
	}

	if request.RefCodeReference != "" {
		// Check reference first
		listRefCodeExist := s.repo.ListRefCode()
		if !utilities.Contains(listRefCodeExist, request.RefCodeReference) {
			httpCode = http.StatusBadRequest
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", errors.New("referral code reference not exist").Error())
			return response, httpCode, errInfo
		}

		if request.RefCodeReference != "" {
			level, err = s.repo.GetLevelReferenceCode(request.RefCodeReference)
			if err != nil {
				logrus.Error(err.Error())
			}

			if err == nil {
				level = level + 1
			}

			if level == 0 {
				level = level + 1
			}
		}

		if request.RefCodeReference == "" {
			level = 0

		}
	}

	newID, err := uuid.NewUUID()
	if err != nil {
		logrus.Error(err.Error())
	}

	// WRITING FOR REFERENCE CODE
	model := entities.AccountRewards{
		ID:               newID,
		RefCode:          request.RefCode,
		RefCodeReference: request.RefCodeReference,
		Level:            level,
	}

	err = s.repo.WriteRewardsList(&model)
	if err != nil {
		logrus.Error(err.Error())
	}

	response.Customer = dtos.CustomerDetail{
		ID:       idPersonalAccount,
		Name:     request.Name,
		Username: request.Username,
		Email:    request.Email,
	}

	response.Account = dtos.Account{
		Role: role,
	}

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	return response, httpCode, errInfo
}

func (s *AccountUseCase) SignIn(request *dtos.AccountSignInRequest) (response dtos.AccountSignInResponse, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		err error
	)

	authentication := entities.AccountSignInAuthenticationEntity{
		Email:    request.Email,
		Password: request.Password,
	}

	data := s.repo.SignInAuth(authentication)
	resultOfCompare := password.Compare(data.Password, []byte(request.Password))

	if !resultOfCompare {
		response.Customer.CustomerID = ""
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "email or password doesn't match")
		return dtos.AccountSignInResponse{}, http.StatusBadRequest, errInfo
	}

	response.Token, response.TokenExp, err = token.JWTBuilder(request.Email, data.Roles)
	if err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "can not give proper token")
		return dtos.AccountSignInResponse{}, http.StatusUnprocessableEntity, errInfo
	}

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	response.Customer.Email = request.Email
	response.Customer.CustomerID = data.ID.String()
	response.Account.Role = data.Roles

	return response, httpCode, errInfo
}

func (s *AccountUseCase) SignOut() {

}

func (s *AccountUseCase) GetProfile(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var dtoResponse dtos.AccountProfile

	usrEmail := ctx.MustGet("email").(string)
	personalAccount := personalaccounts.Informations(ctx, usrEmail)

	if personalAccount.ID == uuid.Nil {
		httpCode = http.StatusUnauthorized
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "token contains invalid information")
		return response, httpCode, errInfo
	}

	dataProfile := s.repo.GetProfile(personalAccount.ID)

	dtoResponse.AccountCustomer.ID = dataProfile.ID
	dtoResponse.AccountCustomer.Name = dataProfile.Name
	dtoResponse.AccountCustomer.ReferType = dataProfile.ReferType
	dtoResponse.AccountCustomer.Username = dataProfile.Username
	dtoResponse.AccountCustomer.Email = dataProfile.Email

	dtoResponse.AccountCustomer.Gender.ID = dataProfile.IDGender
	dtoResponse.AccountCustomer.Gender.Value = dataProfile.Gender

	dtoResponse.AccountDetail.AccountType = dataProfile.AccountType
	dtoResponse.AccountDetail.UserRoles = dataProfile.UserRoles

	dtoResponse.AccountAvatar.URL = os.Getenv("APP_HOST") + "/v1/" + dataProfile.ImagePath
	dtoResponse.AccountAvatar.FileName = dataProfile.FileName

	dtoResponse.AccountCustomer.DOB = dataProfile.DOB

	return dtoResponse, http.StatusOK, []errorsinfo.Errors{}
}

func (s *AccountUseCase) UpdateProfile(ctx *gin.Context, customerID uuid.UUID, request map[string]interface{}) (response map[string]bool, httpCode int, errInfo []errorsinfo.Errors) {
	var err error
	dtoResponse := make(map[string]bool)
	dtoResponse["success"] = false

	if request["id"] != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", errors.New("can not change customer ID").Error())
	}

	if request["refer_code"] != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", errors.New("can not change refer code").Error())
	}

	if request["email"] != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", errors.New("can not change email").Error())
	}

	if request["image_path"] != nil || request["filename"] != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", errors.New("can not change avatar. use set avatar API instead").Error())
	}

	if request["id_master_account_type"] != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", errors.New("can not change account type. use subscription API for switch to PRO account").Error())
	}

	if len(errInfo) > 0 {
		return dtoResponse, http.StatusBadRequest, errInfo
	}

	err = s.repo.UpdateProfile(customerID, request)
	if err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return dtoResponse, http.StatusInternalServerError, errInfo
	}

	dtoResponse["success"] = true
	return dtoResponse, http.StatusOK, []errorsinfo.Errors{}
}

func (s *AccountUseCase) ChangePassword(ctx *gin.Context, customerID uuid.UUID, request *dtos.AccountChangePassword) (response map[string]bool, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		update      = make(map[string]interface{})
		dtoResponse = make(map[string]bool)
	)

	dtoResponse["success"] = false
	dataProfile := s.repo.GetProfilePassword(customerID)

	resultOfCompare := password.Compare(dataProfile.Password, []byte(request.OldPassword))
	if !resultOfCompare {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "old password incorrect")
		return dtoResponse, http.StatusUnprocessableEntity, errInfo
	}

	newPassword := password.Generate(request.NewPassword)
	update["password"] = newPassword

	err := s.repo.UpdatePassword(customerID, update)
	if err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return dtoResponse, http.StatusInternalServerError, errInfo
	}

	dtoResponse["success"] = true
	return dtoResponse, http.StatusOK, errInfo
}

func (s *AccountUseCase) ValidateRefCode(request *dtos.AccountRefCodeValidationRequest) (response dtos.AccountRefCodeValidationResponse, httpCode int, errInfo []errorsinfo.Errors) {
	var dtoResponse dtos.AccountRefCodeValidationResponse
	RefCodeList := s.repo.ListRefCode()

	if utilities.Contains(RefCodeList, request.RefCode) {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", errors.New("referral code already exist on system. please use another code to be registered").Error())
		return dtos.AccountRefCodeValidationResponse{}, http.StatusOK, errInfo
	}

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	dtoResponse.Available = true
	return dtoResponse, http.StatusOK, errInfo
}

func (s *AccountUseCase) SetAvatar(ctx *gin.Context, request *dtos.AccountAvatarRequest) (response dtos.AccountAvatarResponse, httpCode int, errInfo []errorsinfo.Errors) {
	var dtoResponse dtos.AccountAvatarResponse
	updateProfile := make(map[string]interface{})

	imageData, err := base64.StdEncoding.DecodeString(request.ImageBase64)
	if err != nil {
		httpCode = http.StatusBadRequest
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return dtos.AccountAvatarResponse{}, httpCode, errInfo
		return
	}

	usrEmail := ctx.MustGet("email").(string)
	personalAccount := personalaccounts.Informations(ctx, usrEmail)

	if personalAccount.ID == uuid.Nil {
		httpCode = http.StatusUnauthorized
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "token contains invalid information")
		return response, httpCode, errInfo
	}

	dataProfile := s.repo.GetProfile(personalAccount.ID)

	if dataProfile.ImagePath != "" || dataProfile.FileName != "" {
		target := "assets/avatar/" + dataProfile.FileName
		err = os.Remove(target)
		if err != nil {
			logrus.Error(err.Error())
		}
	}

	updateProfile["image_path"] = ""
	updateProfile["file_name"] = ""
	err = s.repo.UpdateProfile(personalAccount.ID, updateProfile)
	if err != nil {
		logrus.Error(err.Error())
	}

	filename := fmt.Sprintf("%d", time.Now().Unix()) + ".png"
	targetPath := "assets/avatar/" + filename

	err = utilities.SaveImage(imageData, targetPath)
	if err != nil {
		logrus.Error(err.Error())
	}

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	updateProfile["image_path"] = "images/avatar/" + filename
	updateProfile["file_name"] = filename
	err = s.repo.UpdateProfile(personalAccount.ID, updateProfile)
	if err != nil {
		logrus.Error(err.Error())
	}

	dtoResponse.Success = true
	return dtoResponse, http.StatusOK, errInfo
}

func (s *AccountUseCase) RemoveAvatar(ctx *gin.Context, customerID uuid.UUID) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		err         error
		dtoResponse dtos.AccountAvatarResponse
	)

	updateProfile := make(map[string]interface{})
	dataProfile := s.repo.GetProfile(customerID)

	if dataProfile.ImagePath != "" || dataProfile.FileName != "" {
		target := "assets/avatar/" + dataProfile.FileName
		err = os.Remove(target)
		if err != nil {
			logrus.Error(err.Error())
			dtoResponse.Success = false
			errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
			return dtoResponse, http.StatusInternalServerError, errInfo
		}
	}

	updateProfile["image_path"] = ""
	updateProfile["file_name"] = ""
	err = s.repo.UpdateProfile(customerID, updateProfile)
	if err != nil {
		logrus.Error(err.Error())
		dtoResponse.Success = false
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return dtoResponse, http.StatusInternalServerError, errInfo
	}

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	dtoResponse.Success = true
	return dtoResponse, http.StatusOK, errInfo
}