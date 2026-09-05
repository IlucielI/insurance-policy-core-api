package middleware

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		token, err := GenerateToken("user-123", "test@example.com", []string{"customer"})
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		// Verify token can be parsed
		parsedToken, err := jwt.ParseWithClaims(token, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		assert.NoError(t, err)
		assert.True(t, parsedToken.Valid)

		claims, ok := parsedToken.Claims.(*UserClaims)
		assert.True(t, ok)
		assert.Equal(t, "user-123", claims.UserID)
		assert.Equal(t, "test@example.com", claims.Email)
		assert.Equal(t, []string{"customer"}, claims.Roles)
	})
}

func TestAuthRequired(t *testing.T) {
	t.Run("valid token in header", func(t *testing.T) {
		app := fiber.New()
		app.Use(AuthRequired())
		app.Get("/test", func(c *fiber.Ctx) error {
			return c.SendString("ok")
		})

		token, _ := GenerateToken("user-123", "test@example.com", []string{"customer"})
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("missing token", func(t *testing.T) {
		app := fiber.New()
		app.Use(AuthRequired())
		app.Get("/test", func(c *fiber.Ctx) error {
			return c.SendString("ok")
		})

		req := httptest.NewRequest("GET", "/test", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 401, resp.StatusCode)
	})

	t.Run("invalid token", func(t *testing.T) {
		app := fiber.New()
		app.Use(AuthRequired())
		app.Get("/test", func(c *fiber.Ctx) error {
			return c.SendString("ok")
		})

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer invalid_token")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 401, resp.StatusCode)
	})

	t.Run("expired token", func(t *testing.T) {
		app := fiber.New()
		app.Use(AuthRequired())
		app.Get("/test", func(c *fiber.Ctx) error {
			return c.SendString("ok")
		})

		// Create expired token
		claims := UserClaims{
			UserID: "user-123",
			Email:  "test@example.com",
			Roles:  []string{"customer"},
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, _ := token.SignedString(jwtSecret)

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+tokenString)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 401, resp.StatusCode)
	})
}

func TestRequireRole(t *testing.T) {
	t.Run("has required role", func(t *testing.T) {
		app := fiber.New()
		app.Use(AuthRequired())
		app.Use(RequireRole("admin", "super_admin"))
		app.Get("/test", func(c *fiber.Ctx) error {
			return c.SendString("ok")
		})

		token, _ := GenerateToken("user-123", "admin@example.com", []string{"admin"})
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("missing required role", func(t *testing.T) {
		app := fiber.New()
		app.Use(AuthRequired())
		app.Use(RequireRole("admin", "super_admin"))
		app.Get("/test", func(c *fiber.Ctx) error {
			return c.SendString("ok")
		})

		token, _ := GenerateToken("user-123", "customer@example.com", []string{"customer"})
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 403, resp.StatusCode)
	})

	t.Run("not authenticated", func(t *testing.T) {
		app := fiber.New()
		// Skip AuthRequired middleware
		app.Use(RequireRole("admin"))
		app.Get("/test", func(c *fiber.Ctx) error {
			return c.SendString("ok")
		})

		req := httptest.NewRequest("GET", "/test", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 401, resp.StatusCode)
	})
}

func TestRequirePermission(t *testing.T) {
	t.Run("super_admin has all permissions", func(t *testing.T) {
		app := fiber.New()
		app.Use(AuthRequired())
		app.Use(RequirePermission("delete_users"))
		app.Get("/test", func(c *fiber.Ctx) error {
			return c.SendString("ok")
		})

		token, _ := GenerateToken("user-123", "superadmin@example.com", []string{"super_admin"})
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("not authenticated", func(t *testing.T) {
		app := fiber.New()
		app.Use(RequirePermission("test_permission"))
		app.Get("/test", func(c *fiber.Ctx) error {
			return c.SendString("ok")
		})

		req := httptest.NewRequest("GET", "/test", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 401, resp.StatusCode)
	})
}

func TestGetUserFromContext(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		app := fiber.New()
		app.Use(AuthRequired())
		app.Get("/test", func(c *fiber.Ctx) error {
			user, err := GetUserFromContext(c)
			assert.NoError(t, err)
			assert.NotNil(t, user)
			assert.Equal(t, "user-123", user.UserID)
			return c.SendString("ok")
		})

		token, _ := GenerateToken("user-123", "test@example.com", []string{"customer"})
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("not authenticated", func(t *testing.T) {
		app := fiber.New()
		app.Get("/test", func(c *fiber.Ctx) error {
			user, err := GetUserFromContext(c)
			assert.Error(t, err)
			assert.Nil(t, user)
			return c.SendStatus(401)
		})

		req := httptest.NewRequest("GET", "/test", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 401, resp.StatusCode)
	})
}
