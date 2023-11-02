package entities

type (
	Statistic struct {
		ExpenseWeekOne   string `gorm:"column:expense_01-04" json:"expense_01-04"`
		ExpenseWeekTwo   string `gorm:"column:expense_05-11" json:"expense_05-11"`
		ExpenseWeekThree string `gorm:"column:expense_12-18" json:"expense_12-18"`
		ExpenseWeekFour  string `gorm:"column:expense_19-25" json:"expense_19-25"`
		ExpenseWeekFive  string `gorm:"column:expense_26-30" json:"expense_26-30"`
		IncomeWeekOne    string `gorm:"column:income_01-04" json:"income_01-04"`
		IncomeWeekTwo    string `gorm:"column:income_05-11" json:"income_05-11"`
		IncomeWeekThree  string `gorm:"column:income_12-18" json:"income_12-18"`
		IncomeWeekFour   string `gorm:"column:income_19-25" json:"income_19-25"`
		IncomeWeekFive   string `gorm:"column:income_26-30" json:"income_26-30"`
		InvestWeekOne    string `gorm:"column:invest_01-04" json:"invest_01-04"`
		InvestWeekTwo    string `gorm:"column:invest_05-11" json:"invest_05-11"`
		InvestWeekThree  string `gorm:"column:invest_12-18" json:"invest_12-18"`
		InvestWeekFour   string `gorm:"column:invest_19-25" json:"invest_19-25"`
		InvestWeekFive   string `gorm:"column:invest_26-30" json:"invest_26-30"`
		TotalIncome      string `gorm:"column:total_income" json:"total_income"`
		TotalExpense     string `gorm:"column:total_expense" json:"total_expense"`
		TotalNetIncome   string `gorm:"column:total_net_income" json:"total_net_income"`
		TotalInvest      string `gorm:"column:total_invest" json:"total_invest"`
	}

	StatisticTransactionPriority struct {
		TransactionNeedPercentage string `gorm:"column:transaction_need_percentage" json:"transaction_need_percentage"`
		TransactionMustPercentage string `gorm:"column:transaction_must_percentage" json:"transaction_must_percentage"`
		TransactionWantPercentage string `gorm:"column:transaction_want_percentage" json:"transaction_want_percentage"`
	}

	StatisticTrend struct {
		WeekOne            string `gorm:"column:01-04" json:"week_one_01-04"`
		WeekTwo            string `gorm:"column:05-11" json:"week_two_05-11"`
		WeekThree          string `gorm:"column:12-18" json:"week_three_12-18"`
		WeekFour           string `gorm:"column:19-25" json:"week_four_19-25"`
		WeekFive           string `gorm:"column:26-30" json:"week_five_26-30"`
		TotalAverageWeekly int    `gorm:"column:total_average_weekly" json:"total_average_weekly"`
		TotalAverageDaily  int    `gorm:"column:total_average_daily" json:"total_average_daily"`
	}

	StatisticExpenseWeekly struct {
		DateRange0104 int `gorm:"column:date_range_01_04"`
		DateRange0511 int `gorm:"column:date_range_05_11"`
		DateRange1218 int `gorm:"column:date_range_12_18"`
		DateRange1925 int `gorm:"column:date_range_19_25"`
		DateRange2630 int `gorm:"column:date_range_26_30"`
	}

	StatisticIncomeWeekly struct {
		DateRange0104 int `gorm:"column:date_range_01_04"`
		DateRange0511 int `gorm:"column:date_range_05_11"`
		DateRange1218 int `gorm:"column:date_range_12_18"`
		DateRange1925 int `gorm:"column:date_range_19_25"`
		DateRange2630 int `gorm:"column:date_range_26_30"`
	}

	StatisticInvestmentWeekly struct {
		DateRange0104 int `gorm:"column:date_range_01_04"`
		DateRange0511 int `gorm:"column:date_range_05_11"`
		DateRange1218 int `gorm:"column:date_range_12_18"`
		DateRange1925 int `gorm:"column:date_range_19_25"`
		DateRange2630 int `gorm:"column:date_range_26_30"`
	}

	StatisticSummaryMonthly struct {
		TotalExpense  int `gorm:"column:total_expense"`
		TotalIncome   int `gorm:"column:total_income"`
		TotalTransfer int `gorm:"column:total_transfer"`
		TotalInvest   int `gorm:"column:total_invest"`
	}
)