package beta

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ExpiredData struct {
	Expired string `gorm:"column:expired"`
}

func ExpiredPromotion(ctx *gin.Context) (data ExpiredData) {
	db := ctx.MustGet("db").(*gorm.DB)
	db.Raw("SELECT expired FROM tbl_beta_promotion LIMIT 1").Scan(&data)
	return data
}