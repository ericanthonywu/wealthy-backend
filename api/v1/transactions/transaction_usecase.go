package transactions

import (
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/transactions/dtos"
	"github.com/semicolon-indonesia/wealthy-backend/utils/errorsinfo"
)

type (
	TransactionUseCase struct {
		repo ITransactionRepository
	}

	ITransactionUseCase interface {
		Expense(request *dtos.TransactionExpenseRequest) (data interface{}, httpCode int, errInfo []errorsinfo.Errors)
	}
)

func NewTransactionUseCase(repo ITransactionRepository) *TransactionUseCase {
	return &TransactionUseCase{repo: repo}
}

func (s *TransactionUseCase) Expense(request *dtos.TransactionExpenseRequest) (data interface{}, httpCode int, errInfo []errorsinfo.Errors) {

	return
}
