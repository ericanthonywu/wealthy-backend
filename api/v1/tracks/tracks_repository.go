package tracks

import "gorm.io/gorm"

type (
	TrackRepository struct {
		db *gorm.DB
	}

	ITrackRepository interface {
		ScreenTime()
	}
)

func NewTrackRepository(db *gorm.DB) *TrackRepository {
	return &TrackRepository{db: db}
}

func (r *TrackRepository) ScreenTime() {

}