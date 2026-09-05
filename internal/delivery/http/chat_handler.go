package http

import (
	"bufio"
	"fmt"

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
	Stream    bool   `json:"stream"`
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

	// Handle streaming response
	if req.Stream {
		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")
		c.Set("Transfer-Encoding", "chunked")

		c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
			err := h.chatUsecase.SendMessageStream(c.Context(), req.SessionID, req.Message, userID, w)
			if err != nil {
				fmt.Fprintf(w, "data: {\"error\": \"%s\"}\n\n", err.Error())
			}
			w.Flush()
		})

		return nil
	}

	// Non-streaming response
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
