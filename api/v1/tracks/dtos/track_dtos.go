package dtos

type (
	ScreenTimeRequest struct {
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
	}
)