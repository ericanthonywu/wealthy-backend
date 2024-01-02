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
		GetCatagoriesList(ginContext *gin.Context) (response []dtos.CategoryListResponse, httpCode int, errInfo []string)
	}
)

func NewCategoryUseCase(repo ICategoryRepository) *CategoryUseCase {
	return &CategoryUseCase{repo: repo}
}

func (s *CategoryUseCase) GetCatagoriesList(ginContext *gin.Context) (response []dtos.CategoryListResponse, httpCode int, errInfo []string) {
	var subCategoryList []dtos.SubCategoryList

	// get account ID
	accountUUID := ginContext.MustGet("accountID").(uuid.UUID)

	// fetch category by account id
	dataCategories, err := s.repo.GetCategoriesByAccountID(accountUUID)
	if err != nil {
		logrus.Errorf(err.Error())
		errInfo = errorsinfo.ErrorInfoWrapper(errInfo, err.Error())
		return []dtos.CategoryListResponse{}, http.StatusInternalServerError, errInfo
	}

	// if data categories not empty
	if len(dataCategories) > 0 {
		for _, v := range dataCategories {
			// get category id
			categoryID := v.CategoryID

			dataSubCategories, err := s.repo.GetSubCategoryByCategoryID(accountUUID, categoryID)
			if err != nil {
				logrus.Error(err.Error())
			}

			// if sub categories not empty
			if len(dataSubCategories) > 0 {
				for _, w := range dataSubCategories {
					subCategoryList = append(subCategoryList, dtos.SubCategoryList{
						SubcategoryName: w.SubCategoryName,
						SubcategoryID:   w.SubCategoryID,
						SubcategoryIcon: w.SubCategoryIcon,
					})
				}
			}

			// if sub categories empty
			if len(dataSubCategories) == 0 {
				subCategoryList = []dtos.SubCategoryList{}
			}

			// mapping response
			response = append(response, dtos.CategoryListResponse{
				CategoryName:    v.CategoryName,
				CategoryID:      v.CategoryID,
				CategoryIcon:    v.CategoryIcon,
				SubCategoryList: subCategoryList,
			})

			// reset
			subCategoryList = nil
		}
	}

	// if not errors
	if len(errInfo) == 0 {
		errInfo = []string{}
	}

	return response, http.StatusOK, errInfo
}