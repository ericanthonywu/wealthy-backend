package budgets

import (
	"github.com/google/uuid"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/budgets/entities"
	"gorm.io/gorm"
)

type (
	BudgetRepository struct {
		db *gorm.DB
	}

	IBudgetRepository interface {
		AllCategories(idPersonal uuid.UUID, month, year string) (budgetCategories []entities.BudgetAllCategoriesEntities)
		Set()
	}
)

func NewBudgetRepository(db *gorm.DB) *BudgetRepository {
	return &BudgetRepository{db: db}
}

func (r *BudgetRepository) AllCategories(idPersonal uuid.UUID, month, year string) (budgetCategories []entities.BudgetAllCategoriesEntities) {

	if err := r.db.Raw(`SELECT tmec.id,
       tmec.expense_types                                                          as categories,
       COALESCE(SUM(budget.amount)
                FILTER ( WHERE budget.id_personal_accounts = '?' AND
                               to_char(budget.created_at, 'MM') = '?' AND to_char(budget.created_at, 'YYYY') = '?'),
                0)                                                                 as total,
       (SELECT json_agg(r)::jsonb as sub_categories
        FROM (SELECT tmes.subcategories::text as subcategory_name,
                     COALESCE(SUM(b.amount::numeric)
                              FILTER ( WHERE b.id_personal_accounts = '?' AND
                                             to_char(b.created_at, 'MM') = '?' AND
                                             to_char(b.created_at, 'YYYY') = '?' ),
                              0)              as limit_amount
              FROM tbl_master_expense_subcategories tmes
                       LEFT JOIN tbl_budgets b ON b.id_master_subcategories = tmes.id
              WHERE tmes.active = TRUE
                AND tmes.id_master_expense_categories = tmec.id
              GROUP BY b.id_master_subcategories, tmes.subcategories, b.amount) r) as sub_categories
FROM tbl_master_expense_categories tmec
         LEFT JOIN tbl_master_expense_subcategories tmes ON tmec.id = tmes.id_master_expense_categories
         LEFT JOIN tbl_budgets budget ON budget.id_master_subcategories = tmes.id
WHERE tmec.active = true
GROUP BY tmec.expense_types, tmec.id`, idPersonal, month, year, idPersonal, month, year).Scan(&budgetCategories).Error; err != nil {
		return []entities.BudgetAllCategoriesEntities{}
	}
	return budgetCategories
}

func (r *BudgetRepository) Set() {

}
