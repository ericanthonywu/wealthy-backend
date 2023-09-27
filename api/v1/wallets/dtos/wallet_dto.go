package dtos

import "github.com/google/uuid"

type (
	WalletAddRequest struct {
		InvestType       string `json:"invest_type"`
		InvestName       string `json:"invest_name"`
		InvestInstrument string `json:"invest_instrument"`
		WalletType       string `json:"wallet_type"`
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

	WalletUpdateAmountRequest struct {
		Amount int64 `json:"amount"`
	}

	WalletUpdateAmountResponse struct {
		ID               uuid.UUID `json:"id,omitempty"`
		IDAccount        uuid.UUID `json:"id_account,omitempty"`
		InvestType       string    `json:"invest_type,omitempty"`
		InvestName       string    `json:"invest_name,omitempty"`
		InvestInstrument string    `json:"invest_instrument,omitempty"`
		Active           bool      `json:"active,omitempty"`
		FeeInvestBuy     int64     `json:"fee_invest_buy,omitempty"`
		FeeInvestSell    int64     `json:"fee_invest_sell,omitempty"`
		Amount           int64     `json:"amount,omitempty"`
	}
)
