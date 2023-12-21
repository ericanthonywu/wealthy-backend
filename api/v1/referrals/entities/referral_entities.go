package entities

import (
	"github.com/google/uuid"
	"time"
)

type (
	ReferralAccountProfile struct {
		ID          uuid.UUID `gorm:"column:id" json:"id"`
		Email       string    `gorm:"column:email" json:"email"`
		Username    string    `gorm:"column:username" json:"username"`
		Name        string    `gorm:"column:name" json:"name"`
		DOB         string    `gorm:"column:date_of_birth" json:"date_of_birth"`
		ReferType   string    `gorm:"column:refer_code" json:"referral_code"`
		AccountType string    `gorm:"column:account_type" json:"account_type"`
		IDGender    uuid.UUID `gorm:"column:id_gender" json:"id_gender"`
		Gender      string    `gorm:"column:gender" json:"gender"`
		UserRoles   string    `gorm:"column:user_roles" json:"user_roles"`
		ImagePath   string    `gorm:"column:image_path" json:"image_path"`
		FileName    string    `gorm:"column:file_name" json:"file_name"`
	}

	ReferralAccountProfileRefCode struct {
		Username    string `gorm:"column:username"`
		Name        string `gorm:"column:name"`
		AccountType string `gorm:"column:account_type"`
	}

	ReferralUserReward struct {
		ID               uuid.UUID `gorm:"column:id"`
		RefCode          string    `gorm:"column:ref_code"`
		RefCodeReference string    `gorm:"column:ref_code_reference"`
		Level            int       `gorm:"column:level"`
		CreatedAt        time.Time `gorm:"column:created_at"`
	}

	PreviousCommission struct {
		Commission float64 `gorm:"column:total_comission"`
	}

	WithdrawEntities struct {
		ID                 uuid.UUID `gorm:"column:id"`
		IDPersonalAccounts uuid.UUID `gorm:"column:id_personal_accounts"`
		AccountNumber      int       `gorm:"column:account_number"`
		AccountName        string    `gorm:"column:account_name"`
		BankIssue          string    `gorm:"column:bank_issue"`
		Amount             float64   `gorm:"column:amount"`
		Status             int       `gorm:"column:status"`
	}
)

func (WithdrawEntities) TableName() string {
	return "tbl_withdraw_request_transaction"
}