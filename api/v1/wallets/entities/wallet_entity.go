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
		ID               uuid.UUID `gorm:"column:id"`
		IDAccount        uuid.UUID `gorm:"id_account"`
		InvestType       string    `gorm:"column:invest_type"`
		InvestName       string    `gorm:"column:invest_name"`
		InvestInstrument string    `gorm:"column:invest_instrument"`
		Active           bool      `gorm:"column:active"`
		FeeInvestBuy     int64     `gorm:"column:fee_invest_buy"`
		FeeInvestSell    int64     `gorm:"column:fee_invest_sell"`
		Amount           int64     `gorm:"column:amount"`
	}
)

func (WalletEntity) TableName() string {
	return "tbl_wallets"
}
