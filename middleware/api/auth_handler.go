package api

import (
	"database/sql"
	"log"
	"os"
	"time"

	"middleware/database"
	"middleware/models"
	"middleware/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// ===== AUTHENTICATION HANDLERS =====

// Signup handles user registration
// @POST /v1/auth/signup
func SignupHandler(c *fiber.Ctx) error {
	var req models.SignupRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid_request", "Invalid request body")
	}

	// Validate input
	if (req.Email == nil || *req.Email == "") && (req.Phone == nil || *req.Phone == "") {
		return utils.BadRequest(c, "missing_contact", "Email or phone number is required")
	}

	if req.Email != nil && *req.Email != "" {
		if err := utils.ValidateEmail(*req.Email); err != nil {
			return utils.BadRequest(c, "invalid_email", err.Error())
		}
	}

	if req.Phone != nil && *req.Phone != "" {
		if err := utils.ValidatePhone(*req.Phone); err != nil {
			return utils.BadRequest(c, "invalid_phone", err.Error())
		}
	}

	if err := utils.ValidateName(req.Name); err != nil {
		return utils.BadRequest(c, "invalid_name", err.Error())
	}

	if err := utils.ValidatePassword(req.Password); err != nil {
		return utils.BadRequest(c, "weak_password", err.Error())
	}

	// ===== REGISTRATION KEY VALIDATION (On-site setup) =====
	var isVerifiedByKey bool
	var regKeyID *uuid.UUID

	if req.RegistrationKey != nil && *req.RegistrationKey != "" {
		// Validate registration key
		var keyIsUsed bool
		var keyExpires *time.Time
		var storedKeyID uuid.UUID

		query := `
		SELECT id, is_used, expires_at 
		FROM registration_keys 
		WHERE key_code = $1
		`

		err := database.DB.QueryRow(query, req.RegistrationKey).Scan(&storedKeyID, &keyIsUsed, &keyExpires)

		if err == sql.ErrNoRows {
			return utils.BadRequest(c, "invalid_key", "Invalid registration key")
		}

		if err != nil {
			log.Printf("Registration key check error: %v", err)
			return utils.InternalError(c, "Failed to validate registration key")
		}

		// Check if already used
		if keyIsUsed {
			return utils.BadRequest(c, "key_used", "Registration key has already been used")
		}

		// Check expiry
		if keyExpires != nil && keyExpires.Before(time.Now()) {
			return utils.BadRequest(c, "key_expired", "Registration key has expired")
		}

		// Key is valid!
		isVerifiedByKey = true
		regKeyID = &storedKeyID
	}

	// Validate language
	language := "km" // Default to Khmer
	if req.Language != nil && *req.Language != "" {
		if err := utils.ValidateLanguage(*req.Language); err != nil {
			return utils.BadRequest(c, "invalid_language", err.Error())
		}
		language = *req.Language
	}

	// Check if email/phone already exists
	if req.Email != nil && *req.Email != "" {
		var exists bool
		err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", req.Email).Scan(&exists)
		if err != nil {
			log.Printf("Email check error: %v", err)
			return utils.InternalError(c, "Database error")
		}
		if exists {
			return utils.Conflict(c, "email_exists", "Email already registered")
		}
	}

	if req.Phone != nil && *req.Phone != "" {
		var exists bool
		err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE phone = $1)", req.Phone).Scan(&exists)
		if err != nil {
			log.Printf("Phone check error: %v", err)
			return utils.InternalError(c, "Database error")
		}
		if exists {
			return utils.Conflict(c, "phone_exists", "Phone number already registered")
		}
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Password hash error: %v", err)
		return utils.InternalError(c, "Failed to process password")
	}

	// Create user
	userID := uuid.New()
	var verificationType string
	if req.Email != nil && *req.Email != "" {
		verificationType = "email"
	} else if req.Phone != nil && *req.Phone != "" {
		verificationType = "phone"
	}

	query := `
	INSERT INTO users (id, email, phone, password_hash, name, language, timezone, contact_verified, verification_type, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	RETURNING id, email, phone, name, language
	`

	// Auto-verify if:
	// 1. Development mode, OR
	// 2. Valid registration key provided (on-site setup)
	isVerified := os.Getenv("ENVIRONMENT") == "development" || isVerifiedByKey

	var user models.User
	err = database.DB.QueryRow(
		query,
		userID,
		req.Email,
		req.Phone,
		string(hashedPassword),
		req.Name,
		language,
		"Asia/Phnom_Penh", // Default timezone
		isVerified,        // Auto-verify in development
		verificationType,
	).Scan(&user.ID, &user.Email, &user.Phone, &user.Name, &user.Language)

	if err != nil {
		log.Printf("User creation error: %v", err)
		return utils.InternalError(c, "Failed to create user")
	}

	// Mark registration key as used (if provided)
	if regKeyID != nil {
		_, err := database.DB.Exec(`
		UPDATE registration_keys 
		SET is_used = true, used_by_user_id = $1, used_at = CURRENT_TIMESTAMP
		WHERE id = $2
		`, user.ID, regKeyID)

		if err != nil {
			log.Printf("Failed to mark registration key as used: %v", err)
			// Don't fail the signup, just log the error
		}
	}

	// TODO: Send verification email/SMS (only if not verified by key)

	// Return success response
	return utils.SuccessResponse(c, fiber.StatusCreated, fiber.Map{
		"user_id": user.ID,
		"message": "Verification code sent to " + getContactMethodFromSignup(req),
	}, "User registered successfully")
}

