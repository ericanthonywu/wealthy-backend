package personalaccounts

import (
	"fmt"
	"github.com/wealthy-app/wealthy-backend/constants"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PersonalAccountEntities struct {
	ID                   uuid.UUID `gorm:"column:id"`
	IDMasterAccountTypes uuid.UUID `gorm:"column:id_master_account_types"`
	AccountTypes         string    `gorm:"column:account_type"`
	TotalWallets         int64     `gorm:"column:total_wallet"`
	ReferCode            string    `gorm:"column:refer_code"`
}

func Informations(db *gorm.DB, email string) (data PersonalAccountEntities) {
	db.Raw(`SELECT pa.id,
       pa.id_master_account_types,
       mat.account_type,
       (SELECT COUNT(id) FROM tbl_wallets WHERE id_account = pa.id) as total_wallet,
       pa.refer_code
FROM tbl_personal_accounts pa
         INNER JOIN tbl_master_account_types mat ON mat.id = pa.id_master_account_types
WHERE pa.email = ?`, email).Scan(&data)
	return data
}

func AccountInformation(ctx *gin.Context) (accountType string, accountUUID uuid.UUID) {
	accountType = fmt.Sprintf("%v", ctx.MustGet("accountType"))
	accountUUID = ctx.MustGet("accountID").(uuid.UUID)
	return
}

func PremiumFeature(ctx *gin.Context) bool {
	accountType, _ := AccountInformation(ctx)
	return accountType == constants.AccountBasic
}
