package images

import "gorm.io/gorm"

type (
	ShowImageRepository struct {
		db *gorm.DB
	}

	IShowImageRepository interface{}
)

func NewShowImageRepository(db *gorm.DB) *ShowImageRepository {
	return &ShowImageRepository{db: db}
}