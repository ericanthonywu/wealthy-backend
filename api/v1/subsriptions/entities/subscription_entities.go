package entities

import "github.com/google/uuid"

type (
	SubsPlan struct {
		AccountType string  `gorm:"column:account_type" json:"account_type"`
		PeriodName  string  `gorm:"column:period_name" json:"period_name"`
		Price       float64 `gorm:"column:price" json:"price"`
		Description string  `gorm:"column:description" json:"description"`
	}

	SubsFAQ struct {
		ID        uuid.UUID `gorm:"column:id" json:"id"`
		Questions string    `gorm:"column:questions" json:"questions"`
		Answer    string    `gorm:"column:answer" json:"answer"`
	}
)