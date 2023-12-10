package dtos

type (
	PaymentSubscription struct {
		PackageID string `json:"package_id"`
	}

	PaymentSnapRequest struct {
		Details PaymentSnapDetails `json:"transaction_details"`
	}

	PaymentSnapDetails struct {
		OrderId     string  `json:"order_id"`
		GrossAmount float64 `json:"gross_amount"`
	}

	MidTansResponse struct {
		Token         string   `json:"token,omitempty"`
		RedirectUrl   string   `json:"redirect_url,omitempty"`
		ErrorMessages []string `json:"error_messages,omitempty"`
	}

	MidTransWebhook struct {
		VaNumbers []struct {
			VaNumber string `json:"va_number"`
			Bank     string `json:"bank"`
		} `json:"va_numbers"`
		TransactionTime   string        `json:"transaction_time"`
		TransactionStatus string        `json:"transaction_status"`
		TransactionId     string        `json:"transaction_id"`
		StatusMessage     string        `json:"status_message"`
		StatusCode        string        `json:"status_code"`
		SignatureKey      string        `json:"signature_key"`
		SettlementTime    string        `json:"settlement_time"`
		PaymentType       string        `json:"payment_type"`
		PaymentAmounts    []interface{} `json:"payment_amounts"`
		OrderId           string        `json:"order_id"`
		MerchantId        string        `json:"merchant_id"`
		GrossAmount       string        `json:"gross_amount"`
		FraudStatus       string        `json:"fraud_status"`
		ExpiryTime        string        `json:"expiry_time"`
		Currency          string        `json:"currency"`
	}
)