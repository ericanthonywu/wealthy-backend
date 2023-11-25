package tracks

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/tracks/dtos"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/tracks/entities"
	"github.com/semicolon-indonesia/wealthy-backend/utils/errorsinfo"
	"github.com/semicolon-indonesia/wealthy-backend/utils/personalaccounts"
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

	usrEmail := ctx.MustGet("email").(string)
	personalAccount := personalaccounts.Informations(ctx, usrEmail)

	if personalAccount.ID == uuid.Nil {
		httpCode = http.StatusUnauthorized
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", "token contains invalid information")
		return response, httpCode, errInfo
	}

	model.ID = uuid.New()
	model.IDPersonal = personalAccount.ID
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