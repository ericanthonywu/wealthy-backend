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
		ID            uuid.UUID `gorm:"column:id" json:"id"`
		IDAccount     uuid.UUID `gorm:"id_account" json:"id_account"`
		InvestType    string    `gorm:"column:invest_type" json:"invest_type"`
		InvestName    string    `gorm:"column:invest_name" json:"invest_name"`
		WalletType    string    `gorm:"column:id_master_wallet_types" json:"id_master_wallet_types"`
		Active        bool      `gorm:"column:active" json:"active"`
		FeeInvestBuy  int64     `gorm:"column:fee_invest_buy" json:"fee_invest_buy"`
		FeeInvestSell int64     `gorm:"column:fee_invest_sell" json:"fee_invest_sell"`
		Amount        int64     `gorm:"column:amount" json:"amount"`
	}
)

func (WalletEntity) TableName() string {
	return "tbl_wallets"
}
