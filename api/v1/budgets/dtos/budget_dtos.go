package dtos

import "github.com/google/uuid"

type (
	AllBudgetLimit struct {
		Period          string            `json:"period"`
		AllBudgetDetail []AllBudgetDetail `json:"budget_info"`
	}

	AllBudgetDetail struct {
		CategoryID   uuid.UUID         `json:"category_id"`
		CategoryName string            `json:"category_name"`
		SubCategory  []SubCategoryInfo `json:"sub_category_info"`
	}

	SubCategoryInfo struct {
		SubCategoryID   uuid.UUID `json:"sub_category_id"`
		SubCategoryName string    `json:"sub_category_name"`
		BudgetLimit     Limit     `json:"budget_limit"`
	}

	BudgetSetRequest struct {
		IDCategory    uuid.UUID `json:"category_id"`
		IDSubCategory uuid.UUID `json:"sub_category_id"`
		Amount        int       `json:"budget_amount"`
	}

	BudgetSetResponse struct {
		ID     uuid.UUID `json:"budget_id"`
		Status bool      `json:"status"`
	}

	BudgetOverview struct {
		Period  string           `json:"period"`
		Details []OverviewDetail `json:"details"`
	}

	OverviewDetail struct {
		CategoryName        string      `json:"category_name"`
		CategoryID          uuid.UUID   `json:"category_id"`
		BudgetLimit         Limit       `json:"budget_limit"`
		TransactionSpending Transaction `json:"transaction_spending"`
		NumberOfCategories  int         `json:"number_of_categories"`
	}

	Limit struct {
		CurrencyCode string `json:"currency_code"`
		Value        int    `json:"value"`
	}

	Transaction struct {
		CurrencyCode string `json:"currency_code"`
		Value        int    `json:"value"`
	}
)