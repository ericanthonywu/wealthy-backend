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
		MidtransWebhook(orderID string) (err error)
	}
)

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) GetPrice(subs uuid.UUID) (data entities.DataPriceInfo) {
	if err := r.db.Raw(`SELECT tmp.price, tmat.account_type, tmp.description, tmsp.period_name FROM tbl_master_price tmp
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

func (r *PaymentRepository) MidtransWebhook(orderID string) (err error) {
	var model interface{}

	if err := r.db.Raw(`UPDATE tbl_subscriptions_transaction SET status=? WHERE order_id=?`, 1, orderID).Scan(&model).Error; err != nil {
		return err
	}

	return nil
}