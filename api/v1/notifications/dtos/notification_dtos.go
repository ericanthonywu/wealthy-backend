package dtos

import (
	"github.com/google/uuid"
)

type (
	Notification struct {
		ID                      uuid.UUID `json:"id"`
		NotificationTitle       string    `json:"notification_title"`
		NotificationDescription string    `json:"notification_description"`
		IDPersonalAccounts      uuid.UUID `json:"id_personal_accounts"`
		IsRead                  bool      `json:"is_read"`
		IDGroupSender           uuid.UUID `json:"id_group_sender"`
		IDGroupRecipient        uuid.UUID `json:"id_group_recipient"`
		CreatedAt               string    `json:"created_at"`
	}
)