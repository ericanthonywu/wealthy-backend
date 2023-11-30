package notifications

type (
	NotificationUseCase struct {
		repo INotificationRepository
	}

	INotificationUseCase interface {
	}
)

func NewNotificationUseCase(repo INotificationRepository) *NotificationUseCase {
	return &NotificationUseCase{repo: repo}
}