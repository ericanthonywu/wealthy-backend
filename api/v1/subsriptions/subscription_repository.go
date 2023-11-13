package subsriptions

import (
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/subsriptions/entities"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type (
	SubscriptionRepository struct {
		db *gorm.DB
	}

	ISubscriptionRepository interface {
		Plan() (data []entities.SubsPlan)
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