package http

import (
	"github.com/IlucielI/insurance-policy-core-api/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authUsecase *usecase.AuthUsecase
}

func NewAuthHandler(authUsecase *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	user, err := h.authUsecase.Register(c.Context(), req.Email, req.Password, req.FullName, req.Phone)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"user":    user,
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	user, err := h.authUsecase.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Set session (simplified - should use proper session store)
	// TODO: Implement Redis session store
	c.Cookie(&fiber.Cookie{
		Name:  "user_id",
		Value: user.ID,
	})

	return c.JSON(fiber.Map{
		"message": "Login successful",
		"user":    user,
	})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	c.ClearCookie("user_id")
	return c.JSON(fiber.Map{
		"message": "Logout successful",
	})
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	userID := c.Cookies("user_id")
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Not authenticated",
		})
	}

	user, err := h.authUsecase.GetUserByID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(user)
}
