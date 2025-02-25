package internals

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/wealthy-app/wealthy-backend/utils/errorsinfo"
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
	month := fmt.Sprintf("%02s", ctx.Query("month"))
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
