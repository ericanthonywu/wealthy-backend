package entities

import "github.com/google/uuid"

type (
	DataPriceInfo struct {
		Price       float64 `json:"price"`
		AccountType string  `json:"account_type"`
		Description string  `json:"description"`
		PeriodName  string  `json:"period_name"`
	}

	SubsTransaction struct {
		ID                uuid.UUID `gorm:"column:id"`
		IDPersonalAccount uuid.UUID `gorm:"column:id_personal_accounts"`
		OrderID           string    `gorm:"column:order_id"`
		SubscriptionID    uuid.UUID `gorm:"column:subscription_id"`
		Amount            float64   `gorm:"column:amount"`
		Token             string    `gorm:"column:token"`
		RedirectURL       string    `gorm:"column:redirect_url"`
	}
)

func (SubsTransaction) TableName() string {
	return "tbl_subscriptions_transaction"
}