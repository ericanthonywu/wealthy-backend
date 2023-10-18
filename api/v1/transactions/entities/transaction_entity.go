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
		IDWallet                      uuid.UUID `gorm:"column:id_wallets"`
		IDMasterIncomeCategories      uuid.UUID `gorm:"column:id_master_income_categories"`
		IDMasterExpenseCategories     uuid.UUID `gorm:"column:id_master_expense_categories"`
		IDMasterInvest                uuid.UUID `gorm:"column:id_master_invest"`
		IDMasterBroker                uuid.UUID `gorm:"column:id_master_broker"`
		IDMasterReksanadaTypes        uuid.UUID `gorm:"column:id_master_reksadana_types"`
		IDMasterTransactionPriorities uuid.UUID `gorm:"column:id_master_transaction_priorities"`
		IDMasterTransactionTypes      uuid.UUID `gorm:"column:id_master_transaction_types"`
	}

	TransactionDetailEntity struct {
		IDTransaction     uuid.UUID `gorm:"column:id_transactions"`
		Repeat            bool      `gorm:"column:repeat"`
		Note              string    `gorm:"column:note"`
		From              string    `gorm:"column:from"`
		To                string    `gorm:"column:to"`
		MutualFundProduct string    `gorm:"column:mutual_fund_product"`
		StockCode         string    `gorm:"column:stock_code"`
		Lot               int64     `gorm:"column:lot"`
	}

	TransactionExpenseTotalHistory struct {
		TotalExpense int `gorm:"column:total_expense" json:"transaction_total"`
	}

	TransactionIncomeTotalHistory struct {
		TotalIncome int `gorm:"column:total_income" json:"transaction_total"`
	}

	TransactionDetailHistory struct {
		TransactionDate     string `gorm:"column:transaction_date" json:"transaction_date"`
		TransactionCategory string `gorm:"column:transaction_category" json:"transaction_category"`
		TransactionAmount   int    `gorm:"column:transaction_amount" json:"transaction_amount"`
		TransactionNote     string `gorm:"column:transaction_note" json:"transaction_note"`
	}

	TransactionDetailTransfer struct {
		TransactionDate        string `gorm:"column:transaction_date" json:"transaction_date"`
		TransactionAmount      int    `gorm:"column:transaction_amount" json:"transaction_amount"`
		TransactionNote        string `gorm:"column:transaction_note" json:"transaction_note"`
		TransactionDestination string `gorm:"column:transaction_destination" json:"transaction_destination"`
		TransactionSource      string `gorm:"column:transaction_source" json:"transaction_source"`
	}

	TransactionDetailInvest struct {
		TransactionDate        string `gorm:"column:transaction_date" json:"transaction_date"`
		TransactionAmountTotal int    `gorm:"column:transaction_amount_total" json:"transaction_amount_total"`
		TransactionNote        string `gorm:"column:transaction_note" json:"transaction_note"`
		Price                  int    `gorm:"column:price" json:"price"`
		Lot                    int    `gorm:"column:lot" json:"lot"`
		StockCode              string `gorm:"column:stock_code" json:"stock_code"`
		SellBuy                string `gorm:"column:sell_buy" json:"sell_buy"`
	}

	TransactionIncomeSpendingTotalMonthly struct {
		Month         string `gorm:"column:month" json:"month"`
		Year          int    `gorm:"column:year" json:"year"`
		TotalIncome   int    `gorm:"column:total_income" json:"total_income"`
		TotalSpending int    `gorm:"column:total_spending" json:"total_spending"`
		NetIncome     int    `gorm:"column:net_income" json:"net_income"`
	}

	TransactionIncomeSpendingDetailMonthly struct {
		TransactionCategory string `gorm:"column:transaction_category" json:"transaction_category"`
		TransactionType     string `gorm:"column:transaction_type" json:"transaction_type"`
		TransactionAmount   int    `gorm:"column:transaction_amount" json:"transaction_amount"`
		TransactionNote     string `gorm:"column:transaction_note" json:"transaction_note"`
	}

	TransactionIncomeSpendingTotalAnnually struct {
		TransactionPeriod string `gorm:"column:transaction_period" json:"transaction_period"`
		TotalIncome       int    `gorm:"column:total_income" json:"total_income"`
		TotalSpending     int    `gorm:"column:total_spending" json:"total_spending"`
		NetIncome         int    `gorm:"column:net_income" json:"net_income"`
	}

	TransactionIncomeSpendingDetailAnnually struct {
		MonthYear       string `gorm:"column:month_year" json:"month_year"`
		TotalDayInMonth int    `gorm:"column:total_day_in_month" json:"total_day_in_month"`
		TotalIncome     int    `gorm:"column:total_income" json:"total_income"`
		TotalSpending   int    `gorm:"column:total_spending" json:"total_spending"`
		NetIncome       int    `gorm:"column:net_income" json:"net_income"`
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
	}
)

func (TransactionEntity) TableName() string {
	return "tbl_transactions"
}

func (TransactionDetailEntity) TableName() string {
	return "tbl_transaction_details"
}
