package entities

import (
	"github.com/google/uuid"
	"time"
)

type (
	NotificationEntities struct {
		ID                      uuid.UUID `gorm:"column:id"`
		NotificationTitle       string    `gorm:"column:notification_title"`
		NotificationDescription string    `gorm:"column:notification_description"`
		IDPersonalAccounts      uuid.UUID `gorm:"column:id_personal_accounts"`
		IsRead                  bool      `gorm:"column:is_read"`
		IDGroupSharing          uuid.UUID `gorm:"column:id_group_sharing"`
		ImagePath               string    `gorm:"column:image_path"`
		Type                    string    `gorm:"column:type"`
		CreatedAt               time.Time `gorm:"column:created_at"`
	}
)