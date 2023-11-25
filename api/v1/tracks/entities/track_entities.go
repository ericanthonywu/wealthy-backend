package entities

import "github.com/google/uuid"

type (
	ScreenTime struct {
		ID         uuid.UUID `gorm:"column:id"`
		IDPersonal uuid.UUID `gorm:"column:id_personal_accounts"`
		StartDate  string    `gorm:"column:start_date"`
		EndDate    string    `gorm:"column:end_date"`
	}
)

func (ScreenTime) TableName() string {
	return "tbl_access_time"
}