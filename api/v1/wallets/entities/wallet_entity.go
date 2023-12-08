package entities

import "github.com/google/uuid"

type (
	WalletPersonalInformationEntity struct {
		ID                   uuid.UUID `gorm:"column:id"`
		IDMasterAccountTypes uuid.UUID `gorm:"column:id_master_account_types"`
		AccountTypes         string    `gorm:"column:account_type"`
		TotalWallets         int64     `gorm:"column:total_wallet"`
	}

	WalletEntity struct {
		ID                 uuid.UUID `gorm:"column:id" json:"wallet_id"`
		IDAccount          uuid.UUID `gorm:"column:id_account" json:"id_personal_accounts"`
		WalletType         string    `gorm:"column:wallet_type" json:"wallet_type"`
		WalletName         string    `gorm:"column:wallet_name" json:"wallet_name"`
		IDMasterWalletType uuid.UUID `gorm:"column:id_master_wallet_types" json:"id_master_wallet_types"`
		Active             bool      `gorm:"column:active" json:"is_active"`
		FeeInvestBuy       float64   `gorm:"column:fee_invest_buy" json:"fee_invest_buy"`
		FeeInvestSell      float64   `gorm:"column:fee_invest_sell" json:"fee_invest_sell"`
		TotalAssets        int64     `gorm:"column:amount" json:"total_assets"`
	}

	WalletInitTransaction struct {
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

	WalletInitTransactionDetail struct {
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
)

func (WalletEntity) TableName() string {
	return "tbl_wallets"
}

func (WalletInitTransaction) TableName() string {
	return "tbl_transactions"
}

func (WalletInitTransactionDetail) TableName() string {
	return "tbl_transaction_details"
}