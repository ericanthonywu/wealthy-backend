package entities

import "github.com/google/uuid"

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
		Email    string `gorm:"column:email"`
		Password string `gorm:"column:password"`
		Active   bool   `gorm:"column:active"`
		Roles    string `gorm:"column:role"`
	}
)

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
