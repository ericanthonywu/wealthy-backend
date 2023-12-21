package dtos

import (
	"github.com/google/uuid"
)

type (
	Notification struct {
		ID                      uuid.UUID                 `json:"id"`
		NotificationTitle       string                    `json:"notification_title"`
		NotificationDescription string                    `json:"notification_description"`
		IDPersonalAccounts      uuid.UUID                 `json:"id_personal_accounts"`
		IsRead                  bool                      `json:"is_read"`
		IDGroupSharing          uuid.UUID                 `json:"id_group_sharing"`
		AccountDetail           NotificationAccountDetail `json:"account_detail"`
		CreatedAt               string                    `json:"created_at"`
	}

	NotificationAccountDetail struct {
		AccountImage string `json:"account_image"`
		AccountType  string `json:"account_type"`
	}
)