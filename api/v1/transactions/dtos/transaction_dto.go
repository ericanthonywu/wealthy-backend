package dtos

import (
	"github.com/google/uuid"
)

type (
	TransactionRequest struct {
		Date                          string `json:"date_time_transaction"`
		Fees                          int64  `json:"fees,omitempty"`
		Amount                        int64  `json:"amount"`
		Repeat                        bool   `json:"repeat,omitempty"`
		Note                          string `json:"note"`
		TransferFrom                  string `json:"transfer_from,omitempty"`
		TransferTo                    string `json:"transfer_to,omitempty"`
		IDWallet                      string `json:"id_wallets,omitempty"`
		IDMasterIncomeCategories      string `json:"id_master_income_categories,omitempty"`
		IDMasterExpenseCategories     string `json:"id_master_expense_categories,omitempty"`
		IDMasterExpenseSubCategories  string `json:"id_master_expense_subcategories,omitempty"`
		IDMasterTransactionPriorities string `json:"id_master_transaction_priorities,omitempty"`
		IDMasterTransactionTypes      string `json:"id_master_transaction_types,omitempty"`
		IDTravel                      string `json:"id_travel,omitempty"`
	}

	TransactionRequestInvestment struct {
		Date                          string `json:"date_time_transaction"`
		Price                         int64  `json:"price"`
		StockCode                     string `json:"stock_code,omitempty"`
		Lot                           int64  `json:"lot,omitempty"`
		SellBuy                       int    `json:"sellbuy,omitempty"`
		IDWallet                      string `json:"id_wallets,omitempty"`
		IDMasterInvest                string `json:"id_master_invest,omitempty"`
		IDMasterTransactionPriorities string `json:"id_master_transaction_priorities,omitempty"`
		IDMasterTransactionTypes      string `json:"id_master_transaction_types,omitempty"`
	}

	TransactionTotalIncomeSpending struct {
		Month int `json:"month"`
	}

	TransactionHistoryForIncomeExpenses struct {
		Total  float64     `json:"transaction_total"`
		Detail interface{} `json:"detail"`
	}

	TransactionHistoryForTravel struct {
		Detail []TransactionHistoryForTravelDetail `json:"details"`
	}

	TransactionHistoryForTravelDetail struct {
		DateTransaction string    `json:"transaction_date"`
		IDTransaction   uuid.UUID `json:"transaction_id"`
		Amount          Amount    `json:"amount"`
		Category        string    `json:"transaction_category"`
		CategoryIcon    string    `json:"transaction_category_icon"`
		Note            string    `json:"transaction_note"`
	}

	Amount struct {
		CurrencyCode string  `json:"currency_code"`
		Value        float64 `json:"value"`
	}

	TransactionHistoryForTransfer struct {
		TotalMoneyIn  int         `json:"total_money_in,omitempty"`
		TotalMoneyOut int         `json:"total_money_out,omitempty"`
		Detail        interface{} `json:"detail"`
	}

	TransactionHistoryForInvest struct {
		Detail interface{} `json:"detail"`
	}

	TransactionIncomeSpendingInvestment struct {
		Summary TransactionSummaryIncomeSpending            `json:"summary"`
		Detail  []TransactionIncomeSpendingInvestmentDetail `json:"details"`
	}

	TransactionSummaryIncomeSpending struct {
		Month         string  `json:"month"`
		Year          int     `json:"year"`
		TotalIncome   float64 `json:"total_income"`
		TotalSpending float64 `json:"total_spending"`
		NetIncome     float64 `json:"net_income"`
	}

	TransactionIncomeSpendingInvestmentAnnually struct {
		Summary interface{}                 `json:"summary"`
		Detail  []TransactionDetailAnnually `json:"details"`
	}

	TransactionDetailAnnually struct {
		LastDayInMonth  int     `json:"last_day_month"`
		MonthYear       string  `json:"month_year"`
		TotalDayInMonth int     `json:"total_day_in_month"`
		TotalIncome     int     `json:"total_income"`
		TotalSpending   int     `json:"total_spending"`
		NetIncome       float64 `json:"net_income"`
	}

	TransactionIncomeSpendingInvestmentDetail struct {
		TransactionDate    string               `json:"transaction_date"`
		TransactionDetails []TransactionDetails `json:"transaction_details"`
	}

	TransactionInvestment struct {
		Summary interface{}                   `json:"summary"`
		Detail  []TransactionInvestmentDetail `json:"details"`
	}

	TransactionInvestmentDetail struct {
		TransactionDate    string                     `json:"transaction_date"`
		TransactionDetails []TransactionInvestDetails `json:"transaction_details"`
	}

	TransactionInvestDetails struct {
		LotWithPrice float64 `json:"lot_with_price"`
		Name         string  `json:"invest_name"`
		Lot          int     `json:"lot"`
		StockCode    string  `json:"stock_code"`
		Price        int64   `json:"price"`
		SellBuy      string  `json:"sell_buy"`
	}

	TransactionDetails struct {
		TransactionCategory     string `json:"transaction_category"`
		TransactionType         string `json:"transaction_type"`
		TransactionCategoryIcon string `json:"transaction_category_icon"`
		TransactionAmount       Amount `json:"transaction_amount"`
		TransactionNote         string `json:"transaction_note"`
	}

	TransactionNotes struct {
		TransactionDate        string                   `json:"transaction_date"`
		TransactionNotesDetail []TransactionNotesDetail `json:"details"`
	}

	TransactionNotesDetail struct {
		TransactionCategory        string                       `json:"transaction_category"`
		TransactionNotesDeepDetail []TransactionNotesDeepDetail `json:"info"`
	}

	TransactionNotesDeepDetail struct {
		TransactionNote   string `json:"transaction_note"`
		TransactionAmount Amount `json:"transaction_amount"`
		TransactionLimit  Amount `json:"transaction_limit"`
	}

	CashFlowResponse struct {
		CountIncome         int64   `json:"count_income"`
		CountExpense        int64   `json:"count_expense"`
		CashFlow            float64 `json:"cashflow"`
		TotalAverageIncome  float64 `json:"total_income"`
		TotalAverageExpense float64 `json:"total_expense"`
		AverageDay          struct {
			Income  float64 `json:"income"`
			Expense float64 `json:"expense"`
		} `json:"average_day"`
		AverageMonth struct {
			Income  float64 `json:"income"`
			Expense float64 `json:"expense"`
		} `json:"average_month"`
	}

	WalletListResponse struct {
		IDAccount     uuid.UUID     `json:"id_personal_accounts"`
		WalletDetails WalletDetails `json:"details"`
	}

	WalletDetails struct {
		WalletID           uuid.UUID `json:"wallet_id"`
		WalletType         string    `json:"wallet_type"`
		WalletName         string    `json:"wallet_name"`
		IDMasterWalletType uuid.UUID `json:"id_master_wallet_types"`
	}
)