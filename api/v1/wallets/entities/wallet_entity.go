package entities

import "github.com/google/uuid"

type (
	WalletPersonalInformationEntity struct {
		ID                   uuid.UUID `gorm:"column:id"`
		IDMasterAccountTypes uuid.UUID `gorm:"column:id_master_account_types"`
		AccountTypes         string    `gorm:"column:account_types"`
	}
)
