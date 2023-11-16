package referrals

import "gorm.io/gorm"

type (
	ReferralRepository struct {
		db *gorm.DB
	}

	IReferralRepository interface {
	}
)

func NewReferralRepository(db *gorm.DB) *ReferralRepository {
	return &ReferralRepository{db: db}
}