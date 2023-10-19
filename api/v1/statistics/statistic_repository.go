package statistics

import "gorm.io/gorm"

type (
	StatisticRepository struct {
		db *gorm.DB
	}

	IStatisticRepository interface {
	}
)

func NewStatisticRepository(db *gorm.DB) *StatisticRepository {
	return &StatisticRepository{db: db}
}
