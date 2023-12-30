package models

import "github.com/google/uuid"

type (
	AllPreviousTrxInvestment struct {
		Date           string    `gorm:"column:date_time_transaction"`
		IDMasterBroker uuid.UUID `gorm:"column:id_master_broker"`
		Amount         float64   `gorm:"column:amount"`
		IDTransaction  uuid.UUID `gorm:"column:id_transactions"`
		Lot            int64     `gorm:"column:lot"`
		SellBuy        int       `gorm:"column:sellbuy"`
		StockCode      string    `gorm:"column:stock_code"`
		FeeBuy         float64   `gorm:"column:fee_buy"`
		FeeSell        float64   `gorm:"column:fee_sell"`
		IDWallet       uuid.UUID `gorm:"column:wallet_id"`
	}

	TrxInvest []AllPreviousTrxInvestment
)

func (AllPreviousTrxInvestment) TableName() string {
	return "tbl_transaction_details"
}
func (a TrxInvest) Len() int      { return len(a) }
func (a TrxInvest) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a TrxInvest) Less(i, j int) bool {
	return a[i].IDMasterBroker.String() < a[j].IDMasterBroker.String() ||
		(a[i].IDMasterBroker.String() == a[j].IDMasterBroker.String() && a[i].StockCode < a[j].StockCode) ||
		(a[i].IDMasterBroker.String() == a[j].IDMasterBroker.String() && a[i].StockCode == a[j].StockCode && a[i].SellBuy > a[j].SellBuy)
}