package transactions

import (
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/transactions/entities"
	"gorm.io/gorm"
)

type (
	TransactionRepository struct {
		db *gorm.DB
	}

	ITransactionRepository interface {
		Add(trx *entities.TransactionEntity, trxDetail *entities.TransactionDetailEntity) (err error)
	}
)

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Add(trx *entities.TransactionEntity, trxDetail *entities.TransactionDetailEntity) (err error) {
	err = r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&trx).Error; err != nil {
			return err
		}

		if err := tx.Create(&trxDetail).Error; err != nil {
			return err
		}
		return nil
	})

	return err
}