// Login handles user login and generates tokens
// @POST /v1/auth/login
func LoginHandler(c *fiber.Ctx) error {
	var req models.LoginRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid_request", "Invalid request body")
	}

	// Validate input
	if (req.Email == nil || *req.Email == "") && (req.Phone == nil || *req.Phone == "") {
		return utils.BadRequest(c, "missing_contact", "Email or phone number is required")
	}

	if req.Password == "" {
		return utils.BadRequest(c, "missing_password", "Password is required")
	}

	// Lookup user by email or phone
	var user models.User
	var storedHash string

	query := `
	SELECT id, email, phone, password_hash, name, language, is_active, contact_verified
	FROM users
	WHERE ($1::TEXT IS NULL OR email = $1) OR ($2::TEXT IS NULL OR phone = $2)
	`

	err := database.DB.QueryRow(query, req.Email, req.Phone).Scan(
		&user.ID,
		&user.Email,
		&user.Phone,
		&storedHash,
		&user.Name,
		&user.Language,
		&user.IsActive,
		&user.ContactVerified,
	)

	if err == sql.ErrNoRows {
		return utils.Unauthorized(c, "Invalid credentials")
	}

	if err != nil {
		log.Printf("Login query error: %v", err)
		return utils.InternalError(c, "Database error")
	}

	// Check if user is active and verified
	if !user.IsActive {
		return utils.Unauthorized(c, "Account is inactive")
	}

	if !user.ContactVerified {
		return utils.Unauthorized(c, "Contact not verified. Please verify before logging in")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(req.Password)); err != nil {
		// Log failed attempt (TODO: implement rate limiting)
		return utils.Unauthorized(c, "Invalid credentials")
	}

	// Get user's first farm (for JWT claims)
	var farmID uuid.UUID
	var role string
	err = database.DB.QueryRow(`
	SELECT farm_id, role FROM farm_users
	WHERE user_id = $1 AND is_active = true
	LIMIT 1
	`, user.ID).Scan(&farmID, &role)

	if err == sql.ErrNoRows {
		// User has no farm yet - create default farm for first login
		farmID = uuid.New()
		role = "owner"

		_, err := database.DB.Exec(`
		INSERT INTO farms (id, owner_id, name, timezone, created_at, updated_at)
		VALUES ($1, $2, $3, 'Asia/Phnom_Penh', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		`, farmID, user.ID)

		if err != nil {
			log.Printf("Farm creation error: %v", err)
			return utils.InternalError(c, "Failed to create default farm")
		}

		// Add user to farm
		_, err = database.DB.Exec(`
		INSERT INTO farm_users (id, farm_id, user_id, role, invited_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		`, uuid.New(), farmID, user.ID, role)

		if err != nil {
			log.Printf("Farm user creation error: %v", err)
			return utils.InternalError(c, "Failed to add user to farm")
		}
	} else if err != nil {
		log.Printf("Farm lookup error: %v", err)
		return utils.InternalError(c, "Database error")
	}

	// Generate tokens
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Email, user.Phone, farmID, role)
	if err != nil {
		log.Printf("Access token generation error: %v", err)
		return utils.InternalError(c, "Failed to generate token")
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		log.Printf("Refresh token generation error: %v", err)
		return utils.InternalError(c, "Failed to generate token")
	}

	// Update last login
	_, err = database.DB.Exec("UPDATE users SET last_login = CURRENT_TIMESTAMP WHERE id = $1", user.ID)
	if err != nil {
		log.Printf("Last login update error: %v", err)
	}

	// Return tokens
	response := models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(utils.AccessTokenExpiry.Seconds()),
		User: models.UserInfo{
			ID:       user.ID,
			Name:     user.Name,
			Email:    user.Email,
			Phone:    user.Phone,
			Role:     role,
			Language: user.Language,
		},
	}

	return utils.SuccessResponse(c, fiber.StatusOK, response, "Login successful")
}

