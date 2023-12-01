package notifications

import (
	"github.com/gin-gonic/gin"
	"github.com/semicolon-indonesia/wealthy-backend/utils/response"
)

type (
	NotificationController struct {
		useCase INotificationUseCase
	}

	INotificationController interface {
		GetNotification(ctx *gin.Context)
	}
)

func NewNotificationController(useCase INotificationUseCase) *NotificationController {
	return &NotificationController{useCase: useCase}
}

func (c *NotificationController) GetNotification(ctx *gin.Context) {
	data, httpCode, errInfo := c.useCase.GetNotification(ctx)
	response.SendBack(ctx, data, errInfo, httpCode)
	return
}