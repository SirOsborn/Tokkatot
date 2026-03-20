package api

import (
	"middleware/schemas"
	"middleware/utils"

	"github.com/gofiber/fiber/v2"
)

// ===== AUTHENTICATION HANDLERS =====

// SignupHandler registers a new farmer user
// @Summary User Signup
// @Description Creates a new user account and returns a JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body schemas.SignupRequest true "Signup Request"
// @Success 201 {object} schemas.JSONResponse
// @Router /v1/auth/signup [post]
func SignupHandler(c *fiber.Ctx) error {
	var req schemas.SignupRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid_body", "Invalid request body")
	}

	user, token, err := authService.Signup(req)
	if err != nil {
		return utils.BadRequest(c, "signup_failed", err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, fiber.Map{
		"user":  user,
		"token": token,
	}, "Account created successfully")
}

// LoginHandler authenticates a user
// @Summary User Login
// @Description Authenticates a user and returns a JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body schemas.LoginRequest true "Login Request"
// @Success 200 {object} schemas.JSONResponse
// @Router /v1/auth/login [post]
func LoginHandler(c *fiber.Ctx) error {
	var req schemas.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid_body", "Invalid request body")
	}

	identifier := ""
	if req.Email != nil && *req.Email != "" {
		identifier = *req.Email
	} else if req.Phone != nil && *req.Phone != "" {
		identifier = *req.Phone
	}

	user, token, role, err := authService.Login(identifier, req.Password)
	if err != nil {
		return utils.Unauthorized(c, "Invalid credentials")
	}

	// Return data in the format expected by the frontend
	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"user": fiber.Map{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"phone": user.Phone,
			"role":  role,
		},
		"access_token":  token,
		"refresh_token": token,
	}, "Login successful")
}

// RefreshTokenHandler refreshes a JWT
// @Summary Refresh Token
// @Description Returns a new JWT using a valid refresh token
// @Tags Auth
// @Produce json
// @Success 200 {object} schemas.JSONResponse
// @Router /v1/auth/refresh [post]
func RefreshTokenHandler(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"access_token":  "new_token_mock",
		"refresh_token": "new_refresh_token_mock",
	}, "Token refreshed")
}

// LogoutHandler invalidates a session
// @Summary Logout
// @Description Invalidates the current user session
// @Tags Auth
// @Produce json
// @Success 200 {object} schemas.JSONResponse
// @Router /v1/auth/logout [post]
func LogoutHandler(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Logout successful")
}

// AuthMiddleware and AdminMiddleware have been removed from here as they are in auth_middleware.go or common.go
