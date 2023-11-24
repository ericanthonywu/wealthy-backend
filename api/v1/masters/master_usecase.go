package masters

import (
	"github.com/google/uuid"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/masters/entities"
	"github.com/semicolon-indonesia/wealthy-backend/utils/errorsinfo"
	"net/http"
)

type (
	MasterUseCase struct {
		repo IMasterRepository
	}

	IMasterUseCase interface {
		TransactionType() (data interface{})
		IncomeType() (data interface{})
		ExpenseType() (data interface{})
		ReksadanaType() (data interface{})
		WalletType() (data interface{})
		InvestType() (data interface{})
		Broker() (data interface{})
		TransactionPriority() (data interface{})
		Gender() (data interface{})
		SubExpenseCategories(expenseID uuid.UUID) (data interface{})
		Exchange() (data interface{}, httpCode int, errInfo []errorsinfo.Errors)
	}
)

func NewMasterUseCase(repo IMasterRepository) *MasterUseCase {
	return &MasterUseCase{repo: repo}
}

func (s *MasterUseCase) TransactionType() (data interface{}) {
	return s.repo.TransactionType()
}

func (s *MasterUseCase) IncomeType() (data interface{}) {
	return s.repo.IncomeType()
}

func (s *MasterUseCase) ExpenseType() (data interface{}) {
	return s.repo.ExpenseType()
}

func (s *MasterUseCase) ReksadanaType() (data interface{}) {
	return s.repo.ReksadanaType()
}

func (s *MasterUseCase) WalletType() (data interface{}) {
	return s.repo.WalletType()
}

func (s *MasterUseCase) InvestType() (data interface{}) {
	return s.repo.InvestType()
}

func (s *MasterUseCase) Broker() (data interface{}) {
	return s.repo.Broker()
}

func (s *MasterUseCase) TransactionPriority() (data interface{}) {
	return s.repo.TransactionPriority()
}

func (s *MasterUseCase) Gender() (data interface{}) {
	return s.repo.Gender()
}

func (s *MasterUseCase) SubExpenseCategories(expenseID uuid.UUID) (data interface{}) {
	if s.repo.ExpenseIDExist(expenseID) {
		return s.repo.SubExpenseCategory(expenseID)
	}
	return data
}

func (s *MasterUseCase) Exchange() (data interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	dataExchange, err := s.repo.Exchange()
	if err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return []entities.Exchange{}, http.StatusInternalServerError, errInfo
	}

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	if len(dataExchange) == 0 {
		message := struct {
			Message string `json:"message"`
		}{
			Message: "no data currency exchange",
		}
		return message, http.StatusNotFound, errInfo
	}

	return dataExchange, http.StatusOK, errInfo
}