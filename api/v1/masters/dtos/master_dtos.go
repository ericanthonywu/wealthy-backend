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
)