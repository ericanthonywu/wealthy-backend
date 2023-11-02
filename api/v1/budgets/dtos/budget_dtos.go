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

	BudgetSetRequest struct {
		IDCategory    uuid.UUID `json:"id_master_categories"`
		IDSubCategory uuid.UUID `json:"id_master_subcategories"`
		Amount        int       `json:"amount"`
	}

	BudgetSetResponse struct {
		ID     uuid.UUID `json:"id"`
		Status bool      `json:"status"`
	}

	BudgetOverview struct {
		TransactionCategory string `json:"transaction_category"`
		BudgetLimit         int    `json:"budget_limit"`
		TotalSpending       int    `json:"total_spending"`
		NumberOfCategory    int    `json:"number_of_category"`
	}
)
