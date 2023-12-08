package investments

import "gorm.io/gorm"

type (
	InvestmentRepository struct {
		db *gorm.DB
	}

	IInvestmentRepository interface {
	}
)

func NewInvestmentRepository(db *gorm.DB) *InvestmentRepository {
	return &InvestmentRepository{db: db}
}