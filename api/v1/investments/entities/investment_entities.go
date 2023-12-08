package entities

type (
	InvestmentTransaction struct {
		Price      int64   `gorm:"column:price"`
		Lot        int64   `gorm:"column:lot"`
		StockCode  string  `gorm:"column:stock_code"`
		BrokerName string  `gorm:"column:broker_name"`
		SellBuy    int64   `gorm:"column:sellbuy"`
		FeeBuy     float64 `gorm:"column:fee_buy"`
		FeeSell    float64 `gorm:"column:fee_sell"`
	}

	InvestmentTreding struct {
		Name   string `gorm:"column:name"`
		Symbol string `gorm:"column:symbol"`
		Close  int64  `gorm:"column:close"`
	}
)