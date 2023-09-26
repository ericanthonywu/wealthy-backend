package dtos

type (
	WalletAddRequest struct {
		InvestType       string `json:"invest_type"`
		InvestName       string `json:"invest_name"`
		InvestInstrument string `json:"invest_instrument"`
		Amount           int64  `json:"amount"`
		FeeInvestBuy     int64  `json:"fee_invest_buy"`
		FeeInvestSell    int64  `json:"fee_invest_sell"`
	}

	WalletAddResponse struct {
		InvestType       string `json:"invest_type,omitempty"`
		InvestName       string `json:"invest_name,omitempty"`
		InvestInstrument string `json:"invest_instrument,omitempty"`
		Amount           int64  `json:"amount,omitempty"`
	}
)
