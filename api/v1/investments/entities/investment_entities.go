package entities

import (
	"github.com/google/uuid"
	"time"
)

type (
	InvestmentTransaction struct {
		Price           int64   `gorm:"column:price"`
		Lot             int64   `gorm:"column:lot"`
		StockCode       string  `gorm:"column:stock_code"`
		BrokerName      string  `gorm:"column:broker_name"`
		SellBuy         int64   `gorm:"column:sellbuy"`
		FeeBuy          float64 `gorm:"column:fee_buy"`
		FeeSell         float64 `gorm:"column:fee_sell"`
		DateTransaction string  `gorm:"column:date_transaction"`
		WalletName      string  `gorm:"column:wallet_name"`
	}

	InvestmentDataHelper struct {
		StockCode         string    `gorm:"column:stock_code"`
		TotalLot          int64     `gorm:"column:total_lot"`
		ValueBuy          float64   `gorm:"column:value_buy"`
		AverageBuy        float64   `gorm:"column:average_buy"`
		InitialInvestment float64   `gorm:"column:initial_investment"`
		IDMasterBroker    uuid.UUID `gorm:"column:id_master_broker"`
		GainLoss          float64   `gorm:"column:gain_loss"`
		PotentialReturn   float64   `gorm:"column:potential_return"`
		PercentageReturn  float64   `gorm:"column:percentage_return"`
		DateTransaction   time.Time `gorm:"column:created_at"`
		FeeBuy            float64   `gorm:"column:fee_buy"`
		NetBuy            float64   `gorm:"column:net_buy"`
	}

	InvestmentDataHelperPortfolio struct {
		StockCode         string    `gorm:"column:stock_code"`
		TotalLot          int64     `gorm:"column:total_lot"`
		ValueBuy          float64   `gorm:"column:value_buy"`
		AverageBuy        float64   `gorm:"column:average_buy"`
		AverageBuyOrigin  float64   `gorm:"column:average_buy_origin"`
		InitialInvestment float64   `gorm:"column:initial_investment"`
		IDMasterBroker    uuid.UUID `gorm:"column:id_master_broker"`
		GainLoss          float64   `gorm:"column:gain_loss"`
		PotentialReturn   float64   `gorm:"column:potential_return"`
		PercentageReturn  float64   `gorm:"column:percentage_return"`
		DateTransaction   time.Time `gorm:"column:created_at"`
		FeeBuy            float64   `gorm:"column:fee_buy"`
		NetBuy            float64   `gorm:"column:net_buy"`
		BrokerName        string    `gorm:"column:broker_name"`
		WalletID          uuid.UUID `gorm:"column:wallet_id"`
		WalletName        string    `gorm:"column:wallet_name"`
	}

	BrokerInfo struct {
		ID   uuid.UUID `gorm:"column:id"`
		Name string    `gorm:"column:broker_name"`
	}

	InvestmentTreding struct {
		Name   string `gorm:"column:name"`
		Symbol string `gorm:"column:symbol"`
		Close  int64  `gorm:"column:close"`
	}
)