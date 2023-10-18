package personalaccounts

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Informations(ctx *gin.Context, email string) (data PersonalAccountEntities) {
	db := ctx.MustGet("db").(*gorm.DB)
	db.Raw(" SELECT pa.id,pa.id_master_account_types, mat.account_type, (SELECT COUNT(id) FROM tbl_wallets WHERE id_account = pa.id) as total_wallet "+
		"FROM tbl_personal_accounts pa "+
		"INNER JOIN tbl_master_account_types mat ON mat.id = pa.id_master_account_types "+
		"WHERE pa.email = ?", email).Scan(&data)
	return data
}
