package budgets

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/budgets/dtos"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/budgets/entities"
	"github.com/semicolon-indonesia/wealthy-backend/utils/errorsinfo"
	"github.com/semicolon-indonesia/wealthy-backend/utils/personalaccounts"
	"log"
	"net/http"
)

type (
	BudgetUseCase struct {
		repo IBudgetRepository
	}

	IBudgetUseCase interface {
		All(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		Overview(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		Category(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		LatestSixMonths(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
		Set(ctx *gin.Context, dtoRequest *dtos.BudgetSetRequest) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
	}
)

func NewBudgetUseCase(repo IBudgetRepository) *BudgetUseCase {
	return &BudgetUseCase{repo: repo}
}

func (s *BudgetUseCase) All(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		allCategories []entities.BudgetAllCategoriesEntities
		resp          []dtos.BudgetResponseAllCategories
		subCat        []dtos.BudgetSubCategories
	)
	month := ctx.Query("month")
	year := ctx.Query("year")

	usrEmail := ctx.MustGet("email").(string)
	personalAccount := personalaccounts.Informations(ctx, usrEmail)

	if personalAccount.ID == uuid.Nil {
		httpCode = http.StatusNotFound
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "not found")
		return response, httpCode, errInfo
	}

	allCategories = s.repo.All(personalAccount.ID, month, year)

	if len(allCategories) > 0 {
		for _, vcat := range allCategories {
			var temp dtos.BudgetResponseAllCategories
			if len(vcat.SubCategories) > 0 {
				if err := json.Unmarshal([]byte(vcat.SubCategories), &subCat); err != nil {
					log.Println(err.Error())
				}
				temp.SubCategories = subCat
			}

			if len(vcat.SubCategories) == 0 {
				temp.SubCategories = []dtos.BudgetSubCategories{}
			}

			temp.ID = vcat.ID
			temp.Categories = vcat.Categories
			temp.Total = vcat.Total
			resp = append(resp, temp)
		}
	}

	return resp, http.StatusOK, []errorsinfo.Errors{}
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

func (s *BudgetUseCase) Set(ctx *gin.Context, dtoRequest *dtos.BudgetSetRequest) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var (
		model       entities.BudgetSetEntities
		dtoResponse dtos.BudgetSetResponse
	)

	usrEmail := ctx.MustGet("email").(string)
	personalAccount := personalaccounts.Informations(ctx, usrEmail)

	if personalAccount.ID == uuid.Nil {
		httpCode = http.StatusNotFound
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "body payload required")
		return response, httpCode, errInfo
	}

	model.Amount = dtoRequest.Amount
	model.IDPersonalAccount = personalAccount.ID
	model.IDCategory = dtoRequest.IDCategory
	model.IDSubCategory = dtoRequest.IDSubCategory
	model.ID = uuid.New()

	err := s.repo.Set(&model)

	if err != nil {
		httpCode = http.StatusInternalServerError
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "problem while set budget")
		return response, httpCode, errInfo
	}

	dtoResponse.ID = model.ID
	dtoResponse.Status = true
	return dtoResponse, httpCode, []errorsinfo.Errors{}
}