// RefreshTokenHandler handles token refresh
// @POST /v1/auth/refresh
func RefreshTokenHandler(c *fiber.Ctx) error {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid_request", "Invalid request body")
	}

	if req.RefreshToken == "" {
		return utils.BadRequest(c, "missing_token", "Refresh token is required")
	}

	// Validate refresh token
	userID, err := utils.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		return utils.Unauthorized(c, "Invalid refresh token")
	}

	// Get user details
	var user models.User
	err = database.DB.QueryRow(`
	SELECT id, email, phone, name, language
	FROM users WHERE id = $1
	`, userID).Scan(&user.ID, &user.Email, &user.Phone, &user.Name, &user.Language)

	if err != nil {
		return utils.Unauthorized(c, "User not found")
	}

	// Get user's farm
	var farmID uuid.UUID
	var role string
	err = database.DB.QueryRow(`
	SELECT farm_id, role FROM farm_users
	WHERE user_id = $1 AND is_active = true
	LIMIT 1
	`, user.ID).Scan(&farmID, &role)

	if err != nil {
		return utils.Unauthorized(c, "No active farm found")
	}

	// Generate new tokens
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Email, user.Phone, farmID, role)
	if err != nil {
		return utils.InternalError(c, "Failed to generate token")
	}

	newRefreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return utils.InternalError(c, "Failed to generate token")
	}

	response := models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int64(utils.AccessTokenExpiry.Seconds()),
		User: models.UserInfo{
			ID:       user.ID,
			Name:     user.Name,
			Email:    user.Email,
			Phone:    user.Phone,
			Role:     role,
			Language: user.Language,
		},
	}

	return utils.SuccessResponse(c, fiber.StatusOK, response, "Token refreshed")
}

// LogoutHandler handles user logout
// @POST /v1/auth/logout
func LogoutHandler(c *fiber.Ctx) error {
	// JWT is stateless, so logout just requires client to discard token
	// Could add token to blacklist for more control (future enhancement)
	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Logged out successfully")
}

