package utilities

import "github.com/google/uuid"

type (
	NotificationEntities struct {
		ID                      uuid.UUID `gorm:"column:id"`
		NotificationTitle       string    `gorm:"column:notification_title"`
		NotificationDescription string    `gorm:"column:notification_description"`
		IDPersonalAccounts      uuid.UUID `gorm:"column:id_personal_accounts"`
		IsRead                  bool      `gorm:"column:is_read"`
		IDGroupSharing          string    `gorm:"column:id_group_sharing"`
	}
)

func (NotificationEntities) TableName() string {
	return "tbl_notifications"
}