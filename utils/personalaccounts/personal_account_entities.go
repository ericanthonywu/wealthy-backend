package personalaccounts

import "github.com/google/uuid"

type PersonalAccountEntities struct {
	ID                   uuid.UUID `gorm:"column:id"`
	IDMasterAccountTypes uuid.UUID `gorm:"column:id_master_account_types"`
	AccountTypes         string    `gorm:"column:account_type"`
	TotalWallets         int64     `gorm:"column:total_wallet"`
}
