package entities

import "github.com/google/uuid"

type (
	CategoryExpenseInformation struct {
		CategoryName string    `json:"category_name"`
		CategoryID   uuid.UUID `json:"category_id"`
		CategoryIcon string    `json:"category_icon"`
	}

	CategoryIncomeInformation struct {
		CategoryName string    `json:"category_name"`
		CategoryID   uuid.UUID `json:"category_id"`
		CategoryIcon string    `json:"category_icon"`
	}

	SubCategoryExpenseInformation struct {
		SubCategoryName string `json:"sub_category_name"`
		SubCategoryID   string `json:"sub_category_id"`
		SubCategoryIcon string `json:"sub_category_icon"`
	}

	SubCategoryIncomeInformation struct{}
)

func (CategoryExpenseInformation) TableName() string {
	return "tbl_master_expense_categories_editable"
}

func (CategoryIncomeInformation) TableName() string {
	return "tbl_master_income_categories_editable"
}

func (SubCategoryExpenseInformation) TableName() string {
	return "tbl_master_expense_subcategories_editable"
}