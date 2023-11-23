package dtos

import "github.com/google/uuid"

type (
	PaymentSubscription struct {
		PackageID uuid.UUID `json:"package_id"`
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
)