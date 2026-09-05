package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestExtractAuditContext(t *testing.T) {
	t.Run("with user_id in locals", func(t *testing.T) {
		app := fiber.New()
		app.Get("/test", func(c *fiber.Ctx) error {
			c.Locals("user_id", "user-123")
			ctx := ExtractAuditContext(c)
			assert.NotNil(t, ctx)
			assert.NotNil(t, ctx.UserID)
			assert.Equal(t, "user-123", *ctx.UserID)
			assert.NotEmpty(t, ctx.IPAddress)
			return c.SendString("ok")
		})

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("User-Agent", "TestAgent/1.0")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("without user_id", func(t *testing.T) {
		app := fiber.New()
		app.Get("/test", func(c *fiber.Ctx) error {
			ctx := ExtractAuditContext(c)
			assert.NotNil(t, ctx)
			assert.Nil(t, ctx.UserID)
			assert.NotEmpty(t, ctx.IPAddress)
			return c.SendString("ok")
		})

		req := httptest.NewRequest("GET", "/test", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("captures user agent", func(t *testing.T) {
		app := fiber.New()
		app.Get("/test", func(c *fiber.Ctx) error {
			ctx := ExtractAuditContext(c)
			assert.Equal(t, "CustomAgent/2.0", ctx.UserAgent)
			return c.SendString("ok")
		})

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("User-Agent", "CustomAgent/2.0")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})
}
