package subsriptions

import (
	"github.com/sirupsen/logrus"
	"github.com/wealthy-app/wealthy-backend/api/v1/subsriptions/entities"
	"gorm.io/gorm"
)

type (
	SubscriptionRepository struct {
		db *gorm.DB
	}

	ISubscriptionRepository interface {
		FAQ() (data []entities.SubsFAQ)
	}
)

func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) FAQ() (data []entities.SubsFAQ) {
	if err := r.db.Raw(`SELECT tf.id as id, tf.questions, tf.answer FROM tbl_faq tf WHERE tf.active=true`).Scan(&data).Error; err != nil {
		logrus.Error(err.Error())
		return []entities.SubsFAQ{}
	}

	return data
}