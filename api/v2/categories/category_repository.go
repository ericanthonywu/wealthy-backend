package categories

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/wealthy-app/wealthy-backend/api/v2/categories/entities"
	"gorm.io/gorm"
)

type (
	CategoryRepository struct {
		db *gorm.DB
	}

	ICategoryRepository interface {
		GetCategoriesByAccountID(accountID uuid.UUID) (data []entities.CategoryInformation, err error)
		GetSubCategoryByCategoryID(accountUUID, categoryID uuid.UUID) (data []entities.SubCategoryInformation, err error)
	}
)

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetCategoriesByAccountID(accountID uuid.UUID) (data []entities.CategoryInformation, err error) {

	if err := r.db.Model(&entities.CategoryInformation{}).
		Select("expense_types as category_name, id as category_id, image_path as category_icon").
		Where("id_personal_accounts = ?", accountID).
		Where("active=?", true).
		Find(&data).Error; err != nil {
		logrus.Errorf(err.Error())
		return []entities.CategoryInformation{}, err
	}
	return data, nil
}

func (r *CategoryRepository) GetSubCategoryByCategoryID(accountUUID, categoryID uuid.UUID) (data []entities.SubCategoryInformation, err error) {
	if err := r.db.Model(&entities.SubCategoryInformation{}).
		Select("subcategories as sub_category_name, id as sub_category_id, image_path as sub_category_icon").
		Where("id_personal_accounts = ?", accountUUID).
		Where("id_master_expense_categories=?", categoryID).
		Where("active=?", true).
		Find(&data).Error; err != nil {
		logrus.Errorf(err.Error())
		return []entities.SubCategoryInformation{}, err
	}
	return data, nil
}