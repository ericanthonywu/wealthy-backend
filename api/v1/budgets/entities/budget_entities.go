package entities

import "github.com/google/uuid"

type (
	BudgetAllCategoriesEntities struct {
		ID            uuid.UUID                   `gorm:"column:id"`
		Categories    string                      `gorm:"column:categories"`
		Total         string                      `gorm:"column:total"`
		SubCategories BudgetSubCategoriesEntities `gorm:"column:sub_categories type:jsonb"`
	}

	BudgetSubCategoriesEntities []struct {
		LimitAmount     int    `json:"limit_amount"`
		SubcategoryName string `json:"subcategory_name"`
	}
)
