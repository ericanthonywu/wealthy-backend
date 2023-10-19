package dtos

import "github.com/google/uuid"

type (
	BudgetResponseAllCategories struct {
		ID            uuid.UUID             `json:"id"`
		Categories    string                `json:"categories"`
		Total         string                `json:"total"`
		SubCategories []BudgetSubCategories `json:"sub-categories"`
	}

	BudgetSubCategories struct {
		LimitAmount     int    `json:"limit_amount"`
		SubCategoryName string `json:"subcategory_name"`
	}
)
