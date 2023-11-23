package subsriptions

import (
	"github.com/semicolon-indonesia/wealthy-backend/infrastructures/subsriptions/entities"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type (
	SubscriptionRepository struct {
		db *gorm.DB
	}

	ISubscriptionRepository interface {
		Plan() (data []entities.SubsPlan)
		FAQ() (data []entities.SubsFAQ)
	}
)

func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) Plan() (data []entities.SubsPlan) {
	if err := r.db.Raw(`SELECT tmat.account_type, tmsp.period_name, tmp.price, tmp.description FROM tbl_master_price tmp
INNER JOIN tbl_master_account_types tmat ON tmat.id = tmp.id_master_account_types
INNER JOIN tbl_master_subs_period tmsp ON tmsp.id = tmp.id_master_subs_period
WHERE tmp.active=true`).Scan(&data).Error; err != nil {
		logrus.Error(err.Error())
		return []entities.SubsPlan{}
	}
	return data
}

func (r *SubscriptionRepository) FAQ() (data []entities.SubsFAQ) {
	if err := r.db.Raw(`SELECT tf.id as id, tf.questions, tf.answer FROM tbl_faq tf WHERE tf.active=true`).Scan(&data).Error; err != nil {
		logrus.Error(err.Error())
		return []entities.SubsFAQ{}
	}

	return data
}