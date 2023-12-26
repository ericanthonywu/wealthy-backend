package masters

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/wealthy-app/wealthy-backend/api/v1/masters/entities"
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
		Exchange() (data []entities.Exchange, err error)
		PersonalIncomeCategory(IDPersonal uuid.UUID) (data []entities.IncomeCategoryEditable, err error)
		PersonalExpenseCategory(IDPersonal uuid.UUID) (data []entities.ExpenseCategoryEditable, err error)
		PersonalExpenseSubCategory(IDPersonal uuid.UUID, expenseIDUUID uuid.UUID) (data []entities.ExpenseSubCategoryEditable, err error)
		RenameIncomeCategory(newName string, id, IDPersonal uuid.UUID) (err error)
		RenameExpenseCategory(newName string, id, IDPersonal uuid.UUID) (err error)
		RenameSubExpenseCategory(newName string, id, IDPersonal uuid.UUID) (err error)
		AddIncomeCategory(newCategory string, IDPersonal uuid.UUID) (data entities.AddEntities, err error)
		AddExpenseCategory(newCategory string, IDPersonal uuid.UUID) (data entities.AddEntities, err error)
		AddSubExpenseCategory(newCategory string, ExpenseID uuid.UUID, IDPersonal uuid.UUID) (data entities.AddEntities, err error)
		Price() (data []entities.Price)
		UserSubscriptionInfo(IDAccount uuid.UUID) (data entities.SubscriptionInfo, err error)
		GetAllIncomeCategories(accountUUID uuid.UUID) (data []string, err error)
		GetAllIExpenseCategories(accountUUID uuid.UUID) (data []string, err error)
		GetAllISubExpenseCategories(accountUUID uuid.UUID) (data []string, err error)
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
	r.db.Where("id_master_expense_categories = ?", expenseID).
		Where("active=?", true).Find(&data)
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

func (r *MasterRepository) Exchange() (data []entities.Exchange, err error) {
	if err := r.db.Raw(`SELECT tmec.id, tmec.currency_name as currency, currency_value as value FROM tbl_master_exchange_currency tmec
WHERE tmec.active = true`).Scan(&data).Error; err != nil {
		logrus.Error(err.Error())
		return []entities.Exchange{}, err
	}
	return data, nil
}

func (r *MasterRepository) PersonalIncomeCategory(IDPersonal uuid.UUID) (data []entities.IncomeCategoryEditable, err error) {
	if err = r.db.Raw(`SELECT tmice.id, tmice.income_types as category, tmice.image_path FROM tbl_master_income_categories_editable tmice
WHERE tmice.id_personal_accounts=? AND tmice.active=true`, IDPersonal).Scan(&data).Error; err != nil {
		return []entities.IncomeCategoryEditable{}, err
	}
	return data, nil
}

func (r *MasterRepository) PersonalExpenseCategory(IDPersonal uuid.UUID) (data []entities.ExpenseCategoryEditable, err error) {
	if err = r.db.Raw(`SELECT tmece.id, tmece.expense_types, tmece.image_path FROM tbl_master_expense_categories_editable tmece
WHERE tmece.id_personal_accounts=? AND
tmece.active=true;`, IDPersonal).Scan(&data).Error; err != nil {
		return []entities.ExpenseCategoryEditable{}, err
	}
	return data, nil
}

func (r *MasterRepository) PersonalExpenseSubCategory(IDPersonal, expenseIDUUID uuid.UUID) (data []entities.ExpenseSubCategoryEditable, err error) {
	if err = r.db.Raw(`SELECT tmese.id, tmese.subcategories, tmese.image_path , tmese.id_master_expense_categories FROM tbl_master_expense_subcategories_editable tmese
WHERE tmese.id_personal_accounts=?  AND tmese.id_master_expense_categories=? AND
tmese.active=true`, IDPersonal, expenseIDUUID).Scan(&data).Error; err != nil {
		return []entities.ExpenseSubCategoryEditable{}, err
	}
	return data, nil
}

func (r *MasterRepository) RenameIncomeCategory(newName string, id, IDPersonal uuid.UUID) (err error) {
	var model interface{}

	if err = r.db.Raw(`UPDATE tbl_master_income_categories_editable SET income_types=? WHERE id=? AND id_personal_accounts=?`, newName, id, IDPersonal).Scan(&model).
		Error; err != nil {
		logrus.Error(err.Error())
		return err
	}
	return nil
}

