package entities

import "github.com/google/uuid"

type (
	CategoryInformation struct {
		CategoryName string    `json:"category_name"`
		CategoryID   uuid.UUID `json:"category_id"`
		CategoryIcon string    `json:"category_icon"`
	}

	SubCategoryInformation struct {
		SubCategoryName string `json:"sub_category_name"`
		SubCategoryID   string `json:"sub_category_id"`
		SubCategoryIcon string `json:"sub_category_icon"`
	}
)

func (CategoryInformation) TableName() string {
	return "tbl_master_expense_categories_editable"
}

func (SubCategoryInformation) TableName() string {
	return "tbl_master_expense_subcategories_editable"
}