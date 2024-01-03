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
		ImagePath  string    `json:"image_path" gorm:"column:image_path"`
	}

	ExpenseType struct {
		ID          uuid.UUID `json:"id" gorm:"column:id"`
		ExpenseType string    `json:"expense_type" gorm:"column:expense_types"`
		ImagePath   string    `json:"image_path" gorm:"column:image_path"`
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

	StockCode struct {
		StockCode string `json:"stock_code" gorm:"column:symbol"`
		Name      string `json:"name" gorm:"column:name"`
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
		ID            uuid.UUID `json:"id" gorm:"column:id"`
		SubCategories string    `json:"sub_categories" gorm:"column:subcategories"`
		ImagePath     string    `json:"image_path" gorm:"column:image_path"`
	}

	Exchange struct {
		ID       uuid.UUID `gorm:"column:id" json:"id"`
		Currency string    `gorm:"column:currency" json:"currency"`
		Value    float64   `gorm:"column:value" json:"value"`
	}

	ExpenseCategoryEditable struct {
		ID          uuid.UUID `gorm:"column:id" json:"id"`
		ExpenseType string    `gorm:"column:expense_types" json:"category"`
		ImagePath   string    `gorm:"column:image_path" json:"image_path"`
	}

	ExpenseSubCategoryEditable struct {
		ID            uuid.UUID `gorm:"column:id" json:"id"`
		Subcategories string    `gorm:"column:subcategories" json:"category"`
		ImagePath     string    `gorm:"column:image_path" json:"image_path"`
	}

	IncomeCategoryEditable struct {
		ID         uuid.UUID `gorm:"column:id" json:"id"`
		IncomeType string    `gorm:"column:category" json:"income_type"`
		ImagePath  string    `gorm:"column:image_path" json:"image_path"`
	}

	AddEntities struct {
		ID uuid.UUID `gorm:"column:id" json:"id"`
	}

	Price struct {
		ID            uuid.UUID `gorm:"column:id"`
		Title         string    `gorm:"column:title"`
		Price         float64   `gorm:"column:price"`
		ActualPrice   float64   `gorm:"column:actual_price"`
		Description   string    `gorm:"column:description"`
		IsRecommended bool      `gorm:"column:is_recommended"`
	}

	SubscriptionInfo struct {
		ID uuid.UUID `gorm:"column:subscription_id"`
	}

	IncomeCategory struct {
		CategoryName string `gorm:"income_types"`
	}

	ExpenseCategory struct {
		CategoryName string `gorm:"expense_types"`
	}

	SubExpenseCategory struct {
		CategoryName string `gorm:"expense_types"`
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

func (Price) TableName() string {
	return "tbl_master_price"
}

func (IncomeCategory) TableName() string {
	return "tbl_master_income_categories_editable"
}

func (ExpenseCategory) TableName() string {
	return "tbl_master_income_categories_editable"
}

func (SubExpenseCategory) TableName() string {
	return "tbl_master_expense_subcategories_editable"
}

func (StockCode) TableName() string {
	return "tbl_master_trading"
}
