package dtos

type (
	InvestmentResponse struct {
		TotalInvestment           float64             `json:"total_investment"`
		TotalPotentialReturn      float64             `json:"total_potential_return"`
		PercentagePotentialReturn string              `json:"percentage_potential_return"`
		Details                   []InvestmentDetails `json:"details"`
	}

	InvestmentDetails struct {
		BrokerName string           `json:"broker_name"`
		Info       []InvestmentInfo `json:"info"`
	}

	InvestmentInfo struct {
		Name                string  `json:"name"`
		InitialInvestment   float64 `json:"initial_investment"`
		StockCode           string  `json:"stock_code"`
		Lot                 int64   `json:"lot"`
		AverageBuy          float64 `json:"average_buy"`
		PotentialReturn     float64 `json:"potential_return"`
		PercentageReturn    string  `json:"percentage_potential_return"`
		UnrealizedPotential float64 `json:"unrealized_potential"`
		TotalDays           int64   `json:"total_days"`
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
		Percentage        string  `json:"percentage"`
		TotalDays         int     `json:"total_days"`
	}
)