package statistics

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/semicolon-indonesia/wealthy-backend/utils/errorsinfo"
	"github.com/semicolon-indonesia/wealthy-backend/utils/personalaccounts"
	"net/http"
)

type (
	StatisticUseCase struct {
		repo IStatisticRepository
	}

	IStatisticUseCase interface {
		Statistic(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		TransactionPriority(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		Trend(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		Category(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
	}
)

func NewStatisticUseCase(repo IStatisticRepository) *StatisticUseCase {
	return &StatisticUseCase{repo: repo}
}

func (s *StatisticUseCase) Statistic(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	usrEmail := ctx.MustGet("email").(string)
	personalAccount := personalaccounts.Informations(ctx, usrEmail)

	if personalAccount.ID == uuid.Nil {
		httpCode = http.StatusNotFound
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "not found")
		return response, httpCode, errInfo
	}

	response = s.repo.Statistic(personalAccount.ID)
	return response, http.StatusOK, []errorsinfo.Errors{}
}

func (s *StatisticUseCase) TransactionPriority(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	usrEmail := ctx.MustGet("email").(string)
	personalAccount := personalaccounts.Informations(ctx, usrEmail)

	if personalAccount.ID == uuid.Nil {
		httpCode = http.StatusNotFound
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "not found")
		return response, httpCode, errInfo
	}

	response = s.repo.TransactionPriority(personalAccount.ID)
	return response, http.StatusOK, []errorsinfo.Errors{}
}

func (s *StatisticUseCase) Trend(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	usrEmail := ctx.MustGet("email").(string)
	personalAccount := personalaccounts.Informations(ctx, usrEmail)

	if personalAccount.ID == uuid.Nil {
		httpCode = http.StatusNotFound
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "not found")
		return response, httpCode, errInfo
	}

	response = s.repo.Trend(personalAccount.ID)
	return response, http.StatusOK, []errorsinfo.Errors{}
}

func (s *StatisticUseCase) Category(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var categoryUUID uuid.UUID

	category := ctx.Query("category")
	usrEmail := ctx.MustGet("email").(string)

	personalAccount := personalaccounts.Informations(ctx, usrEmail)

	if personalAccount.ID == uuid.Nil {
		httpCode = http.StatusNotFound
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "not found")
		return response, httpCode, errInfo
	}

	if category == "" {
		httpCode = http.StatusBadRequest
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "category required")
		return response, httpCode, errInfo
	}

	if category != "" {
		categoryUUID, _ = uuid.Parse(category)
	}

	response = s.repo.Category(personalAccount.ID, categoryUUID)
	return response, http.StatusOK, []errorsinfo.Errors{}
}