// ForgotPasswordHandler initiates password reset
// @POST /v1/auth/forgot-password
func ForgotPasswordHandler(c *fiber.Ctx) error {
	var req struct {
		Email *string `json:"email,omitempty"`
		Phone *string `json:"phone,omitempty"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid_request", "Invalid request body")
	}

	if (req.Email == nil || *req.Email == "") && (req.Phone == nil || *req.Phone == "") {
		return utils.BadRequest(c, "missing_contact", "Email or phone is required")
	}

	// Check if user exists
	var userID uuid.UUID
	if req.Email != nil && *req.Email != "" {
		err := database.DB.QueryRow("SELECT id FROM users WHERE email = $1", req.Email).Scan(&userID)
		if err == sql.ErrNoRows {
			// For security, don't reveal if email exists
			return utils.SuccessResponse(c, fiber.StatusOK, nil, "If email exists, password reset code sent")
		}
		if err != nil {
			return utils.InternalError(c, "Database error")
		}
	}

	if req.Phone != nil && *req.Phone != "" {
		err := database.DB.QueryRow("SELECT id FROM users WHERE phone = $1", req.Phone).Scan(&userID)
		if err == sql.ErrNoRows {
			return utils.SuccessResponse(c, fiber.StatusOK, nil, "If phone exists, password reset code sent")
		}
		if err != nil {
			return utils.InternalError(c, "Database error")
		}
	}

	// TODO: Generate reset code and send via email/SMS
	// TODO: Store reset code with expiry (30 minutes)

	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Password reset code sent")
}

// ResetPasswordHandler completes password reset
// @POST /v1/auth/reset-password
func ResetPasswordHandler(c *fiber.Ctx) error {
	var req struct {
		ResetCode   string `json:"reset_code"`
		NewPassword string `json:"new_password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid_request", "Invalid request body")
	}

	if req.ResetCode == "" {
		return utils.BadRequest(c, "missing_code", "Reset code is required")
	}

	if err := utils.ValidatePassword(req.NewPassword); err != nil {
		return utils.BadRequest(c, "weak_password", err.Error())
	}

	// TODO: Verify reset code and expiry
	// TODO: Update password and invalidate reset code

	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Password reset successfully")
}

// VerifyContact handles email/phone verification
// @POST /v1/auth/verify
func VerifyContactHandler(c *fiber.Ctx) error {
	var req struct {
		Email *string `json:"email,omitempty"`
		Phone *string `json:"phone,omitempty"`
		Code  string  `json:"code"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid_request", "Invalid request body")
	}

	if (req.Email == nil || *req.Email == "") && (req.Phone == nil || *req.Phone == "") {
		return utils.BadRequest(c, "missing_contact", "Email or phone is required")
	}

	if req.Code == "" {
		return utils.BadRequest(c, "missing_code", "Verification code is required")
	}

	// In development: Accept any 6-digit code for testing
	if os.Getenv("ENVIRONMENT") == "development" && len(req.Code) == 6 {
		query := `
		UPDATE users 
		SET contact_verified = true, updated_at = CURRENT_TIMESTAMP
		WHERE ($1::TEXT IS NULL OR email = $1) OR ($2::TEXT IS NULL OR phone = $2)
		RETURNING id, email, phone, name
		`

		var user models.User
		err := database.DB.QueryRow(query, req.Email, req.Phone).Scan(
			&user.ID,
			&user.Email,
			&user.Phone,
			&user.Name,
		)

		if err == sql.ErrNoRows {
			return utils.BadRequest(c, "user_not_found", "User not found")
		}

		if err != nil {
			log.Printf("Verification error: %v", err)
			return utils.InternalError(c, "Failed to verify contact")
		}

		return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
			"user_id": user.ID,
		}, "Contact verified successfully")
	}

	// Production: Verify code from database
	// TODO: Check verification_codes table for matching code and expiry
	// query := `
	// SELECT user_id FROM verification_codes
	// WHERE code = $1 AND expires_at > CURRENT_TIMESTAMP
	// AND (email = $2 OR phone = $3)
	// `

	return utils.BadRequest(c, "invalid_code", "Invalid or expired verification code")
}

// ===== HELPER FUNCTIONS =====

func getContactMethodFromSignup(req models.SignupRequest) string {
	if req.Email != nil && *req.Email != "" {
		return *req.Email
	}
	if req.Phone != nil && *req.Phone != "" {
		return *req.Phone
	}
	return "provided contact"
}

func getContactMethod(req models.LoginRequest) string {
	if req.Email != nil && *req.Email != "" {
		return *req.Email
	}
	if req.Phone != nil && *req.Phone != "" {
		return *req.Phone
	}
	return "provided contact"
}
