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
		IDGroupSender           uuid.UUID `gorm:"column:id_group_sender"`
		IDGroupRecipient        uuid.UUID `gorm:"column:id_group_recipient"`
		ImagePath               string    `gorm:"column:image_path"`
		Type                    string    `gorm:"column:type"`
		CreatedAt               time.Time `gorm:"column:created_at"`
	}
)