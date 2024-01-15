package transactions

import (
	"github.com/gin-gonic/gin"
	"github.com/wealthy-app/wealthy-backend/api/v1/transactions/dtos"
)

type (
	TransactionUseCase struct {
		repo ITransactionRepository
	}

	ITransactionUseCase interface {
		InvestmentRecords(ctx *gin.Context, request *dtos.TransactionRequestInvestment) (response interface{}, httpCode int, errInfo []string)
	}
)

func NewTransactionUseCase(repo ITransactionRepository) *TransactionUseCase {
	return &TransactionUseCase{repo: repo}
}

func (s *TransactionUseCase) InvestmentRecords(ctx *gin.Context, request *dtos.TransactionRequestInvestment) (response interface{}, httpCode int, errInfo []string) {
	return
}
