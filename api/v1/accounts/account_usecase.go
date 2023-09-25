package accounts

import (
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/accounts/dtos"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/accounts/entities"
	"github.com/semicolon-indonesia/wealthy-backend/utils/errorsinfo"
	"github.com/semicolon-indonesia/wealthy-backend/utils/password"
)

type (
	AccountUseCase struct {
		repo IAccountRepository
	}

	IAccountUseCase interface {
		SignIn()
		SignUp(request *dtos.AccountSignUpRequest) (response dtos.AccountSignUpResponse, httpCode int, errInfo []errorsinfo.Errors)
		SignOut()
	}
)

func NewAccountUseCase(repo IAccountRepository) *AccountUseCase {
	return &AccountUseCase{repo: repo}
}

func (s *AccountUseCase) SignIn() {

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

func (s *AccountUseCase) SignOut() {

}
