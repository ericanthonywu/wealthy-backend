package tracks

import (
	"github.com/sirupsen/logrus"
	"github.com/wealthy-app/wealthy-backend/api/v1/tracks/entities"
	"gorm.io/gorm"
)

type (
	TrackRepository struct {
		db *gorm.DB
	}

	ITrackRepository interface {
		ScreenTime(model *entities.ScreenTime) (err error)
	}
)

func NewTrackRepository(db *gorm.DB) *TrackRepository {
	return &TrackRepository{db: db}
}

func (r *TrackRepository) ScreenTime(model *entities.ScreenTime) (err error) {
	if err := r.db.Create(&model).Error; err != nil {
		logrus.Error(err.Error())
		return err
	}
	return nil
}