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
		Departure       string `json:"departure"`
		Arrival         string `json:"arrival"`
		Amount          Amount `json:"budget_amount"`
		TravelStartDate string `json:"travel_start_date"`
		TravelEndDate   string `json:"travel_end_date"`
		ImagePath       string `json:"image_path"`
		Filename        string `json:"filename"`
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
		Summary interface{} `json:"summary"`
		Detail  interface{} `json:"detail"`
	}
)