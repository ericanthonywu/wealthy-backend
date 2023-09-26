package masters

import "gorm.io/gorm"

type (
	MasterRepository struct {
		db *gorm.DB
	}

	IMasterRepository interface {
	}
)

func NewMasterRepository(db *gorm.DB) *MasterRepository {
	return &MasterRepository{db: db}
}
