package notifications

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/semicolon-indonesia/wealthy-backend/api/v1/notifications/dtos"
	"github.com/semicolon-indonesia/wealthy-backend/utils/datecustoms"
	"github.com/semicolon-indonesia/wealthy-backend/utils/errorsinfo"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

type (
	NotificationUseCase struct {
		repo INotificationRepository
	}

	INotificationUseCase interface {
		GetNotification(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors)
	}
)

func NewNotificationUseCase(repo INotificationRepository) *NotificationUseCase {
	return &NotificationUseCase{repo: repo}
}

func (s *NotificationUseCase) GetNotification(ctx *gin.Context) (response interface{}, httpCode int, errInfo []errorsinfo.Errors) {
	var dtoResponse []dtos.Notification

	accountUUID := ctx.MustGet("accountID").(uuid.UUID)

	// get notification info by personal ID
	dataNotification, err := s.repo.GetNotification(accountUUID)
	if err != nil {
		logrus.Error(err.Error())
		errInfo = errorsinfo.ErrorWrapper(errInfo, "", err.Error())
		resp := struct {
			Message string `json:"message,omitempty"`
		}{}
		return resp, http.StatusInternalServerError, errInfo
	}

	// if not error message
	if len(errInfo) == 0 {
		errInfo = []errorsinfo.Errors{}
	}

	// if data not found
	if len(dataNotification) == 0 {
		resp := struct {
			Message string `json:"message,omitempty"`
		}{
			Message: "no new notification update",
		}
		return resp, http.StatusOK, errInfo
	}

	for _, v := range dataNotification {
		dtoResponse = append(dtoResponse, dtos.Notification{
			ID:                      v.ID,
			NotificationTitle:       v.NotificationTitle,
			NotificationDescription: v.NotificationDescription,
			IDPersonalAccounts:      v.IDPersonalAccounts,
			IsRead:                  v.IsRead,
			IDGroupSharing:          v.IDGroupSharing,
			AccountDetail: dtos.NotificationAccountDetail{
				AccountImage: os.Getenv("APP_HOST") + "/v1/" + v.ImagePath,
				AccountType:  v.Type,
			},
			CreatedAt: datecustoms.TimeRFC3339ToString(v.CreatedAt),
		})
	}

	return dtoResponse, http.StatusOK, errInfo
}