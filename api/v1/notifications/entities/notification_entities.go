package entities

import (
	"github.com/google/uuid"
	"time"
)

type (
	NotificationEntities struct {
		ID                      uuid.UUID `gorm:"column:id" json:"id"`
		NotificationTitle       string    `gorm:"column:notification_title" json:"notification_title"`
		NotificationDescription string    `gorm:"column:notification_description" json:"notification_description"`
		IDPersonalAccounts      uuid.UUID `gorm:"column:id_personal_accounts" json:"id_personal_accounts"`
		IsRead                  bool      `gorm:"column:is_read" json:"is_read"`
		IDGroupSender           uuid.UUID `gorm:"column:id_group_sender" json:"id_group_sender"`
		IDGroupRecipient        uuid.UUID `gorm:"column:id_group_recipient" json:"id_group_recipient"`
		CreatedAt               time.Time `gorm:"column:created_at" json:"created_at"`
	}
)
