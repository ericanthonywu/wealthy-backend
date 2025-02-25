package entities

import (
	"github.com/google/uuid"
	"time"
)

type (
	AccountMasterRoles struct {
		ID    uuid.UUID `gorm:"column:id"`
		Roles string    `gorm:"column:roles"`
	}

	AccountMasterAccountType struct {
		ID          uuid.UUID `gorm:"column:id"`
		AccountType string    `gorm:"column:account_type"`
	}

	AccountMasterGender struct {
		ID uuid.UUID `gorm:"column:id"`
	}

	AccountSignUpPersonalAccountEntity struct {
		ID                  uuid.UUID `gorm:"column:id"`
		Username            string    `gorm:"column:username"`
		Name                string    `gorm:"column:name"`
		Email               string    `gorm:"column:email"`
		ReferCode           string    `gorm:"column:refer_code"`
		IDMasterAccountType uuid.UUID `gorm:"column:id_master_account_types"`
	}

	AccountSignUpAuthenticationsEntity struct {
		Password           string    `gorm:"column:password"`
		Active             bool      `gorm:"column:active"`
		IDMasterRoles      uuid.UUID `gorm:"column:id_master_roles"`
		IDPersonalAccounts uuid.UUID `gorm:"column:id_personal_accounts"`
	}

	AccountSignInAuthenticationEntity struct {
		ID       uuid.UUID `gorm:"column:id"`
		Email    string    `gorm:"column:email"`
		Password string    `gorm:"column:password"`
		Active   bool      `gorm:"column:active"`
		Roles    string    `gorm:"column:role"`
		Type     string    `gorm:"column:type"`
	}

	AccountProfile struct {
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
		Latitude    string    `gorm:"column:lat" json:"latitude"`
		Longitude   string    `gorm:"column:long" json:"longitude"`
	}

	AccountProfilePassword struct {
		ID          uuid.UUID `gorm:"column:id" json:"id"`
		Email       string    `gorm:"column:email" json:"email"`
		Username    string    `gorm:"column:username" json:"username"`
		Password    string    `gorm:"column:password" json:"password"`
		Name        string    `gorm:"column:name" json:"name"`
		DOB         string    `gorm:"column:date_of_birth" json:"date_of_birth"`
		ReferType   string    `gorm:"column:refer_code" json:"referral_code"`
		AccountType string    `gorm:"column:account_type" json:"account_type"`
		Gender      string    `gorm:"column:gender" json:"gender"`
		UserRoles   string    `gorm:"column:user_roles" json:"user_roles"`
	}

	AccountSetProfileEntity struct {
		ID       uuid.UUID `gorm:"column:id"`
		Name     string    `gorm:"column:name"`
		Username string    `gorm:"column:username"`
		DOB      string    `gorm:"column:dob"`
		Gender   uuid.UUID `gorm:"column:id_master_gender"`
	}

	AccountAuthorization struct {
		ID uuid.UUID `gorm:"column:id_personal_accounts"`
	}

	AccountRewards struct {
		ID               uuid.UUID `gorm:"column:id"`
		RefCode          string    `gorm:"column:ref_code"`
		RefCodeReference string    `gorm:"column:ref_code_reference"`
		Level            int       `gorm:"column:level"`
	}

	AccountSearchEmail struct {
		ID uuid.UUID `gorm:"column:id"`
	}

	AccountGroupSharing struct {
		ID         uuid.UUID `gorm:"column:id"`
		ShareFrom  uuid.UUID `gorm:"column:id_personal_accounts_share_from"`
		ShareTo    uuid.UUID `gorm:"column:id_personal_accounts_share_to"`
		IsAccepted bool      `gorm:"column:is_accepted"`
	}

	AccountPersonalIDGroupSharing struct {
		ID         uuid.UUID `gorm:"column:id_personal_accounts"`
		IsAccepted bool      `gorm:"column:is_accepted"`
		Name       string    `gorm:"column:name"`
	}

	AccountGroupSharingWithProfileInfo struct {
		ID          uuid.UUID `gorm:"column:id"`
		Name        string    `gorm:"column:name"`
		Email       string    `gorm:"column:email"`
		Filename    string    `gorm:"column:file_name"`
		ImagePath   string    `gorm:"column:image_path"`
		Status      string    `gorm:"column:status"`
		Type        string    `gorm:"column:type"`
		IDShareFrom uuid.UUID `gorm:"column:id_share_from"`
		IDShareTo   uuid.UUID `gorm:"column:id_share_to"`
	}

	AccountForgotPassword struct {
		ID                uuid.UUID `gorm:"column:id"`
		OTPCode           string    `gorm:"column:otp_code"`
		IDPersonalAccount uuid.UUID `gorm:"column:id_personal_accounts"`
		IsVerified        bool      `gorm:"column:is_verified"`
		Expired           time.Time `gorm:"column:expired"`
		CreatedAt         time.Time `gorm:"column:created_at"`
	}

	AccountGender struct {
		Exists bool `gorm:"column:exists"`
	}

	AccountAlreadySharing struct {
		Exists bool `gorm:"column:exists"`
	}

	AccountNotificationEntities struct {
		ID                      uuid.UUID `gorm:"column:id"`
		Name                    string    `gorm:"name"`
		NotificationTitle       string    `gorm:"column:notification_title"`
		NotificationDescription string    `gorm:"column:notification_description"`
		IDPersonalAccounts      uuid.UUID `gorm:"column:id_personal_accounts"`
		IsRead                  bool      `gorm:"column:is_read"`
		IDGroupSender           uuid.UUID `gorm:"column:id_group_sender"`
		IDGroupRecipient        uuid.UUID `gorm:"column:id_group_recipient"`
		ImagePath               string    `gorm:"column:image_path"`
		Type                    string    `gorm:"column:type"`
		CreatedAt               time.Time `gorm:"column:created_at"`
	}

	AccountWallet struct {
		ID uuid.UUID `gorm:"column:id"`
	}

	AccountTransaction struct {
		ID uuid.UUID `gorm:"primaryKey;column:id"`
	}

	AccountTransactionDetail struct {
		ID            uuid.UUID `gorm:"primaryKey;column:id"`
		IDTransaction uuid.UUID `gorm:"foreignKey:ID;column:id_transactions"`
	}

	AccountMasterExpenseCategory struct {
		IDPersonalAccounts uuid.UUID `gorm:"column:id_personal_accounts"`
	}

	AccountMasterSubExpenseCategory struct {
		IDPersonalAccounts uuid.UUID `gorm:"column:id_personal_accounts"`
	}

	AccountMasterIncomeCategory struct {
		IDPersonalAccounts uuid.UUID `gorm:"column:id_personal_accounts"`
	}

	AccountBudget struct {
		IDPersonalAccounts uuid.UUID `gorm:"column:id_personal_accounts"`
	}

	AccountGroupSharings struct {
		IDPersonalAccounts uuid.UUID `gorm:"column:id_personal_accounts_share_from"`
	}

	AccountSubscriptions struct {
		IDPersonalAccounts uuid.UUID `gorm:"column:id_personal_accounts"`
	}

	AccountWithdraw struct {
		IDPersonalAccounts uuid.UUID `gorm:"column:id_personal_accounts"`
	}

	AccountAuthorizations struct {
		IDPersonalAccounts uuid.UUID `gorm:"column:id_personal_accounts"`
	}

	AccountPersonal struct {
		ID uuid.UUID `gorm:"column:id"`
	}
)

