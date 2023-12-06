package notifications

import (
	"github.com/google/uuid"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/notifications/entities"
	"gorm.io/gorm"
)

type (
	NotificationRepository struct {
		db *gorm.DB
	}

	INotificationRepository interface {
		GetNotification(personalAccount uuid.UUID) (data []entities.NotificationEntities, err error)
	}
)

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) GetNotification(personalAccount uuid.UUID) (data []entities.NotificationEntities, err error) {
	if err := r.db.Raw(`SELECT  tn.id, tn.notification_title, tn.notification_description, tn.id_personal_accounts, tn.is_read, tn.id_group_sender, tn.id_group_recipient
FROM tbl_notifications tn WHERE tn.id_personal_accounts=? AND tn.is_read=false ORDER BY created_at DESC`, personalAccount).Scan(&data).Error; err != nil {
		return []entities.NotificationEntities{}, err
	}
	return data, nil
}