package dtos

import "github.com/google/uuid"

type (
	TransactionRequest struct {
		Date                          string    `json:"date_time_transaction"`
		Fees                          int64     `json:"fees,omitempty"`
		Amount                        int64     `json:"amount"`
		Repeat                        bool      `json:"repeat,omitempty"`
		Note                          string    `json:"note"`
		From                          string    `json:"from,omitempty"`
		To                            string    `json:"to,omitempty"`
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
)
