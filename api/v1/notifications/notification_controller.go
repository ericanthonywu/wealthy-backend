package notifications

type (
	NotificationController struct {
		useCase INotificationUseCase
	}

	INotificationController interface {
	}
)

func NewNotificationController(useCase INotificationUseCase) *NotificationController {
	return &NotificationController{useCase: useCase}
}