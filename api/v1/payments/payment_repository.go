package payments

import (
	"github.com/google/uuid"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/payments/entities"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type (
	PaymentRepository struct {
		db *gorm.DB
	}

	IPaymentRepository interface {
		GetPrice(subs uuid.UUID) (data entities.DataPriceInfo)
		SaveSubscriptionPayment(model *entities.SubsTransaction) (result bool, err error)
		UpdateStatusTransaction(orderID string) (err error)
		GetSubscriptionInformation(IDPersonalAccount uuid.UUID) (data entities.SubsInfo, err error)
		CheckPackageID(IDPackage uuid.UUID) (data entities.CheckPackage)
		GetTransactionInfoByOrderID(OrderID string) (data entities.SubsTransaction, err error)
		GetPeriodName(ID uuid.UUID) (data entities.GetPeriodName, err error)
		WriteUserSubscription(model entities.SubsInfo) (err error)
		ChangeAccountUser(IDPersonalAccount, IDProAccountUUID uuid.UUID) (err error)
	}
)

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) GetPrice(subs uuid.UUID) (data entities.DataPriceInfo) {
	if err := r.db.Raw(`SELECT tmp.price, tmat.account_type, tmp.description, tmsp.period_name, tmp.id_master_subs_period FROM tbl_master_price tmp
         INNER JOIN tbl_master_account_types tmat ON tmp.id_master_account_types = tmat.id
         INNER JOIN tbl_master_subs_period tmsp ON tmp.id_master_subs_period = tmsp.id WHERE tmp.id = ? AND tmp.active=true`, subs).Scan(&data).Error; err != nil {
		logrus.Error(err.Error())
	}
	return data
}

func (r *PaymentRepository) SaveSubscriptionPayment(model *entities.SubsTransaction) (result bool, err error) {
	if err := r.db.Create(&model).Error; err != nil {
		logrus.Error(err.Error())
		return false, err
	}
	return true, nil
}

func (r *PaymentRepository) UpdateStatusTransaction(orderID string) (err error) {
	var model interface{}

	if err := r.db.Raw(`UPDATE tbl_subscriptions_transaction SET status=? WHERE order_id=?`, 1, orderID).Scan(&model).Error; err != nil {
		return err
	}

	return nil
}

func (r *PaymentRepository) GetSubscriptionInformation(IDPersonalAccount uuid.UUID) (data entities.SubsInfo, err error) {
	if err := r.db.Raw(`SELECT * FROM tbl_user_subscription WHERE id_personal_accounts=? ORDER BY created_at DESC LIMIT 1`, IDPersonalAccount).Scan(&data).Error; err != nil {
		return entities.SubsInfo{}, err
	}

	return data, nil
}

func (r *PaymentRepository) CheckPackageID(IDPackage uuid.UUID) (data entities.CheckPackage) {
	if err := r.db.Raw(`SELECT EXISTS( SELECT 1 FROM tbl_master_price WHERE id = ?)`, IDPackage).Scan(&data).Error; err != nil {
		return data
	}
	return data
}

func (r *PaymentRepository) GetTransactionInfoByOrderID(OrderID string) (data entities.SubsTransaction, err error) {
	if err := r.db.Raw(`SELECT * FROM tbl_subscriptions_transaction WHERE order_id=?`, OrderID).Scan(&data).Error; err != nil {
		return entities.SubsTransaction{}, err
	}
	return data, nil
}

func (r *PaymentRepository) GetPeriodName(ID uuid.UUID) (data entities.GetPeriodName, err error) {
	if err := r.db.Raw(`SELECT tmsp.id, tmsp.period_name FROM tbl_master_subs_period tmsp WHERE tmsp.id=?`, ID).Scan(&data).Error; err != nil {
		return entities.GetPeriodName{}, err
	}
	return data, nil
}

func (r *PaymentRepository) WriteUserSubscription(model entities.SubsInfo) (err error) {
	if err := r.db.Create(&model).Error; err != nil {
		return err
	}
	return nil
}

func (r *PaymentRepository) ChangeAccountUser(IDPersonalAccount, IDProAccountUUID uuid.UUID) (err error) {
	var m interface{}

	if err := r.db.Raw(`UPDATE tbl_personal_accounts SET id_master_account_types=?  WHERE id=?`, IDProAccountUUID, IDPersonalAccount).Scan(&m).Error; err != nil {
		return err
	}
	return nil
}