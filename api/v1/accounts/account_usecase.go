package accounts

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/accounts/dtos"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/accounts/entities"
	"github.com/semicolon-indonesia/wealthy-backend/utils/errorsinfo"
	"github.com/semicolon-indonesia/wealthy-backend/utils/password"
	"github.com/semicolon-indonesia/wealthy-backend/utils/personalaccounts"
	"github.com/semicolon-indonesia/wealthy-backend/utils/token"
	"net/http"
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
		SetProfile(ctx *gin.Context, request *dtos.AccountSetProfileRequest) (response map[string]bool, httpCode int, errInfo []errorsinfo.Errors)
	}
)

func NewAccountUseCase(repo IAccountRepository) *AccountUseCase {
	return &AccountUseCase{repo: repo}
}

func (s *AccountUseCase) SignUp(request *dtos.AccountSignUpRequest) (response dtos.AccountSignUpResponse, httpCode int, errInfo []errorsinfo.Errors) {
	var role string

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

	role, httpCode, errInfo = s.repo.SignUp(&personalAccountEntity, &authAccountActivity)

	response.Username = request.Username
	response.Name = request.Name
	response.Email = request.Email
	response.Role = role

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	return response, httpCode, errInfo
}

func (s *AccountUseCase) SignIn(request *dtos.AccountSignInRequest) (response dtos.AccountSignInResponse, httpCode int, errInfo []errorsinfo.Errors) {
	var err error

	authentication := entities.AccountSignInAuthenticationEntity{
		Email:    request.Email,
		Password: request.Password,
	}

	data := s.repo.SignInAuth(authentication)
	resultOfCompare := password.Compare(data.Password, []byte(request.Password))

	response.Email = request.Email
	response.Role = data.Roles

	if !resultOfCompare {
		response.Role = ""
		response.Token = ""
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "email or password doesn't match")
		return response, http.StatusUnprocessableEntity, errInfo
	}

	response.Token, err = token.JWTBuilder(response.Email, response.Role)
	if err != nil {
		response.Role = ""
		response.Token = ""
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "can not give proper token")
		return response, http.StatusUnprocessableEntity, errInfo
	}

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	return response, httpCode, errInfo
}

func (s *AccountUseCase) SignOut() {

}

func (s *AccountUseCase) GetProfile(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	usrEmail := ctx.MustGet("email").(string)
	personalAccount := personalaccounts.Informations(ctx, usrEmail)

	if personalAccount.ID == uuid.Nil {
		httpCode = http.StatusNotFound
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "not found")
		return response, httpCode, errInfo
	}

	response = s.repo.GetProfile(personalAccount.ID)
	return response, http.StatusOK, []errorsinfo.Errors{}
}

func (s *AccountUseCase) SetProfile(ctx *gin.Context, request *dtos.AccountSetProfileRequest) (response map[string]bool, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		model      entities.AccountSetProfileEntity
		genderUUID uuid.UUID
	)

	dtoResponse := make(map[string]bool)

	usrEmail := ctx.MustGet("email").(string)
	personalAccount := personalaccounts.Informations(ctx, usrEmail)

	if personalAccount.ID == uuid.Nil {
		httpCode = http.StatusBadRequest
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "cannot access API. invalid personal ID")
		return response, httpCode, errInfo
	}

	_ = copier.Copy(&model, &request)

	if request.Gender != "" {
		genderUUID, _ = uuid.Parse(request.Gender)
		model.Gender = genderUUID
	}

	err := s.repo.SetProfile(personalAccount.ID, &model)
	if err != nil {
		httpCode = http.StatusInternalServerError
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return response, httpCode, errInfo
	}

	dtoResponse["success"] = true
	return dtoResponse, http.StatusOK, []errorsinfo.Errors{}
}
