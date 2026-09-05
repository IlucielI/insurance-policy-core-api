package http

import (
	"log"
	"net/http"

	"github.com/IlucielI/insurance-policy-core-api/internal/infrastructure/websocket"
	"github.com/IlucielI/insurance-policy-core-api/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	ws "github.com/gorilla/websocket"
)

type NotificationHandler struct {
	usecase *usecase.NotificationUsecase
	hub     *websocket.Hub
}

var upgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins (configure domain whitelist in production)
		return true
	},
}

func NewNotificationHandler(usecase *usecase.NotificationUsecase, hub *websocket.Hub) *NotificationHandler {
	return &NotificationHandler{
		usecase: usecase,
		hub:     hub,
	}
}

// WebSocketHandler handles WebSocket connections
func (h *NotificationHandler) WebSocketHandler(c *fiber.Ctx) error {
	userID := c.Query("user_id")
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user_id required",
		})
	}

	// Use Fiber adaptor to convert fasthttp -> net/http for gorilla/websocket
	return adaptor.HTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("WebSocket upgrade failed: %v", err)
			return
		}

		client := websocket.NewClient(h.hub, conn, userID)
		h.hub.Register <- client

		go client.WritePump()
		go client.ReadPump()
	})(c)
}

// GetNotifications returns user notifications
func (h *NotificationHandler) GetNotifications(c *fiber.Ctx) error {
	userID := c.Query("user_id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user_id required",
		})
	}

	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)

	notifications, err := h.usecase.GetByUserID(userID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"notifications": notifications,
	})
}

// MarkAsRead marks a notification as read
func (h *NotificationHandler) MarkAsRead(c *fiber.Ctx) error {
	notificationID := c.Params("id")
	if notificationID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "notification_id required",
		})
	}

	err := h.usecase.MarkAsRead(notificationID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Notification marked as read",
	})
}

// MarkAllAsRead marks all user notifications as read
func (h *NotificationHandler) MarkAllAsRead(c *fiber.Ctx) error {
	userID := c.Query("user_id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user_id required",
		})
	}

	err := h.usecase.MarkAllAsRead(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "All notifications marked as read",
	})
}

// GetUnreadCount returns unread notification count
func (h *NotificationHandler) GetUnreadCount(c *fiber.Ctx) error {
	userID := c.Query("user_id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user_id required",
		})
	}

	count, err := h.usecase.GetUnreadCount(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"unread_count": count,
	})
}

// DeleteNotification deletes a notification
func (h *NotificationHandler) DeleteNotification(c *fiber.Ctx) error {
	notificationID := c.Params("id")
	if notificationID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "notification_id required",
		})
	}

	err := h.usecase.Delete(notificationID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Notification deleted",
	})
}