package http

import (
	"github.com/IlucielI/insurance-policy-core-api/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type ChatHandler struct {
	chatUsecase *usecase.ChatUsecase
}

func NewChatHandler(chatUsecase *usecase.ChatUsecase) *ChatHandler {
	return &ChatHandler{chatUsecase: chatUsecase}
}

type ChatRequest struct {
	SessionID string `json:"session_id"`
	Message   string `json:"message"`
}

func (h *ChatHandler) SendMessage(c *fiber.Ctx) error {
	var req ChatRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.SessionID == "" || req.Message == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "session_id and message are required",
		})
	}

	// Get user ID from session (if logged in)
	var userID *string
	userIDCookie := c.Cookies("user_id")
	if userIDCookie != "" {
		userID = &userIDCookie
	}

	// Send message and get response
	reply, sources, err := h.chatUsecase.SendMessage(c.Context(), req.SessionID, req.Message, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"reply":   reply,
		"sources": sources,
	})
}

func (h *ChatHandler) GetHistory(c *fiber.Ctx) error {
	sessionID := c.Query("session_id")
	if sessionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "session_id is required",
		})
	}

	messages, err := h.chatUsecase.GetChatHistory(c.Context(), sessionID, 50)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"messages": messages,
	})
}
