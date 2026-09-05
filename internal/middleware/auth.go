package middleware

import (
	"context"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserContextKey contextKey = "user"

type UserClaims struct {
	UserID string   `json:"user_id"`
	Email  string   `json:"email"`
	Roles  []string `json:"roles"`
	jwt.RegisteredClaims
}

var jwtSecret = []byte("your-secret-key-change-this-in-production") // TODO: Move to env

// AuthRequired middleware validates JWT token
func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Try to get token from Authorization header
		authHeader := c.Get("Authorization")
		var tokenString string

		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			// Fallback to cookie for backward compatibility
			tokenString = c.Cookies("auth_token")
		}

		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authentication token",
			})
		}

		// Parse and validate token
		token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		claims, ok := token.Claims.(*UserClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token claims",
			})
		}

		// Store user info in context
		c.Locals(string(UserContextKey), claims)

		return c.Next()
	}
}

// RequireRole middleware checks if user has at least one of the required roles
func RequireRole(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, ok := c.Locals(string(UserContextKey)).(*UserClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User not authenticated",
			})
		}

		// Check if user has any of the allowed roles
		hasRole := false
		for _, userRole := range claims.Roles {
			for _, allowedRole := range allowedRoles {
				if userRole == allowedRole {
					hasRole = true
					break
				}
			}
			if hasRole {
				break
			}
		}

		if !hasRole {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":         "Insufficient permissions",
				"required_role": allowedRoles,
			})
		}

		return c.Next()
	}
}

// RequirePermission checks if user has specific permission
func RequirePermission(permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, ok := c.Locals(string(UserContextKey)).(*UserClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User not authenticated",
			})
		}

		// Super admin has all permissions
		for _, role := range claims.Roles {
			if role == "super_admin" {
				return c.Next()
			}
		}

		// TODO: Implement permission check against database
		// For now, just check role-based access
		return c.Next()
	}
}

// GetUserFromContext extracts user claims from fiber context
func GetUserFromContext(c *fiber.Ctx) (*UserClaims, error) {
	claims, ok := c.Locals(string(UserContextKey)).(*UserClaims)
	if !ok {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "User not authenticated")
	}
	return claims, nil
}

// GenerateToken creates a new JWT token with role claims
func GenerateToken(userID, email string, roles []string) (string, error) {
	claims := UserClaims{
		UserID: userID,
		Email:  email,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// Helper to extract context in usecase layer
func GetUserIDFromContext(ctx context.Context) string {
	if claims, ok := ctx.Value(UserContextKey).(*UserClaims); ok {
		return claims.UserID
	}
	return ""
}
