package transactions

import "gorm.io/gorm"

type (
	TransactionRepository struct {
		db *gorm.DB
	}

	ITransactionRepository interface {
	}
)

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}
