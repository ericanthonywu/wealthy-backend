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
		MutualFundProduct             string `json:"mutual_fund_product,omitempty"`
		StockCode                     string `json:"stock_code,omitempty"`
		Lot                           int64  `json:"lot,omitempty"`
		SellBuy                       int    `json:"sellbuy,omitempty"`
		IDWallet                      string `json:"id_wallets,omitempty"`
		IDMasterIncomeCategories      string `json:"id_master_income_categories,omitempty"`
		IDMasterExpenseCategories     string `json:"id_master_expense_categories,omitempty"`
		IDMasterExpenseSubCategories  string `json:"id_master_expense_subcategories,omitempty"`
		IDMasterInvest                string `json:"id_master_invest,omitempty"`
		IDMasterBroker                string `json:"id_master_broker,omitempty"`
		IDMasterReksanadaTypes        string `json:"id_master_reksadana_types,omitempty"`
		IDMasterTransactionPriorities string `json:"id_master_transaction_priorities,omitempty"`
		IDMasterTransactionTypes      string `json:"id_master_transaction_types,omitempty"`
		IDTravel                      string `json:"id_travel,omitempty"`
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
	}

	TransactionDetails struct {
		TransactionCategory string `json:"transaction_category"`
		TransactionType     string `json:"transaction_type"`
		TransactionAmount   Amount `json:"transaction_amount"`
		TransactionNote     string `json:"transaction_note"`
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
)