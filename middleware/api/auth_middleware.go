package api

import (
	"database/sql"
	"middleware/database"
	"middleware/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

// AuthMiddleware validates JWT token or Gateway API Key
func AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	gatewayHeader := c.Get("X-Gateway-Token")
	tokenString := ""

	// 1. Check for dedicated Gateway Header (Priority for IoT)
	if gatewayHeader != "" {
		return validateGatewayToken(c, gatewayHeader)
	}

	if authHeader == "" {
		// Allow websocket auth via query parameter (browser limitation)
		if websocket.IsWebSocketUpgrade(c) {
			tokenString = c.Query("token")
		} else {
			return utils.Unauthorized(c, "Missing authorization header")
		}
	} else {
		// Check for "Gateway <token>" in Authorization header
		if strings.HasPrefix(authHeader, "Gateway ") {
			return validateGatewayToken(c, authHeader[8:])
		}

		// Extract and validate JWT
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			if websocket.IsWebSocketUpgrade(c) {
				tokenString = authHeader // raw token in some WS clients
			} else {
				return utils.Unauthorized(c, "Invalid authorization header format")
			}
		} else {
			tokenString = parts[1]
		}
	}

	if tokenString == "" {
		return utils.Unauthorized(c, "Missing token")
	}

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
	c.Locals("auth_type", "jwt")

	// Hardened Deactivation Check: Verify user is still active in real-time
	isActive, err := adminService.IsUserActive(claims.UserID)
	if err != nil {
		return utils.Unauthorized(c, "Authentication failed: could not verify account status")
	}
	if !isActive {
		return utils.Unauthorized(c, "This account has been deactivated. Access denied.")
	}

	return c.Next()
}

// validateGatewayToken checks token against database and populates context
func validateGatewayToken(c *fiber.Ctx, token string) error {
	// Look up token in database (using a hash for security)
	tokenHash := utils.HashToken(token)

	var farmID, userID uuid.UUID
	err := database.DB.QueryRow(`
		UPDATE gateway_tokens 
		SET last_used_at = CURRENT_TIMESTAMP 
		WHERE token_hash = $1 AND is_active = true
		RETURNING farm_id, user_id
	`, tokenHash).Scan(&farmID, &userID)

	if err == sql.ErrNoRows {
		return utils.Unauthorized(c, "Invalid or inactive gateway token")
	}
	if err != nil {
		return utils.InternalError(c, "Token validation failed")
	}

	// Gateway tokens always have 'worker' permissions for their farm
	c.Locals("user_id", userID)
	c.Locals("farm_id", farmID)
	c.Locals("role", "worker")
	c.Locals("auth_type", "gateway")

	return c.Next()
}

// RequireRole checks if user has the required role
// Usage: RequireRole("farmer") or RequireRole("farmer", "viewer") for multiple allowed roles
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

// RequireOwner checks if user is a farmer (full farm access)
func RequireOwner(c *fiber.Ctx) error {
	return RequireRole("farmer")(c)
}

// RequireManagerOrOwner checks if user has farmer role (full farm access)
func RequireManagerOrOwner(c *fiber.Ctx) error {
	return RequireRole("farmer")(c)
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

// AdminMiddleware ensures the user has administrator privileges
func AdminMiddleware(c *fiber.Ctx) error {
	role := GetRoleFromContext(c)
	if role != "admin" {
		return utils.Forbidden(c, "Administrator privileges required")
	}

	// Verify admin is still active
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return utils.Unauthorized(c, "Authentication failed")
	}

	isActive, err := adminService.IsUserActive(userID)
	if err != nil || !isActive {
		return utils.Unauthorized(c, "Administrator account deactivated or inaccessible")
	}

	return c.Next()
}
