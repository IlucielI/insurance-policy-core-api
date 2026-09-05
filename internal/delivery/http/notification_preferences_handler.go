package http

import (
	"log"

	"github.com/IlucielI/insurance-policy-core-api/internal/domain"
	"github.com/IlucielI/insurance-policy-core-api/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type NotificationPreferencesHandler struct {
	usecase *usecase.NotificationPreferencesUsecase
}

func NewNotificationPreferencesHandler(usecase *usecase.NotificationPreferencesUsecase) *NotificationPreferencesHandler {
	return &NotificationPreferencesHandler{usecase: usecase}
}

type UpdatePreferencesRequest struct {
	PromotionalEmails       bool `json:"promotional_emails"`
	PolicyUpdateEmails      bool `json:"policy_update_emails"`
	ClaimNotificationEmails bool `json:"claim_notification_emails"`
	NewsletterEmails        bool `json:"newsletter_emails"`
}

// GetPreferences godoc
// @Summary Get notification preferences
// @Tags notification-preferences
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/notification-preferences [get]
func (h *NotificationPreferencesHandler) GetPreferences(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	prefs, err := h.usecase.GetPreferences(userID)
	if err != nil {
		log.Printf("Failed to get preferences: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get notification preferences",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    prefs,
	})
}

// UpdatePreferences godoc
// @Summary Update notification preferences
// @Tags notification-preferences
// @Accept json
// @Produce json
// @Param request body UpdatePreferencesRequest true "Preferences"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/notification-preferences [put]
func (h *NotificationPreferencesHandler) UpdatePreferences(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var req UpdatePreferencesRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	prefs := &domain.NotificationPreferences{
		UserID:                  userID,
		PromotionalEmails:       req.PromotionalEmails,
		PolicyUpdateEmails:      req.PolicyUpdateEmails,
		ClaimNotificationEmails: req.ClaimNotificationEmails,
		NewsletterEmails:        req.NewsletterEmails,
	}

	if err := h.usecase.UpdatePreferences(userID, prefs); err != nil {
		log.Printf("Failed to update preferences: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update notification preferences",
		})
	}

	// Get updated preferences to return
	updatedPrefs, _ := h.usecase.GetPreferences(userID)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Preferensi notifikasi berhasil diperbarui",
		"data":    updatedPrefs,
	})
}
