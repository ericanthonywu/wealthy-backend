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
		GetCategoriesExpenseByAccountID(accountID uuid.UUID) (data []entities.CategoryExpenseInformation, err error)
		GetCategoriesIncomeByAccountID(accountID uuid.UUID) (data []entities.CategoryIncomeInformation, err error)
		GetSubCategoryExpenseByCategoryID(accountUUID, categoryID uuid.UUID) (data []entities.SubCategoryExpenseInformation, err error)
	}
)

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetCategoriesExpenseByAccountID(accountID uuid.UUID) (data []entities.CategoryExpenseInformation, err error) {
	if err := r.db.Model(&entities.CategoryExpenseInformation{}).
		Select("expense_types as category_name, id as category_id, image_path as category_icon").
		Where("id_personal_accounts = ?", accountID).
		Where("active=?", true).
		Find(&data).Error; err != nil {
		logrus.Errorf(err.Error())
		return []entities.CategoryExpenseInformation{}, err
	}
	return data, nil
}

func (r *CategoryRepository) GetCategoriesIncomeByAccountID(accountID uuid.UUID) (data []entities.CategoryIncomeInformation, err error) {
	if err := r.db.Model(&entities.CategoryIncomeInformation{}).
		Select("income_types as category_name, id as category_id, image_path as category_icon").
		Where("id_personal_accounts = ?", accountID).
		Where("active=?", true).
		Find(&data).Error; err != nil {
		logrus.Errorf(err.Error())
		return []entities.CategoryIncomeInformation{}, err
	}
	return data, nil
}

func (r *CategoryRepository) GetSubCategoryExpenseByCategoryID(accountUUID, categoryID uuid.UUID) (data []entities.SubCategoryExpenseInformation, err error) {
	if err := r.db.Model(&entities.SubCategoryExpenseInformation{}).
		Select("subcategories as sub_category_name, id as sub_category_id, image_path as sub_category_icon").
		Where("id_personal_accounts = ?", accountUUID).
		Where("id_master_expense_categories=?", categoryID).
		Where("active=?", true).
		Find(&data).Error; err != nil {
		logrus.Errorf(err.Error())
		return []entities.SubCategoryExpenseInformation{}, err
	}
	return data, nil
}