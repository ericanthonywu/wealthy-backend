package budgets

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/semicolon-indonesia/wealthy-backend/utils/errorsinfo"
	"github.com/semicolon-indonesia/wealthy-backend/utils/personalaccounts"
	"net/http"
)

type (
	BudgetUseCase struct {
		repo IBudgetRepository
	}

	IBudgetUseCase interface {
		AllCategories(ctx *gin.Context)
		Overview(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		Category(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		LatestSixMonths(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		Set()
	}
)

func NewBudgetUseCase(repo IBudgetRepository) *BudgetUseCase {
	return &BudgetUseCase{repo: repo}
}

func (s *BudgetUseCase) AllCategories(ctx *gin.Context) {

	//s.repo.AllCategories()
	return
}

func (s *BudgetUseCase) Overview(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	month := ctx.Query("month")
	year := ctx.Query("year")

	usrEmail := ctx.MustGet("email").(string)
	personalAccount := personalaccounts.Informations(ctx, usrEmail)

	if personalAccount.ID == uuid.Nil {
		httpCode = http.StatusNotFound
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "not found")
		return response, httpCode, errInfo
	}

	if month == "" && year == "" {
		httpCode = http.StatusBadGateway
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "need month and or year information")
		return response, httpCode, errInfo
	}

	response = s.repo.Overview(personalAccount.ID, month, year)
	return response, http.StatusOK, []errorsinfo.Errors{}
}

func (s *BudgetUseCase) Category(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var categoryUUID uuid.UUID

	category := ctx.Query("category")
	month := ctx.Query("month")
	year := ctx.Query("year")

	if category != "" {
		categoryUUID, _ = uuid.Parse(category)
	}

	usrEmail := ctx.MustGet("email").(string)
	personalAccount := personalaccounts.Informations(ctx, usrEmail)

	if personalAccount.ID == uuid.Nil {
		httpCode = http.StatusNotFound
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "not found")
		return response, httpCode, errInfo
	}

	if category == "" || month == "" || year == "" {
		httpCode = http.StatusBadGateway
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "need month, year and id_category information")
		return response, httpCode, errInfo
	}

	response = s.repo.Category(personalAccount.ID, month, year, categoryUUID)
	return response, http.StatusOK, []errorsinfo.Errors{}
}

func (s *BudgetUseCase) LatestSixMonths(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var categoryUUID uuid.UUID

	category := ctx.Query("category")

	if category != "" {
		categoryUUID, _ = uuid.Parse(category)
	}

	usrEmail := ctx.MustGet("email").(string)
	personalAccount := personalaccounts.Informations(ctx, usrEmail)

	if personalAccount.ID == uuid.Nil {
		httpCode = http.StatusNotFound
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "not found")
		return response, httpCode, errInfo
	}

	if category == "" {
		httpCode = http.StatusBadGateway
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "need month, year and id_category information")
		return response, httpCode, errInfo
	}

	response = s.repo.LatestSixMonths(personalAccount.ID, categoryUUID)
	return response, http.StatusOK, []errorsinfo.Errors{}
}

func (s *BudgetUseCase) Set() {
	s.repo.Set()
}
