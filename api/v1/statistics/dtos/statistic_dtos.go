package dtos

type (
	WeeklyData struct {
		Period     string             `json:"period"`
		Expense    []ExpenseWeekly    `json:"expense"`
		Income     []IncomeWeekly     `json:"income"`
		Investment []InvestmentWeekly `json:"investment"`
	}

	ExpenseWeekly struct {
		DateRange string `json:"date_range"`
		Amount    int    `json:"amount"`
	}

	IncomeWeekly struct {
		DateRange string `json:"date_range"`
		Amount    int    `json:"amount"`
	}

	InvestmentWeekly struct {
		DateRange string `json:"date_range"`
		Amount    int    `json:"amount"`
	}

	Summary struct {
		Period  string `json:"period"`
		Expense struct {
			TotalAmount int    `json:"total_amount"`
			Percentage  string `json:"percentage"`
		} `json:"expense"`
		Investment struct {
			TotalAmount int    `json:"total_amount"`
			Percentage  string `json:"percentage_string"`
		} `json:"investment"`
		NetIncome struct {
			TotalAmount int    `json:"total_amount"`
			Percentage  string `json:"percentage"`
		} `json:"net_income"`
	}

	TrendsData struct {
		Period        string          `json:"period"`
		AverageWeekly int             `json:"average_weekly"`
		AverageDaily  int             `json:"average_daily"`
		Expense       []ExpenseWeekly `json:"expense"`
	}

	Priority struct {
		Period string `json:"period"`
		Must   string `json:"transaction_must_type"`
		Want   string `json:"transaction_want_type"`
		Need   string `json:"transaction_need_type"`
	}
)