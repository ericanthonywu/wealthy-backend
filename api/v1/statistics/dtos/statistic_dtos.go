package dtos

import "github.com/google/uuid"

type (
	WeeklyData struct {
		Period     string             `json:"period"`
		Expense    []ExpenseWeekly    `json:"expense_details"`
		Income     []IncomeWeekly     `json:"income_details"`
		Investment []InvestmentWeekly `json:"investment_details"`
	}

	ExpenseWeekly struct {
		StartDate string             `json:"transaction_start_date"`
		EndDate   string             `json:"transaction_end_date"`
		Amount    ExpenseTransaction `json:"transaction_amount"`
	}

	ExpenseTransaction struct {
		CurrencyCode string `json:"currency_code"`
		Value        int    `json:"value"`
	}

	IncomeWeekly struct {
		StartDate string            `json:"transaction_start_date"`
		EndDate   string            `json:"transaction_end_date"`
		Amount    IncomeTransaction `json:"transaction_amount"`
	}

	IncomeTransaction struct {
		CurrencyCode string `json:"currency_code"`
		Value        int    `json:"value"`
	}

	InvestmentWeekly struct {
		StartDate string            `json:"transaction_start_date"`
		EndDate   string            `json:"transaction_end_date"`
		Amount    InvestTransaction `json:"transaction_amount"`
	}

	InvestTransaction struct {
		CurrencyCode string `json:"currency_code"`
		Value        int    `json:"value"`
	}

	Summary struct {
		Period  string `json:"period"`
		Expense struct {
			TotalAmount SummaryTransaction `json:"transaction_amount"`
			Percentage  string             `json:"percentage_increase"`
		} `json:"expense"_info`
		Investment struct {
			TotalAmount SummaryTransaction `json:"transaction_amount"`
			Percentage  string             `json:"percentage_increase"`
		} `json:"investment_info"`
		NetIncome struct {
			TotalAmount SummaryTransaction `json:"transaction_amount"`
			Percentage  string             `json:"percentage_increase"`
		} `json:"net_income_info"`
	}

	SummaryTransaction struct {
		CurrencyCode string `json:"currency_code"`
		Value        int    `json:"value"`
	}

	TrendsData struct {
		Period        string          `json:"period"`
		AverageWeekly int             `json:"average_weekly"`
		AverageDaily  int             `json:"average_daily"`
		Expense       []ExpenseWeekly `json:"expense"`
	}

	Priority struct {
		Period string         `json:"period"`
		Info   []PriorityInfo `json:"details"`
	}

	PriorityInfo struct {
		Type       string `json:"transaction_priority"`
		Percentage string `json:"percentage"`
	}

	ExpenseDetail struct {
		Period       string      `json:"period"`
		TotalExpense int64       `json:"total_expense"`
		Expense      []ExpDetail `json:"details"`
	}

	ExpDetail struct {
		ID           uuid.UUID            `json:"category_id"`
		Category     string               `json:"category_name"`
		CategoryIcon string               `json:"category_icon"`
		Amount       ExpDetailTransaction `json:"transaction_amount"`
	}

	ExpDetailTransaction struct {
		CurrencyCode string `json:"currency_code"`
		Value        int64  `json:"value"`
	}

	WeeklySubExpense struct {
		Period       string                   `json:"period"`
		CategoryID   string                   `json:"category_id"`
		CategoryName string                   `json:"category_name"`
		Expense      []WeeklySubExpenseDetail `json:"details"`
	}

	WeeklySubExpenseDetail struct {
		StartDate string                            `json:"transaction_start_date"`
		EndDate   string                            `json:"transaction_end_date"`
		Amount    WeeklySubExpenseDetailTransaction `json:"transaction_amount"`
	}

	WeeklySubExpenseDetailTransaction struct {
		CurrencyCode string `json:"currency_code"`
		Value        int    `json:"value"`
	}
)