func (r *MasterRepository) RenameExpenseCategory(newName string, id, IDPersonal uuid.UUID) (err error) {
	var model interface{}

	if err = r.db.Raw(`UPDATE tbl_master_expense_categories_editable SET expense_types=? WHERE id=? AND id_personal_accounts=?`, newName, id, IDPersonal).Scan(&model).Error; err != nil {
		logrus.Error(err.Error())
		return err
	}
	return nil
}

func (r *MasterRepository) RenameSubExpenseCategory(newName string, id, IDPersonal uuid.UUID) (err error) {
	var model interface{}

	if err = r.db.Raw(`UPDATE tbl_master_expense_subcategories_editable SET subcategories=? WHERE id=? AND id_personal_accounts=?`, newName, id, IDPersonal).Scan(&model).
		Error; err != nil {
		logrus.Error(err.Error())
		return err
	}
	return nil
}

func (r *MasterRepository) AddIncomeCategory(newCategory string, IDPersonal uuid.UUID) (data entities.AddEntities, err error) {
	id, _ := uuid.NewUUID()

	if err = r.db.Raw(`INSERT INTO tbl_master_income_categories_editable (id, income_types, active, id_personal_accounts) VALUES (?,?, ?, ?) RETURNING id`, id, newCategory, true, IDPersonal).Scan(&data).Error; err != nil {
		logrus.Error(err.Error())
		return entities.AddEntities{}, err
	}
	return data, nil
}

func (r *MasterRepository) AddExpenseCategory(newCategory string, IDPersonal uuid.UUID) (data entities.AddEntities, err error) {
	id, _ := uuid.NewUUID()

	if err = r.db.Raw(`INSERT INTO tbl_master_expense_categories_editable (id,expense_types,active, id_personal_accounts) VALUES (?,?, ?, ?) RETURNING id`, id, newCategory, true, IDPersonal).Scan(&data).Error; err != nil {
		logrus.Error(err.Error())
		return entities.AddEntities{}, err
	}
	return data, nil
}

func (r *MasterRepository) AddSubExpenseCategory(newCategory string, ExpenseID uuid.UUID, IDPersonal uuid.UUID) (data entities.AddEntities, err error) {
	id, _ := uuid.NewUUID()

	if err = r.db.Raw(`INSERT INTO tbl_master_expense_subcategories_editable (id, subcategories,id_master_expense_categories, active, id_personal_accounts ) VALUES (?,?, ?, ?, ?) RETURNING id`, id, newCategory, ExpenseID, true, IDPersonal).Scan(&data).Error; err != nil {
		logrus.Error(err.Error())
		return entities.AddEntities{}, err
	}
	return data, nil
}

func (r *MasterRepository) Price() (data []entities.Price) {
	if err := r.db.Find(&data).Error; err != nil {
		logrus.Error(err.Error())
		return []entities.Price{}
	}
	return data
}

func (r *MasterRepository) UserSubscriptionInfo(IDAccount uuid.UUID) (data entities.SubscriptionInfo, err error) {
	if err = r.db.Raw(`SELECT tst.subscription_id FROM tbl_user_subscription tus
INNER JOIN tbl_subscriptions_transaction tst ON tst.id = tus.id_subscriptions_transaction
WHERE tus.id_personal_accounts = ?
ORDER BY tus.period_expired DESC LIMIT 1`, IDAccount).Scan(&data).Error; err != nil {
		return entities.SubscriptionInfo{}, err
	}
	return data, nil
}

func (r *MasterRepository) GetAllIncomeCategories(accountUUID uuid.UUID) (data []string, err error) {
	if err := r.db.Table("tbl_master_income_categories_editable").
		Distinct("income_types").
		Where("id_personal_accounts=?", accountUUID).
		Where("active=?", true).
		Pluck("income_types", &data).Error; err != nil {
		logrus.Error(err.Error())
		return []string{}, err
	}
	return data, nil
}

func (r *MasterRepository) GetAllIExpenseCategories(accountUUID uuid.UUID) (data []string, err error) {
	if err := r.db.Table("tbl_master_expense_categories_editable").
		Distinct("expense_types").
		Where("id_personal_accounts=?", accountUUID).
		Where("active=?", true).
		Pluck("expense_types", &data).Error; err != nil {
		logrus.Error(err.Error())
		return []string{}, err
	}
	return data, nil
}

func (r *MasterRepository) GetAllISubExpenseCategories(accountUUID uuid.UUID) (data []string, err error) {
	if err := r.db.Table("tbl_master_expense_subcategories_editable").
		Distinct("subcategories").
		Where("id_personal_accounts=?", accountUUID).
		Where("active=?", true).
		Pluck("subcategories", &data).Error; err != nil {
		logrus.Error(err.Error())
		return []string{}, err
	}
	return data, nil
}