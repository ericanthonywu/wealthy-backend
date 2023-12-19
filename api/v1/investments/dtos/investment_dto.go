package dtos

type (
	InvestmentResponse struct {
		Details []InvestmentDetails `json:"details"`
	}

	InvestmentDetails struct {
		BokerName           string           `json:"boker_name"`
		UnrealizedPotential float64          `json:"unrealized_potential"`
		Info                []InvestmentInfo `json:"info"`
	}

	InvestmentInfo struct {
		Name              string  `json:"name"`
		InitialInvestment float64 `json:"initial_investment"`
		StockCode         string  `json:"stock_code"`
		Lot               int64   `json:"lot"`
		AverageBuy        float64 `json:"average_buy"`
		PotentialReturn   float64 `json:"potential_return"`
	}

	InvestmentResponseGainLoss struct {
		Details []InvestmentGainLoss `json:"details"`
	}

	InvestmentGainLoss struct {
		DataTransaction   string  `json:"data_transaction"`
		BrokerName        string  `json:"broker_name"`
		InitialInvestment float64 `json:"initial_investment"`
		StockCode         string  `json:"investment_stock_code"`
		Name              string  `json:"investment_name"`
		Lot               int64   `json:"lot"`
		Price             float64 `json:"price"`
		GainLoss          float64 `json:"gain_loss"`
		Percentage        float64 `json:"percentage"`
		TotalDays         int     `json:"total_days"`
	}
)