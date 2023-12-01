package utilities

import "github.com/google/uuid"

type (
	NotificationEntities struct {
		ID                      uuid.UUID `gorm:"column:id"`
		NotificationTitle       string    `gorm:"column:notification_title"`
		NotificationDescription string    `gorm:"column:notification_description"`
		IDPersonalAccounts      uuid.UUID `gorm:"column:id_personal_accounts"`
		IsRead                  bool      `gorm:"column:is_read"`
		IDGroupSender           string    `gorm:"column:id_group_sender"`
		IDGroupReceipt          string    `gorm:"column:id_group_recipient"`
	}
)

func (NotificationEntities) TableName() string {
	return "tbl_notifications"
}