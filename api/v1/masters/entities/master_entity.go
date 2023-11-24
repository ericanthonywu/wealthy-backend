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
		ID         uuid.UUID `json:"id" gorm:"column:id"`
		WalletType string    `json:"wallet" gorm:"column:wallet_type"`
	}

	InvestType struct {
		ID         uuid.UUID `json:"id" gorm:"column:id"`
		InvestName string    `json:"invest_name" gorm:"column:invest_name"`
	}

	Broker struct {
		ID         uuid.UUID `json:"id" gorm:"column:id"`
		BrokerName string    `json:"broker_name" gorm:"column:broker_name"`
	}

	TransactionPriority struct {
		ID       uuid.UUID `json:"id" gorm:"column:id"`
		Priority string    `json:"priority" gorm:"column:priority"`
	}

	Gender struct {
		ID         uuid.UUID `json:"id" gorm:"column:id"`
		GenderName string    `json:"gender_name" gorm:"column:gender_name"`
	}

	SubExpenseCategories struct {
		ID         uuid.UUID `json:"id" gorm:"column:id"`
		GenderName string    `json:"subcategories" gorm:"column:subcategories"`
	}

	Exchange struct {
		ID       uuid.UUID `gorm:"column:id" json:"id"`
		Currency string    `gorm:"column:currency" json:"currency"`
		Value    int64     `gorm:"column:value" json:"value"`
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

func (InvestType) TableName() string {
	return "tbl_master_invest"
}

func (Broker) TableName() string {
	return "tbl_master_broker"
}

func (TransactionPriority) TableName() string {
	return "tbl_master_transaction_priorities"
}

func (Gender) TableName() string {
	return "tbl_master_genders"
}

func (SubExpenseCategories) TableName() string {
	return "tbl_master_expense_subcategories"
}