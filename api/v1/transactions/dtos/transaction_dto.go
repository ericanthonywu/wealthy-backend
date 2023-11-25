package dtos

import (
	"github.com/google/uuid"
)

type (
	TransactionRequest struct {
		Date                          string    `json:"date_time_transaction"`
		Fees                          int64     `json:"fees,omitempty"`
		Amount                        int64     `json:"amount"`
		Repeat                        bool      `json:"repeat,omitempty"`
		Note                          string    `json:"note"`
		TransferFrom                  string    `json:"transfer_from,omitempty"`
		TransferTo                    string    `json:"transfer_to,omitempty"`
		MutualFundProduct             string    `json:"mutual_fund_product,omitempty"`
		StockCode                     string    `json:"stock_code,omitempty"`
		Lot                           int64     `json:"lot,omitempty"`
		SellBuy                       int       `json:"sellbuy,omitempty"`
		IDWallet                      uuid.UUID `json:"id_wallets,omitempty"`
		IDMasterIncomeCategories      uuid.UUID `json:"id_master_income_categories,omitempty"`
		IDMasterExpenseCategories     uuid.UUID `json:"id_master_expense_categories,omitempty"`
		IDMasterExpenseSubCategories  uuid.UUID `json:"id_master_expense_subcategories,omitempty"`
		IDMasterInvest                uuid.UUID `json:"id_master_invest,omitempty"`
		IDMasterBroker                uuid.UUID `json:"id_master_broker,omitempty"`
		IDMasterReksanadaTypes        uuid.UUID `json:"id_master_reksadana_types,omitempty"`
		IDMasterTransactionPriorities uuid.UUID `json:"id_master_transaction_priorities,omitempty"`
		IDMasterTransactionTypes      uuid.UUID `json:"id_master_transaction_types,omitempty"`
		IDTravel                      uuid.UUID `json:"id_travel,omitempty"`
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
		Summary interface{}                                 `json:"summary"`
		Detail  []TransactionIncomeSpendingInvestmentDetail `json:"details"`
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

	TransactionDetails struct {
		TransactionCategory string `json:"transaction_category"`
		TransactionType     string `json:"transaction_type"`
		TransactionAmount   Amount `json:"transaction_amount"`
		TransactionNote     string `json:"transaction_note"`
	}
)