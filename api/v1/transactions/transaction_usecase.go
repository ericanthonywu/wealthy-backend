package transactions

import (
	"github.com/google/uuid"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/transactions/dtos"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/transactions/entities"
	"github.com/semicolon-indonesia/wealthy-backend/utils/errorsinfo"
	"net/http"
)

type (
	TransactionUseCase struct {
		repo ITransactionRepository
	}

	ITransactionUseCase interface {
		Add(request *dtos.TransactionRequest) (data interface{}, httpCode int, errInfo []errorsinfo.Errors)
	}
)

func NewTransactionUseCase(repo ITransactionRepository) *TransactionUseCase {
	return &TransactionUseCase{repo: repo}
}

func (s *TransactionUseCase) Add(request *dtos.TransactionRequest) (data interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		trxID uuid.UUID
		err   error
	)

	trxID, err = uuid.NewUUID()
	if err != nil {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return nil, http.StatusUnprocessableEntity, errInfo
	}

	modelTransaction := entities.TransactionEntity{
		ID:                            trxID,
		Date:                          request.Date,
		Fees:                          float64(request.Fees),
		Amount:                        float64(request.Amount),
		IDWallet:                      request.IDWallet,
		IDMasterIncomeCategories:      request.IDMasterIncomeCategories,
		IDMasterExpenseCategories:     request.IDMasterExpenseCategories,
		IDMasterInvest:                request.IDMasterInvest,
		IDMasterBroker:                request.IDMasterBroker,
		IDMasterReksanadaTypes:        request.IDMasterReksanadaTypes,
		IDMasterTransactionPriorities: request.IDMasterTransactionPriorities,
		IDMasterTransactionTypes:      request.IDMasterTransactionTypes,
	}

	modelTransactionDetail := entities.TransactionDetailEntity{
		IDTransaction:     trxID,
		Repeat:            request.Repeat,
		Note:              request.Note,
		From:              request.From,
		To:                request.To,
		MutualFundProduct: request.MutualFundProduct,
		StockCode:         request.StockCode,
		Lot:               request.Lot,
	}

	err = s.repo.Add(&modelTransaction, &modelTransactionDetail)
	if err != nil {
		data := struct {
			IDTransaction uuid.UUID
		}{
			IDTransaction: uuid.Nil,
		}
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return data, http.StatusInternalServerError, errInfo
	}

	data = struct {
		IDTransaction uuid.UUID
	}{
		IDTransaction: trxID,
	}
	return data, http.StatusOK, []errorsinfo.Errors{}
}
