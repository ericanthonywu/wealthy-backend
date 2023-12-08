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
		BudgetInfo   Limit             `json:"budget_info"`
		SubCategory  []SubCategoryInfo `json:"sub_category_info"`
	}

	SubCategoryInfo struct {
		SubCategoryID   uuid.UUID `json:"sub_category_id"`
		SubCategoryName string    `json:"sub_category_name"`
		BudgetLimit     Limit     `json:"budget_limit"`
	}

	BudgetSetRequest struct {
		IDCategory               uuid.UUID `json:"category_id"`
		IDSubCategory            uuid.UUID `json:"sub_category_id"`
		IDMasterTransactionTypes uuid.UUID `json:"id_master_transaction_types,omitempty"`
		Amount                   int64     `json:"budget_amount"`
		Departure                string    `json:"departure,omitempty"`
		Arrival                  string    `json:"arrival,omitempty"`
		TravelStartDate          string    `json:"travel_start_date,omitempty"`
		TravelEndDate            string    `json:"travel_end_date,omitempty"`
		ImageBase64              string    `json:"image_base64,omitempty"`
		IDMasterExchangeCurrency string    `json:"id_master_exchange_currency,omitempty"`
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

	Trends struct {
		Period       string       `json:"period"`
		CategoryID   uuid.UUID    `json:"category_id"`
		CategoryName string       `json:"category_name"`
		BudgetInfo   Limit        `json:"budget_info"`
		Expense      ExpenseInfo  `json:"expense_info"`
		TrendsInfo   []TrendsInfo `json:"trends_info"`
	}

	ExpenseInfo struct {
		AverageDailySpending            Transaction `json:"average_daily_spending"`
		AverageDailySpendingRecommended Transaction `json:"average_daily_spending_recommended"`
		TransactionSpending             Transaction `json:"transaction_spending"`
		BudgetRemains                   Transaction `json:"budget_remains"`
	}

	TrendsInfo struct {
		StartDate         string      `json:"transaction_start_date"`
		EndDate           string      `json:"transaction_end_date"`
		TransactionAmount Transaction `json:"transaction_amount"`
	}

	LatestMonth struct {
		Details []LatestDetails `json:"details"`
	}

	LatestDetails struct {
		Period              string      `json:"period"`
		StatusComparison    string      `json:"status_comparison"`
		BudgetInfo          Limit       `json:"budget_info"`
		TransactionSpending Transaction `json:"transaction_spending"`
		Percentage          string      `json:"percentage"`
	}

	Travel struct {
		Details []TravelDetails `json:"details"`
	}

	TravelDetails struct {
		ID              uuid.UUID `json:"travel_id"`
		Departure       string    `json:"travel_departure"`
		Arrival         string    `json:"travel_arrival"`
		Budget          Amount    `json:"travel_budget"`
		TravelStartDate string    `json:"travel_start_date"`
		TravelEndDate   string    `json:"travel_end_date"`
		ImagePath       string    `json:"image_path"`
		Filename        string    `json:"filename"`
	}

	Amount struct {
		CurrencyCode string `json:"currency_code"`
		Value        int64  `json:"value"`
	}
)