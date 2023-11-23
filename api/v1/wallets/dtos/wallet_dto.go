package dtos

import "github.com/google/uuid"

type (
	WalletAddRequest struct {
		InvestType    string       `json:"invest_type"`
		InvestName    string       `json:"invest_name"`
		WalletType    string       `json:"id_master_wallet_type"`
		WalletAmount  WalletAmount `json:"wallet_amount"`
		FeeInvestBuy  int64        `json:"fee_invest_buy"`
		FeeInvestSell int64        `json:"fee_invest_sell"`
		Amount        int64        `json:"amount"`
	}

	WalletAmount struct {
		CurrencyCode string  `json:"currency_code"`
		Value        float64 `json:"value"`
	}

	WalletAddResponse struct {
		InvestType string `json:"invest_type,omitempty"`
		InvestName string `json:"invest_name,omitempty"`
		WalletType string `json:"id_master_wallet_type,omitempty"`
		Amount     int64  `json:"amount,omitempty"`
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