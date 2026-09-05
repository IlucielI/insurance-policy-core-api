package usecase

import (
	"github.com/IlucielI/insurance-policy-core-api/internal/infrastructure/websocket"
	"github.com/IlucielI/insurance-policy-core-api/internal/model"
	"github.com/IlucielI/insurance-policy-core-api/internal/repository"
)

type NotificationUsecase struct {
	repo *repository.NotificationRepository
	hub  *websocket.Hub
}

func NewNotificationUsecase(repo *repository.NotificationRepository, hub *websocket.Hub) *NotificationUsecase {
	return &NotificationUsecase{
		repo: repo,
		hub:  hub,
	}
}

func (u *NotificationUsecase) Create(req *model.NotificationCreateRequest) (*model.Notification, error) {
	notification, err := u.repo.Create(req)
	if err != nil {
		return nil, err
	}

	// Send real-time notification via WebSocket
	u.hub.SendNotification(req.UserID, notification)

	return notification, nil
}

func (u *NotificationUsecase) GetByUserID(userID string, limit, offset int) ([]*model.Notification, error) {
	if limit <= 0 {
		limit = 20
	}
	return u.repo.GetByUserID(userID, limit, offset)
}

func (u *NotificationUsecase) MarkAsRead(notificationID string) error {
	return u.repo.MarkAsRead(notificationID)
}

func (u *NotificationUsecase) MarkAllAsRead(userID string) error {
	return u.repo.MarkAllAsRead(userID)
}

func (u *NotificationUsecase) GetUnreadCount(userID string) (int, error) {
	return u.repo.GetUnreadCount(userID)
}

func (u *NotificationUsecase) Delete(notificationID string) error {
	return u.repo.Delete(notificationID)
}
