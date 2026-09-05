package http

import (
	"strconv"

	"github.com/IlucielI/insurance-policy-core-api/internal/middleware"
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

// Register godoc
// @Summary Register user baru
// @Description Mendaftarkan akun user baru dengan email, password, nama lengkap, dan nomor telepon
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Request body untuk registrasi"
// @Success 201 {object} map[string]interface{} "User berhasil didaftarkan"
// @Failure 400 {object} map[string]interface{} "Invalid request body atau email sudah terdaftar"
// @Router /auth/register [post]
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

// Login godoc
// @Summary Login user
// @Description Autentikasi user dengan email dan password, mengembalikan JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Kredensial login"
// @Success 200 {object} map[string]interface{} "Login berhasil dengan token dan informasi user"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 401 {object} map[string]interface{} "Email atau password salah"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	user, roles, err := h.authUsecase.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Generate JWT token with role claims
	token, err := middleware.GenerateToken(user.ID, user.Email, roles)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	// Set token as cookie
	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    token,
		HTTPOnly: true,
		Secure:   false, // Set true in production with HTTPS
		SameSite: "Lax",
	})

	// Also keep legacy cookie for backward compatibility
	c.Cookie(&fiber.Cookie{
		Name:  "user_id",
		Value: user.ID,
	})

	return c.JSON(fiber.Map{
		"message": "Login successful",
		"user":    user,
		"token":   token,
		"roles":   roles,
	})
}

// Logout godoc
// @Summary Logout user
// @Description Menghapus cookie autentikasi dan melakukan logout
// @Tags Auth
// @Produce json
// @Success 200 {object} map[string]interface{} "Logout berhasil"
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	c.ClearCookie("user_id")
	return c.JSON(fiber.Map{
		"message": "Logout successful",
	})
}

// Me godoc
// @Summary Get current user
// @Description Mendapatkan informasi user yang sedang login berdasarkan JWT token
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Security CookieAuth
// @Success 200 {object} map[string]interface{} "Informasi user"
// @Failure 401 {object} map[string]interface{} "Not authenticated"
// @Router /auth/me [get]
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	// Try to get from JWT claims first
	claims, err := middleware.GetUserFromContext(c)
	if err == nil && claims != nil {
		user, err := h.authUsecase.GetUserByID(c.Context(), claims.UserID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(fiber.Map{
			"user":  user,
			"roles": claims.Roles,
		})
	}

	// Fallback to cookie-based auth (backward compatibility)
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

	return c.JSON(fiber.Map{
		"user":  user,
		"roles": []string{user.Role},
	})
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

// ForgotPassword godoc
// @Summary Request password reset
// @Description Mengirim email reset password ke alamat email terdaftar
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body ForgotPasswordRequest true "Email untuk reset password"
// @Success 200 {object} map[string]interface{} "Email reset password telah dikirim"
// @Failure 400 {object} map[string]interface{} "Email is required"
// @Router /auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *fiber.Ctx) error {
	var req ForgotPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email is required",
		})
	}

	if err := h.authUsecase.RequestPasswordReset(c.Context(), req.Email); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "If the email exists, a password reset link has been sent",
	})
}

// ResetPassword godoc
// @Summary Reset password
// @Description Mengatur password baru menggunakan token reset yang diterima via email
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body ResetPasswordRequest true "Token dan password baru"
// @Success 200 {object} map[string]interface{} "Password berhasil direset"
// @Failure 400 {object} map[string]interface{} "Token invalid atau password terlalu pendek"
// @Router /auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	var req ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Token == "" || req.NewPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Token and new password are required",
		})
	}

	if len(req.NewPassword) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Password must be at least 6 characters",
		})
	}

	if err := h.authUsecase.ResetPassword(c.Context(), req.Token, req.NewPassword); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Password reset successfully",
	})
}

// Admin handlers

// GET /admin/users - List all users (super_admin)
func (h *AuthHandler) ListUsers(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	users, total, err := h.authUsecase.ListUsers(c.Context(), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data":   users,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// POST /admin/users/:id/roles - Assign role (super_admin)
type AssignRoleRequest struct {
	RoleID string `json:"role_id"`
}

func (h *AuthHandler) AssignRole(c *fiber.Ctx) error {
	userID := c.Params("id")
	var req AssignRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Get assignedBy from JWT claims
	claims, err := middleware.GetUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "not authenticated",
		})
	}

	if err := h.authUsecase.AssignRole(c.Context(), userID, req.RoleID, claims.UserID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "role assigned",
	})
}

// DELETE /admin/users/:id/roles - Remove role (super_admin)
type RemoveRoleRequest struct {
	RoleID string `json:"role_id"`
}

func (h *AuthHandler) RemoveRole(c *fiber.Ctx) error {
	userID := c.Params("id")
	var req RemoveRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if err := h.authUsecase.RemoveRole(c.Context(), userID, req.RoleID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "role removed",
	})
}

// GET /admin/roles - List all roles (super_admin)
func (h *AuthHandler) GetRoles(c *fiber.Ctx) error {
	roles, err := h.authUsecase.GetRoles(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": roles,
	})
}
