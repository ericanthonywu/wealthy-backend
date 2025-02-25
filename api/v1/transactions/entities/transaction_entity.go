package entities

import (
	"github.com/google/uuid"
)

type (
	TransactionEntity struct {
		ID                            uuid.UUID `gorm:"column:id"`
		Date                          string    `gorm:"column:date_time_transaction"`
		Fees                          float64   `gorm:"column:fees"`
		Amount                        float64   `gorm:"column:amount"`
		IDPersonalAccount             uuid.UUID `gorm:"column:id_personal_account"`
		IDWallet                      uuid.UUID `gorm:"column:id_wallets"`
		IDMasterIncomeCategories      uuid.UUID `gorm:"column:id_master_income_categories"`
		IDMasterExpenseCategories     uuid.UUID `gorm:"column:id_master_expense_categories"`
		IDMasterExpenseSubCategories  uuid.UUID `gorm:"column:id_master_expense_subcategories"`
		IDMasterInvest                uuid.UUID `gorm:"column:id_master_invest"`
		IDMasterReksanadaTypes        uuid.UUID `gorm:"column:id_master_reksadana_types"`
		IDMasterTransactionPriorities uuid.UUID `gorm:"column:id_master_transaction_priorities"`
		IDMasterTransactionTypes      uuid.UUID `gorm:"column:id_master_transaction_types"`
		Balance                       float64   `gorm:"column:balance"`
		Debit                         float64   `gorm:"column:debit"`
		Credit                        float64   `gorm:"column:credit"`
	}

	TransactionDetailEntity struct {
		IDTransaction     uuid.UUID `gorm:"column:id_transactions"`
		Repeat            bool      `gorm:"column:repeat"`
		Note              string    `gorm:"column:note"`
		From              string    `gorm:"column:transfer_from"`
		To                string    `gorm:"column:transfer_to"`
		MutualFundProduct string    `gorm:"column:mutual_fund_product"`
		StockCode         string    `gorm:"column:stock_code"`
		Lot               int64     `gorm:"column:lot"`
		SellBuy           int       `gorm:"column:sellbuy"`
		IDTravel          uuid.UUID `gorm:"column:id_travel"`
	}

	TransactionInvestmentEntity struct {
		StockCode          string    `gorm:"column:stock_code"`
		TotalLot           int64     `gorm:"column:total_lot"`
		ValueBuy           float64   `gorm:"column:value_buy"`
		FeeBuy             float64   `gorm:"column:fee_buy"`
		NetBuy             float64   `gorm:"column:net_buy"`
		AverageBuy         float64   `gorm:"column:average_buy"`
		InitialInvestment  float64   `gorm:"column:initial_investment"`
		IDPersonalAccount  uuid.UUID `gorm:"column:id_personal_accounts"`
		IDMasterBroker     uuid.UUID `gorm:"column:id_master_broker"`
		GainLoss           float64   `gorm:"column:gain_loss"`
		GainLossPercentage float64   `gorm:"column:gain_loss_percentage"`
		PotentialReturn    float64   `gorm:"column:potential_return"`
		PercentageReturn   float64   `gorm:"column:percentage_return"`
		WalletID           uuid.UUID `gorm:"column:wallet_id"`
		DateTime           string    `gorm:"column:date_time"`
	}

	TransactionExpenseTotalHistory struct {
		TotalExpense float64 `gorm:"column:total_expense" json:"transaction_total"`
	}

	TransactionIncomeTotalHistory struct {
		TotalIncome float64 `gorm:"column:total_income" json:"transaction_total"`
	}

	TransactionInvestTotalHistory struct {
		TotalInvest float64 `gorm:"column:total_invest" json:"transaction_total"`
	}

	TransactionDetailHistory struct {
		TransactionDate         string `gorm:"column:transaction_date" json:"transaction_date"`
		TransactionCategory     string `gorm:"column:transaction_category" json:"transaction_category"`
		TransactionCategoryIcon string `gorm:"column:transaction_category_icon" json:"transaction_category_icon"`
		TransactionAmount       int    `gorm:"column:transaction_amount" json:"transaction_amount"`
		TransactionNote         string `gorm:"column:transaction_note" json:"transaction_note"`
	}

	TransactionDetailTransfer struct {
		TransactionDate         string `gorm:"column:transaction_date" json:"transaction_date"`
		TransactionCategoryIcon string `gorm:"column:transaction_category_icon" json:"transaction_category_icon"`
		TransactionAmount       int    `gorm:"column:transaction_amount" json:"transaction_amount"`
		TransactionNote         string `gorm:"column:transaction_note" json:"transaction_note"`
		TransactionDestination  string `gorm:"column:transaction_destination" json:"transaction_destination"`
		TransactionSource       string `gorm:"column:transaction_source" json:"transaction_source"`
	}

	TransactionDetailInvest struct {
		TransactionDate         string  `gorm:"column:transaction_date" json:"transaction_date"`
		TransactionCategoryIcon string  `gorm:"column:transaction_category_icon" json:"transaction_category_icon"`
		TransactionAmountTotal  float64 `gorm:"column:transaction_amount_total" json:"transaction_amount_total"`
		TransactionNote         string  `gorm:"column:transaction_note" json:"transaction_note"`
		Price                   float64 `gorm:"column:price" json:"price"`
		Lot                     int     `gorm:"column:lot" json:"lot"`
		StockCode               string  `gorm:"column:stock_code" json:"stock_code"`
		SellBuy                 string  `gorm:"column:sell_buy" json:"sell_buy"`
	}

	TransactionDetailTravel struct {
		DateTransaction    string    `gorm:"column:date_time_transaction"`
		IDTransaction      uuid.UUID `gorm:"column:id_transaction"`
		TransactionCatIcon string    `gorm:"column:transaction_category_icon"`
		Amount             int64     `gorm:"column:amount"`
		Category           string    `gorm:"column:category"`
		Note               string    `gorm:"column:note"`
	}

	TransactionIncomeSpendingTotalMonthly struct {
		Month         string  `gorm:"column:month" json:"month"`
		Year          int     `gorm:"column:year" json:"year"`
		TotalIncome   int     `gorm:"column:total_income" json:"total_income"`
		TotalSpending int     `gorm:"column:total_spending" json:"total_spending"`
		NetIncome     float64 `gorm:"column:net_income" json:"net_income"`
	}

	TransactionIncomeSpendingDetailMonthly struct {
		TransactionDate         string `gorm:"column:date" json:"transaction_date"`
		TransactionCategory     string `gorm:"column:transaction_category" json:"transaction_category"`
		TransactionCategoryIcon string `gorm:"column:image_path" json:"transaction_category_icon"`
		TransactionType         string `gorm:"column:transaction_type" json:"transaction_type"`
		TransactionAmount       int    `gorm:"column:transaction_amount" json:"transaction_amount"`
		TransactionNote         string `gorm:"column:transaction_note" json:"transaction_note"`
	}

	TransactionIncomeSpendingTotalAnnually struct {
		TransactionPeriod string  `gorm:"column:transaction_period" json:"transaction_period"`
		TotalIncome       int     `gorm:"column:total_income" json:"total_income"`
		TotalSpending     int     `gorm:"column:total_spending" json:"total_spending"`
		NetIncome         float64 `gorm:"column:net_income" json:"net_income"`
	}

	TransactionIncomeSpendingDetailAnnually struct {
		DateOrigin      string  `gorm:"date_origin"`
		MonthYear       string  `gorm:"column:month_year" json:"month_year"`
		TotalDayInMonth int     `gorm:"column:total_day_in_month" json:"total_day_in_month"`
		TotalIncome     int     `gorm:"column:total_income" json:"total_income"`
		TotalSpending   int     `gorm:"column:total_spending" json:"total_spending"`
		NetIncome       float64 `gorm:"column:net_income" json:"net_income"`
	}

	TransactionInvestmentTotals struct {
		TotalBuy              int `gorm:"column:total_buy" json:"total_buy"`
		TotalSell             int `gorm:"column:total_sell" json:"total_sell"`
		TotalCurrentPortfolio int `gorm:"column:total_current_portfolio" json:"total_current_portfolio"`
	}

	TransactionInvestmentDetail struct {
		Date      string `gorm:"column:date" json:"date"`
		TotalBuy  int    `gorm:"column:total_buy" json:"total_buy"`
		TotalSell int    `gorm:"column:total_sell" json:"total_sell"`
		Lot       int    `gorm:"column:lot" json:"lot"`
		StockCode string `gorm:"column:stock_code" json:"stock_code"`
		Price     int64  `gorm:"column:price"`
		SellBuy   int    `gorm:"column:sell_buy"`
	}

	TransactionByNotes struct {
		Budget              float64   `gorm:"column:budget"`
		Amount              float64   `gorm:"column:amount"`
		TransactionNote     string    `gorm:"column:transaction_note"`
		TransactionCategory string    `gorm:"column:expense_types"`
		CategoryID          uuid.UUID `gorm:"column:category_id"`
	}

	TransactionSuggestionNotes struct {
		Note string `gorm:"column:note"`
	}

	TransactionWalletExist struct {
		Exists bool `gorm:"column:exists"`
	}

	TransactionWithCurrency struct {
		CurrencyValue int64 `gorm:"column:currency_value"`
	}

	InvestmentTreding struct {
		Name   string `gorm:"column:name"`
		Symbol string `gorm:"column:symbol"`
		Close  int64  `gorm:"column:close"`
	}

	BrokerInfo struct {
		ID   uuid.UUID `gorm:"column:id"`
		Name string    `gorm:"column:broker_name"`
	}

	IncomeEachDay struct {
		IncomeAverage float64 `gorm:"column:income_average"`
	}

	ExpenseEachDay struct {
		ExpenseAverage float64 `gorm:"column:expense_average"`
	}

	IncomeMonthly struct {
		IncomeAverage float64 `gorm:"column:income_average"`
	}

	ExpenseMonthly struct {
		ExpenseAverage float64 `gorm:"column:expense_average"`
	}

	CountIncomeTrxMonthly struct {
		CountIncome int64 `gorm:"column:count_income"`
	}

	CountExpenseTrxMonthly struct {
		CountExpense int64 `gorm:"column:count_expense"`
	}

	TotalIncomeCashflow struct {
		TotalIncome float64 `gorm:"column:total_income"`
	}

	TotalExpenseCashflow struct {
		TotalExpense float64 `gorm:"column:total_expense"`
	}

	LastBalance struct {
		ID                            uuid.UUID `gorm:"column:id"`
		Date                          string    `gorm:"column:date_time_transaction"`
		Fees                          float64   `gorm:"column:fees"`
		Amount                        float64   `gorm:"column:amount"`
		IDPersonalAccount             uuid.UUID `gorm:"column:id_personal_account"`
		IDWallet                      uuid.UUID `gorm:"column:id_wallets"`
		IDMasterIncomeCategories      uuid.UUID `gorm:"column:id_master_income_categories"`
		IDMasterExpenseCategories     uuid.UUID `gorm:"column:id_master_expense_categories"`
		IDMasterExpenseSubCategories  uuid.UUID `gorm:"column:id_master_expense_subcategories"`
		IDMasterInvest                uuid.UUID `gorm:"column:id_master_invest"`
		IDMasterBroker                uuid.UUID `gorm:"column:id_master_broker"`
		IDMasterReksanadaTypes        uuid.UUID `gorm:"column:id_master_reksadana_types"`
		IDMasterTransactionPriorities uuid.UUID `gorm:"column:id_master_transaction_priorities"`
		IDMasterTransactionTypes      uuid.UUID `gorm:"column:id_master_transaction_types"`
		Credit                        float64   `gorm:"column:credit"`
		Debit                         float64   `gorm:"column:debit"`
		Balance                       float64   `gorm:"column:balance"`
	}

	WalletEntity struct {
		ID                 uuid.UUID `gorm:"column:id"`
		IDAccount          uuid.UUID `gorm:"column:id_account"`
		WalletType         string    `gorm:"column:wallet_type"`
		WalletName         string    `gorm:"column:wallet_name"`
		IDMasterWalletType uuid.UUID `gorm:"column:id_master_wallet_types"`
		Active             bool      `gorm:"column:active"`
		FeeInvestBuy       float64   `gorm:"column:fee_invest_buy"`
		FeeInvestSell      float64   `gorm:"column:fee_invest_sell"`
		TotalAssets        int64     `gorm:"column:amount"`
	}

	Budget struct {
		IDPersonalAccount uuid.UUID `gorm:"column:id_personal_accounts"`
	}

	StockCodeData struct {
		StockCode string `gorm:"column:stock_code"`
	}

	BudgetEachCategory struct {
		Category    string `gorm:"column:category_name"`
		BudgetLimit int    `gorm:"column:total_budget_limit"`
	}

	IncomeSpendingMonthly struct {
		DateTransaction string    `gorm:"column:transaction_date"`
		Amount          float64   `gorm:"column:transaction_amount"`
		Note            string    `gorm:"column:transaction_note"`
		Type            string    `gorm:"column:transaction_type"`
		IDCategory      uuid.UUID `gorm:"column:id_category"`
	}

	Category struct {
		CategoryName string    `gorm:"column:category_name"`
		CategoryID   uuid.UUID `gorm:"column:category_id"`
		CategoryIcon string    `gorm:"column:category_icon"`
	}

	TotalSpendingMonthly struct {
		Amount float64 `gorm:"column:total_spending"`
	}

	TotalIncomeMonthly struct {
		Amount float64 `gorm:"column:total_income"`
	}

	TotalSpendingAnnual struct {
		Amount float64 `gorm:"column:total_spending"`
	}

	TotalIncomeAnnual struct {
		Amount float64 `gorm:"column:total_income"`
	}
)

func (TransactionEntity) TableName() string {
	return "tbl_transactions"
}

func (TransactionDetailEntity) TableName() string {
	return "tbl_transaction_details"
}

func (TransactionInvestmentEntity) TableName() string {
	return "tbl_investment"
}

func (BrokerInfo) TableName() string {
	return "tbl_master_broker"
}

func (LastBalance) TableName() string {
	return "tbl_transactions"
}