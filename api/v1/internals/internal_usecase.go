package internals

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/semicolon-indonesia/wealthy-backend/utils/errorsinfo"
	"github.com/sirupsen/logrus"
	"net/http"
)

type (
	InternalUseCase struct {
		repo IInternalRepository
	}

	IInternalUseCase interface {
		TransactionNotes(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
	}
)

func NewInternalUseCase(repo IInternalRepository) *InternalUseCase {
	return &InternalUseCase{repo: repo}
}

func (s *InternalUseCase) TransactionNotes(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	month := ctx.Query("month")
	year := ctx.Query("year")
	customerID := ctx.Query("customerid")

	if month == "" || year == "" || customerID == "" {
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "need month, year and customer ID in url params")
		return response, http.StatusBadRequest, errInfo
	}

	customerIDUUID, err := uuid.Parse(customerID)
	if err != nil {
		logrus.Error(err.Error())
	}

	response = s.repo.ByNote(customerIDUUID, month, year)
	return response, http.StatusOK, []errorsinfo.Errors{}
}