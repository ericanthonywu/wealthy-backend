package internals

import (
	"github.com/google/uuid"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/transactions/entities"
	"gorm.io/gorm"
)

type (
	InternalRepository struct {
		db *gorm.DB
	}

	IInternalRepository interface {
		ByNote(customerID uuid.UUID, month, year string) (data []entities.TransactionByNotes)
	}
)

func NewInternalRepository(db *gorm.DB) *InternalRepository {
	return &InternalRepository{db: db}
}

func (r *InternalRepository) ByNote(customerID uuid.UUID, month, year string) (data []entities.TransactionByNotes) {
	if err := r.db.Raw(`SELECT to_char(t.date_time_transaction::DATE, 'MM')           as month,
       concat(to_char(to_date(t.date_time_transaction, 'YYYY-MM-DD'), 'Mon'), ' ',
              to_char(t.date_time_transaction::DATE, 'YYYY'))::text as month_year,
       count(td.note)::text                                         as quantity,
       td.note::text                                                as transaction_note,
       tmec.expense_types ::text                                    as transaction_category
FROM tbl_transactions t
         INNER JOIN tbl_transaction_details td ON td.id_transactions = t.id
         INNER JOIN tbl_master_transaction_types tmtt ON tmtt.id = t.id_master_transaction_types
         INNER JOIN tbl_master_expense_categories tmec
                    on t.id_master_expense_categories = tmec.id
WHERE to_char(t.date_time_transaction::DATE, 'MM') = ?
  AND to_char(t.date_time_transaction::DATE, 'YYYY') = ?
  AND t.id_personal_account = ?
  AND tmtt.type = 'EXPENSE'
GROUP BY month_year, td.note, expense_types, month
ORDER BY month DESC, month_year`, month, year, customerID).Scan(&data).Error; err != nil {
		return []entities.TransactionByNotes{}
	}
	return data
}