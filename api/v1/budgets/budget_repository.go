package budgets

import "gorm.io/gorm"

type (
	BudgetRepository struct {
		db *gorm.DB
	}

	IBudgetRepository interface {
		AllCategories()
		Set()
	}
)

func NewBudgetRepository(db *gorm.DB) *BudgetRepository {
	return &BudgetRepository{db: db}
}

func (r *BudgetRepository) AllCategories() {

}

func (r *BudgetRepository) Set() {

}
