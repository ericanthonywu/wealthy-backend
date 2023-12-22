package dtos

import "github.com/google/uuid"

type (
	RenameCatRequest struct {
		NewCategoryName string `json:"new_category_name"`
	}

	AddCategory struct {
		ExpenseID    uuid.UUID `json:"expense_id,omitempty"`
		CategoryName string    `json:"category_name"`
	}

	WalletResponse struct {
		ID         uuid.UUID `json:"wallet_id"`
		WalletName string    `json:"wallet_name"`
	}

	PriceResponse struct {
		ID            uuid.UUID `json:"id"`
		Title         string    `json:"title"`
		Price         float64   `json:"price"`
		ActualPrice   float64   `json:"actual_price"`
		Description   string    `json:"description"`
		IsCurrent     bool      `json:"is_current"`
		IsRecommended bool      `json:"is_recommended"`
	}
)