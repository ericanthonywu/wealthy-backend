package entities

import "github.com/google/uuid"

type (
	TransactionType struct {
		ID   uuid.UUID `json:"id" gorm:"column:id"`
		Type string    `json:"type" gorm:"column:type"`
	}

	IncomeType struct {
		ID         uuid.UUID `json:"id" gorm:"column:id"`
		IncomeType string    `json:"income_type" gorm:"column:income_types"`
	}

	ExpenseType struct {
		ID          uuid.UUID `json:"id" gorm:"column:id"`
		ExpenseType string    `json:"expense_type" gorm:"column:expense_types"`
	}

	ReksadanaType struct {
		ID   uuid.UUID `json:"id" gorm:"column:id"`
		Type string    `json:"type" gorm:"column:types"`
	}

	WalletType struct {
		ID         uuid.UUID `json:"" gorm:"column:id"`
		WalletType string    `json:"" gorm:"column:wallet_type"`
	}
)

func (TransactionType) TableName() string {
	return "tbl_master_transaction_types"
}

func (IncomeType) TableName() string {
	return "tbl_master_income_categories"
}

func (ExpenseType) TableName() string {
	return "tbl_master_expense_categories"
}

func (ReksadanaType) TableName() string {
	return "tbl_master_reksadana_types"
}

func (WalletType) TableName() string {
	return "tbl_master_wallet_types"
}
