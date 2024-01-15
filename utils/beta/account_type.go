package beta

import (
	"gorm.io/gorm"
)

type ExpiredData struct {
	Expired string `gorm:"column:expired"`
}

func ExpiredPromotion(db *gorm.DB) (data ExpiredData) {
	db.Raw("SELECT expired FROM tbl_beta_promotion LIMIT 1").Scan(&data)
	return data
}
