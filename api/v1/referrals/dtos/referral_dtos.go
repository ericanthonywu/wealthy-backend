package dtos

type (
	ReferralResponse struct {
		Tier []TierDetail `json:"tier_details"`
	}

	ReferralResponseWithCustomer struct {
		Tier []TierDetailWithCustomer `json:"tier_details"`
	}

	TierDetail struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	TierDetailWithCustomer struct {
		Name           string           `json:"name"`
		Value          int              `json:"value"`
		CustomerDetail []CustomerDetail `json:"customer_detail"`
	}

	CustomerDetail struct {
		Name        string `json:"name,omitempty"`
		AccountType string `json:"account_type,omitempty"`
	}
)