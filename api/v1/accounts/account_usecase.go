package accounts

import (
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/accounts/dtos"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/accounts/entities"
	"github.com/semicolon-indonesia/wealthy-backend/utils/errorsinfo"
	"github.com/semicolon-indonesia/wealthy-backend/utils/password"
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
