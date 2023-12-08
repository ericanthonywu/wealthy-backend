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
		Name            string  `json:"name"`
		StockCode       string  `json:"stock_code"`
		Lot             int64   `json:"lot"`
		AverageBuy      float64 `json:"average_buy"`
		PotentialReturn float64 `json:"potential_return"`
	}
)