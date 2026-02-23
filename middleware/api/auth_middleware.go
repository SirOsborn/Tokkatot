package api

import (
	"strings"

	"middleware/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// AuthMiddleware validates JWT token from Authorization header
func AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return utils.Unauthorized(c, "Missing authorization header")
	}

	// Extract and validate JWT
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return utils.Unauthorized(c, "Invalid authorization header format")
	}

	tokenString := parts[1]
	claims, err := utils.ValidateToken(tokenString)
	if err != nil {
		return utils.Unauthorized(c, err.Error())
	}

	// Store claims in context for use in handlers
	c.Locals("user_id", claims.UserID)
	c.Locals("email", claims.Email)
	c.Locals("phone", claims.Phone)
	c.Locals("farm_id", claims.FarmID)
	c.Locals("role", claims.Role)

	return c.Next()
}

// RequireRole checks if user has the required role
// Usage: RequireRole("manager", "owner") for multiple allowed roles
func RequireRole(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("role")
		if userRole == nil {
			return utils.Forbidden(c, "User role not found in token")
		}

		role, ok := userRole.(string)
		if !ok {
			return utils.Forbidden(c, "Invalid user role")
		}

		// Check if user role is in allowed roles
		allowed := false
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				allowed = true
				break
			}
		}

		if !allowed {
			return utils.Forbidden(c, "Insufficient permissions for this operation")
		}

		return c.Next()
	}
}

// RequireOwner checks if user is farm owner
func RequireOwner(c *fiber.Ctx) error {
	return RequireRole("owner")(c)
}

// RequireManagerOrOwner checks if user is manager or owner
func RequireManagerOrOwner(c *fiber.Ctx) error {
	return RequireRole("manager", "owner")(c)
}

// RequireAnyRole allows any authenticated role
func RequireAnyRole(c *fiber.Ctx) error {
	userRole := c.Locals("role")
	if userRole == nil {
		return utils.Forbidden(c, "User role not found in token")
	}
	return c.Next()
}

// GetUserIDFromContext extracts user ID from context
func GetUserIDFromContext(c *fiber.Ctx) (uuid.UUID, error) {
	userIDInterface := c.Locals("user_id")
	if userIDInterface == nil {
		return uuid.Nil, fiber.NewError(fiber.StatusUnauthorized, "User ID not found in context")
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		return uuid.Nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid user ID in context")
	}

	return userID, nil
}

// GetFarmIDFromContext extracts farm ID from context
func GetFarmIDFromContext(c *fiber.Ctx) (uuid.UUID, error) {
	farmIDInterface := c.Locals("farm_id")
	if farmIDInterface == nil {
		return uuid.Nil, fiber.NewError(fiber.StatusUnauthorized, "Farm ID not found in context")
	}

	farmID, ok := farmIDInterface.(uuid.UUID)
	if !ok {
		return uuid.Nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid farm ID in context")
	}

	return farmID, nil
}

// GetRoleFromContext extracts role from context
func GetRoleFromContext(c *fiber.Ctx) string {
	roleInterface := c.Locals("role")
	if roleInterface == nil {
		return ""
	}

	role, ok := roleInterface.(string)
	if !ok {
		return ""
	}

	return role
}
