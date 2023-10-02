package entities

import "github.com/google/uuid"

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
)

func (TransactionEntity) TableName() string {
	return "tbl_transactions"
}

func (TransactionDetailEntity) TableName() string {
	return "tbl_transaction_details"
}
