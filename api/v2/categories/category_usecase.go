package categories

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/wealthy-app/wealthy-backend/api/v2/categories/dtos"
	"github.com/wealthy-app/wealthy-backend/utils/errorsinfo"
	"net/http"
)

type (
	CategoryUseCase struct {
		repo ICategoryRepository
	}

	ICategoryUseCase interface {
		GetCategoriesExpenseList(ginContext *gin.Context) (response []dtos.CategoryResponse, httpCode int, errInfo []string)
		GetCategoriesIncomeList(ginContext *gin.Context) (response []dtos.CategoryResponse, httpCode int, errInfo []string)
	}
)

func NewCategoryUseCase(repo ICategoryRepository) *CategoryUseCase {
	return &CategoryUseCase{repo: repo}
}

func (s *CategoryUseCase) GetCategoriesExpenseList(ginContext *gin.Context) (response []dtos.CategoryResponse, httpCode int, errInfo []string) {
	var subCategoryResponse []dtos.SubCategoryResponse

	// get account ID
	accountUUID := ginContext.MustGet("accountID").(uuid.UUID)

	// fetch category by account id
	dataCategories, err := s.repo.GetCategoriesExpenseByAccountID(accountUUID)
	if err != nil {
		logrus.Errorf(err.Error())
		errInfo = errorsinfo.ErrorInfoWrapper(errInfo, err.Error())
		return []dtos.CategoryResponse{}, http.StatusInternalServerError, errInfo
	}

	// if data categories not empty
	if len(dataCategories) > 0 {
		for _, v := range dataCategories {
			// get category id
			categoryID := v.CategoryID

			dataSubCategories, err := s.repo.GetSubCategoryExpenseByCategoryID(accountUUID, categoryID)
			if err != nil {
				logrus.Error(err.Error())
			}

			// if sub categories not empty
			if len(dataSubCategories) > 0 {
				for _, w := range dataSubCategories {
					subCategoryResponse = append(subCategoryResponse, dtos.SubCategoryResponse{
						SubcategoryName: w.SubCategoryName,
						SubcategoryID:   w.SubCategoryID,
						SubcategoryIcon: w.SubCategoryIcon,
					})
				}
			}

			// if sub categories empty
			if len(dataSubCategories) == 0 {
				subCategoryResponse = []dtos.SubCategoryResponse{}
			}

			// mapping response
			response = append(response, dtos.CategoryResponse{
				CategoryName:    v.CategoryName,
				CategoryID:      v.CategoryID,
				CategoryIcon:    v.CategoryIcon,
				SubCategoryList: subCategoryResponse,
			})

			// reset
			subCategoryResponse = nil
		}
	}

	// if not errors
	if len(errInfo) == 0 {
		errInfo = []string{}
	}

	return response, http.StatusOK, errInfo
}

func (s *CategoryUseCase) GetCategoriesIncomeList(ginContext *gin.Context) (response []dtos.CategoryResponse, httpCode int, errInfo []string) {
	var subCategoryResponse []dtos.SubCategoryResponse

	// get account ID
	accountUUID := ginContext.MustGet("accountID").(uuid.UUID)

	// fetch category by account id
	dataCategories, err := s.repo.GetCategoriesIncomeByAccountID(accountUUID)
	if err != nil {
		logrus.Errorf(err.Error())
		errInfo = errorsinfo.ErrorInfoWrapper(errInfo, err.Error())
		return []dtos.CategoryResponse{}, http.StatusInternalServerError, errInfo
	}

	// if data categories not empty
	if len(dataCategories) > 0 {
		for _, v := range dataCategories {
			// get category id
			categoryID := v.CategoryID

			dataSubCategories, err := s.repo.GetSubCategoryExpenseByCategoryID(accountUUID, categoryID)
			if err != nil {
				logrus.Error(err.Error())
			}

			// if sub categories not empty
			if len(dataSubCategories) > 0 {
				for _, w := range dataSubCategories {
					subCategoryResponse = append(subCategoryResponse, dtos.SubCategoryResponse{
						SubcategoryName: w.SubCategoryName,
						SubcategoryID:   w.SubCategoryID,
						SubcategoryIcon: w.SubCategoryIcon,
					})
				}
			}

			// if sub categories empty
			if len(dataSubCategories) == 0 {
				subCategoryResponse = []dtos.SubCategoryResponse{}
			}

			// mapping response
			response = append(response, dtos.CategoryResponse{
				CategoryName:    v.CategoryName,
				CategoryID:      v.CategoryID,
				CategoryIcon:    v.CategoryIcon,
				SubCategoryList: subCategoryResponse,
			})

			// reset
			subCategoryResponse = nil
		}
	}

	// if not errors
	if len(errInfo) == 0 {
		errInfo = []string{}
	}

	return response, http.StatusOK, errInfo
}