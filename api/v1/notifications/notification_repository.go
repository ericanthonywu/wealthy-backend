package notifications

import "gorm.io/gorm"

type (
	NotificationRepository struct {
		db *gorm.DB
	}

	INotificationRepository interface {
	}
)

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}