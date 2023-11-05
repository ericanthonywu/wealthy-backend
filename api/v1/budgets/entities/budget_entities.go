package entities

import "github.com/google/uuid"

type (
	BudgetAllCategoriesEntities struct {
		ID            uuid.UUID `gorm:"column:category_id"`
		Categories    string    `gorm:"column:category_name"`
		Total         string    `gorm:"column:budget_amount"`
		SubCategories string    `gorm:"column:sub_categories"`
	}

	SubCategoryBudget struct {
		CategoryID      uuid.UUID `gorm:"column:category_id"`
		CategoryName    string    `gorm:"column:category_name"`
		SubCategoryID   uuid.UUID `gorm:"column:sub_category_id"`
		SubCategoryName string    `gorm:"column:sub_category_name"`
		BudgetLimit     int       `gorm:"column:budget_limit"`
	}

	BudgetTotalSpendingAndNumberOfCategory struct {
		ID               uuid.UUID `gorm:"column:id"`
		Category         string    `gorm:"column:category"`
		Spending         int       `gorm:"column:spending"`
		NumberOfCategory int       `gorm:"column:number_of_category"`
	}

	BudgetLimit struct {
		IDMasterExpense uuid.UUID `gorm:"column:id_master_expense"`
		BudgetLimit     int       `gorm:"column:budget_limit"`
		ExpenseType     string    `gorm:"column:expense_types"`
	}

	BudgetCategory struct {
		TransactionCategory string `gorm:"column:transaction_category" json:"transaction_category"`
		BudgetLimit         int    `gorm:"column:budget_limit" json:"budget_limit"`
		TotalSpending       int    `gorm:"column:total_spending" json:"total_spending"`
		TotalRemaining      int    `gorm:"column:total_remaining" json:"total_remaining"`
	}

	BudgetLatestSixMonth struct {
		Period        string `gorm:"column:period" json:"period"`
		TotalSpending int    `gorm:"column:total_spending" json:"total_spending"`
		BudgetLimit   int    `gorm:"column:budget_limit" json:"budget_limit"`
		Percentage    string `gorm:"column:percentage" json:"percentage"`
	}

	BudgetSetEntities struct {
		ID                uuid.UUID `gorm:"column:id"`
		IDPersonalAccount uuid.UUID `gorm:"column:id_personal_accounts"`
		IDCategory        uuid.UUID `gorm:"column:id_master_categories"`
		IDSubCategory     uuid.UUID `gorm:"column:id_master_subcategories"`
		Amount            int       `gorm:"column:amount"`
	}

	BudgetExistEntities struct {
		ID uuid.UUID `gorm:"column:id"`
	}

	PersonalBudget struct {
		ID          uuid.UUID `gorm:"column:id"`
		Category    string    `gorm:"column:category"`
		BudgetLimit int       `gorm:"column:budget"`
	}

	PersonalTransaction struct {
		ID       uuid.UUID `gorm:"column:id"`
		Category string    `gorm:"column:category"`
		Amount   int       `gorm:"column:amount"`
		Count    int       `gorm:"column:count"`
	}
)

func (BudgetSetEntities) TableName() string {
	return "tbl_budgets"
}