func (AccountAuthorization) TableName() string {
	return "tbl_authentications"
}

func (AccountSetProfileEntity) TableName() string {
	return "tbl_personal_accounts"

}

func (AccountMasterAccountType) TableName() string {
	return "tbl_master_account_types"
}

func (AccountMasterRoles) TableName() string {
	return "tbl_master_roles"
}

func (AccountSignUpPersonalAccountEntity) TableName() string {
	return "tbl_personal_accounts"
}

func (AccountSignUpAuthenticationsEntity) TableName() string {
	return "tbl_authentications"
}

func (AccountRewards) TableName() string {
	return "tbl_user_rewards"
}

func (AccountGroupSharing) TableName() string {
	return "tbl_group_sharing"
}

func (AccountForgotPassword) TableName() string {
	return "tbl_forgot_password"
}

func (AccountWallet) TableName() string {
	return "tbl_wallets"
}

func (AccountMasterExpenseCategory) TableName() string {
	return "tbl_master_expense_categories_editable"
}

func (AccountMasterSubExpenseCategory) TableName() string {
	return "tbl_master_expense_subcategories_editable"
}

func (AccountMasterIncomeCategory) TableName() string {
	return "tbl_master_income_categories_editable"
}

func (AccountBudget) TableName() string {
	return "tbl_budgets"
}

func (AccountGroupSharings) TableName() string {
	return "tbl_group_sharing"
}

func (AccountSubscriptions) TableName() string {
	return "tbl_user_subscription"
}

func (AccountWithdraw) TableName() string {
	return "tbl_withdraw_request_transaction"
}

func (AccountAuthorizations) TableName() string {
	return "tbl_authentications"
}

func (AccountPersonal) TableName() string {
	return "tbl_personal_accounts"
}