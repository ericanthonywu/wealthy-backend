package masters

import (
	"github.com/google/uuid"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/masters/entities"
	"github.com/sirupsen/logrus"
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
		Gender() (data []entities.Gender)
		SubExpenseCategory(expenseID uuid.UUID) (data []entities.SubExpenseCategories)
		ExpenseIDExist(expenseID uuid.UUID) (exist bool)
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

func (r *MasterRepository) Gender() (data []entities.Gender) {
	r.db.Find(&data).Scan(&data)
	return data
}

func (r *MasterRepository) SubExpenseCategory(expenseID uuid.UUID) (data []entities.SubExpenseCategories) {
	r.db.Where("id_master_expense_categories = ?", expenseID).Find(&data)
	return data
}

func (r *MasterRepository) ExpenseIDExist(expenseID uuid.UUID) (exist bool) {
	var model entities.ExpenseType

	err := r.db.First(&model, "id = ?", expenseID).Error
	if err != nil {
		logrus.Error(err.Error())
	}

	if model.ID != uuid.Nil {
		exist = true
	}

	return exist
}