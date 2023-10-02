package masters

import (
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/masters/entities"
	"gorm.io/gorm"
)

type (
	MasterRepository struct {
		db *gorm.DB
	}

	IMasterRepository interface {
		TransactionType() (data []entities.TransactionType)
		IncomeType() (data []entities.IncomeType)
		ExpenseType() (data []entities.ExpenseType)
		ReksadanaType() (data []entities.ReksadanaType)
		WalletType() (data []entities.WalletType)
		InvestType() (data []entities.InvestType)
		Broker() (data []entities.Broker)
		TransactionPriority() (data []entities.TransactionPriority)
	}
)

func NewMasterRepository(db *gorm.DB) *MasterRepository {
	return &MasterRepository{db: db}
}

func (r *MasterRepository) TransactionType() (data []entities.TransactionType) {
	r.db.Find(&data)
	return data
}

func (r *MasterRepository) IncomeType() (data []entities.IncomeType) {
	r.db.Where("active = ?", true).Find(&data)
	return data
}

func (r *MasterRepository) ExpenseType() (data []entities.ExpenseType) {
	r.db.Find(&data)
	return data
}

func (r *MasterRepository) ReksadanaType() (data []entities.ReksadanaType) {
	r.db.Where("active = ?", true).Find(&data)
	return data
}

func (r *MasterRepository) WalletType() (data []entities.WalletType) {
	r.db.Where("active=?", true).Find(&data)
	return data
}

func (r *MasterRepository) InvestType() (data []entities.InvestType) {
	r.db.Where("active=?", true).Find(&data)
	return data
}

func (r *MasterRepository) Broker() (data []entities.Broker) {
	r.db.Where("active=?", true).Find(&data)
	return data
}

func (r *MasterRepository) TransactionPriority() (data []entities.TransactionPriority) {
	r.db.Where("active=?", true).Find(&data)
	return
}
