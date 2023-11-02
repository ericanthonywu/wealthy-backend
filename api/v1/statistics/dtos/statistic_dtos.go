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
)