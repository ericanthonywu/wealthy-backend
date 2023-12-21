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
	if err := r.db.Raw(`SELECT tn.id,
       tn.notification_title,
       tn.notification_description,
       tn.id_personal_accounts,
       tn.is_read,
       tn.id_group_sharing,
       tpa.image_path,
       tmat.account_type as type,
       tn.created_at
FROM tbl_notifications tn
         INNER JOIN tbl_personal_accounts tpa ON tpa.id = tn.id_personal_accounts
         INNER JOIN tbl_master_account_types tmat ON tmat.id = tpa.id_master_account_types
WHERE tn.id_personal_accounts = ?
  AND tn.is_read = false
ORDER BY created_at DESC`, personalAccount).Scan(&data).Error; err != nil {
		return []entities.NotificationEntities{}, err
	}
	return data, nil
}