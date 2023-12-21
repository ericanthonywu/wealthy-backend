package entities

import (
	"github.com/google/uuid"
	"time"
)

type (
	DataPriceInfo struct {
		Price              float64   `gorm:"column:price"`
		AccountType        string    `gorm:"column:account_type"`
		Description        string    `gorm:"column:description"`
		PeriodName         string    `gorm:"column:period_name"`
		IDMasterSubsPeriod uuid.UUID `json:"column:id_master_subs_period"`
	}

	SubsTransaction struct {
		ID                 uuid.UUID `gorm:"column:id"`
		IDPersonalAccount  uuid.UUID `gorm:"column:id_personal_accounts"`
		OrderID            string    `gorm:"column:order_id"`
		SubscriptionID     uuid.UUID `gorm:"column:subscription_id"`
		Amount             float64   `gorm:"column:amount"`
		Token              string    `gorm:"column:token"`
		RedirectURL        string    `gorm:"column:redirect_url"`
		IDMasterSubsPeriod uuid.UUID `gorm:"column:id_master_subs_period"`
	}

	SubsInfo struct {
		ID                 uuid.UUID `gorm:"column:id"`
		IDPersonalAccounts uuid.UUID `gorm:"column:id_personal_accounts"`
		IDSubsTransaction  uuid.UUID `gorm:"column:id_subscriptions_transaction"`
		PeriodEExpired     time.Time `gorm:"column:period_expired"`
		CreatedAt          time.Time `gorm:"created_at"`
	}

	CheckPackage struct {
		Exists bool `gorm:"column:exists"`
	}

	GetPeriodName struct {
		ID         uuid.UUID `gorm:"column:id"`
		PeriodName string    `gorm:"column:period_name"`
	}

	GetReferralInfo struct {
		RefCode          string `gorm:"column:ref_code"`
		RefCodeReference string `gorm:"column:ref_code_reference"`
		Level            int    `gorm:"column:level"`
	}

	RewardInfo struct {
		Level      int     `gorm:"column:level"`
		Percentage float64 `gorm:"column:percentage"`
	}

	PersonalInfo struct {
		RefCode string `gorm:"column:refer_code"`
	}

	PreviousCommission struct {
		Commission float64 `gorm:"column:total_comission"`
	}
)

func (SubsTransaction) TableName() string {
	return "tbl_subscriptions_transaction"
}

func (SubsInfo) TableName() string {
	return "tbl_user_subscription"
}