package dtos

import "github.com/google/uuid"

type (
	WalletAddRequest struct {
		WalletType    string  `json:"wallet_type"`
		WalletName    string  `json:"wallet_name"`
		TotalAsset    int64   `json:"total_assets"`
		FeeInvestBuy  float64 `json:"fee_invest_buy"`
		FeeInvestSell float64 `json:"fee_invest_sell"`
	}

	WalletAddResponse struct {
		WalletID      uuid.UUID `json:"wallet_id"`
		WalletName    string    `json:"wallet_name"`
		WalletType    string    `json:"wallet_type"`
		FeeInvestBuy  float64   `json:"fee_invest_buy"`
		FeeInvestSell float64   `json:"fee_invest_sell"`
		TotalAssets   int64     `json:"total_assets"`
	}

	WalletListResponse struct {
		IDAccount     uuid.UUID     `json:"id_personal_accounts"`
		WalletDetails WalletDetails `json:"details"`
		Active        bool          `json:"is_active"`
		FeeInvestBuy  float64       `json:"fee_invest_buy"`
		FeeInvestSell float64       `json:"fee_invest_sell"`
		TotalAssets   int64         `json:"total_assets"`
	}

	WalletDetails struct {
		WalletID           uuid.UUID `json:"wallet_id"`
		WalletType         string    `json:"wallet_type"`
		WalletName         string    `json:"wallet_name"`
		IDMasterWalletType uuid.UUID `json:"id_master_wallet_types"`
	}

	WalletUpdateAmountRequest struct {
		Amount     int64  `json:"amount,omitempty"`
		WalletName string `json:"wallet_name"`
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