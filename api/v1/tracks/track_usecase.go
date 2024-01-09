package tracks

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wealthy-app/wealthy-backend/api/v1/tracks/dtos"
	"github.com/wealthy-app/wealthy-backend/api/v1/tracks/entities"
	"github.com/wealthy-app/wealthy-backend/utils/errorsinfo"
	"net/http"
)

type (
	TrackUseCase struct {
		repo ITrackRepository
	}

	ITrackUseCase interface {
		ScreenTime(ctx *gin.Context, dtoRequest *dtos.ScreenTimeRequest) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
	}
)

func NewTrackUseCase(repo ITrackRepository) *TrackUseCase {
	return &TrackUseCase{repo: repo}
}

func (s *TrackUseCase) ScreenTime(ctx *gin.Context, dtoRequest *dtos.ScreenTimeRequest) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var model entities.ScreenTime

	accountUUID := ctx.MustGet("accountID").(uuid.UUID)

	model.ID = uuid.New()
	model.IDPersonal = accountUUID
	model.StartDate = dtoRequest.StartTime
	model.EndDate = dtoRequest.EndTime

	err := s.repo.ScreenTime(&model)
	if err != nil {
		result := struct{}{}
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		return result, http.StatusInternalServerError, errInfo
	}

	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	result := struct {
		Message string `json:"message"`
	}{
		Message: "save successfully",
	}

	return result, http.StatusOK, errInfo
}