package dtos

import "github.com/google/uuid"

type (
	CategoryResponse struct {
		CategoryName    string                `json:"category_name"`
		CategoryID      uuid.UUID             `json:"category_id"`
		CategoryIcon    string                `json:"category_icon"`
		SubCategoryList []SubCategoryResponse `json:"sub_category"`
	}

	SubCategoryResponse struct {
		SubcategoryName string `json:"sub_cateogry_name"`
		SubcategoryID   string `json:"sub_cateogry_id"`
		SubcategoryIcon string `json:"sub_cateogry_icon"`
	}
)