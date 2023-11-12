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
		Departure                     string    `json:"departure,omitempty"`
		Arrival                       string    `json:"arrival,omitempty"`
		MutualFundProduct             string    `json:"mutual_fund_product,omitempty"`
		StockCode                     string    `json:"stock_code,omitempty"`
		Lot                           int64     `json:"lot,omitempty"`
		IDWallet                      uuid.UUID `json:"id_wallets,omitempty"`
		IDMasterIncomeCategories      uuid.UUID `json:"id_master_income_categories,omitempty"`
		IDMasterExpenseCategories     uuid.UUID `json:"id_master_expense_categories,omitempty"`
		IDMasterInvest                uuid.UUID `json:"id_master_invest,omitempty"`
		IDMasterBroker                uuid.UUID `json:"id_master_broker,omitempty"`
		IDMasterReksanadaTypes        uuid.UUID `json:"id_master_reksadana_types,omitempty"`
		IDMasterTransactionPriorities uuid.UUID `json:"id_master_transaction_priorities,omitempty"`
		IDMasterTransactionTypes      uuid.UUID `json:"id_master_transaction_types,omitempty"`
	}

	TransactionTotalIncomeSpending struct {
		Month int `json:"month"`
	}

	TransactionHistoryForIncomeExpenses struct {
		Total  int         `json:"transaction_total"`
		Detail interface{} `json:"detail"`
	}

	TransactionHistoryForTransfer struct {
		TotalMoneyIn  int         `json:"total_money_in"`
		TotalMoneyOut int         `json:"total_money_out"`
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