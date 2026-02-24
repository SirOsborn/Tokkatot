package api

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"time"

	"middleware/database"
	"middleware/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// AdminMiddleware ensures the requester has role="admin"
func AdminMiddleware(c *fiber.Ctx) error {
	role, ok := c.Locals("role").(string)
	if !ok || role != "admin" {
		return utils.Forbidden(c, "Admin access required")
	}
	return c.Next()
}

// generateRegKey creates a formatted registration key like "A3F7-B2C9"
func generateRegKey() (string, error) {
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	s := strings.ToUpper(hex.EncodeToString(b))
	return fmt.Sprintf("%s-%s", s[:4], s[4:]), nil
}

// ===== GET /v1/admin/stats =====
// GetAdminStatsHandler returns system-level stats
func GetAdminStatsHandler(c *fiber.Ctx) error {
	stats := fiber.Map{}

	var farmers, viewers, farms, regKeysUsed, regKeysUnused, admins int
	database.DB.QueryRow(`SELECT COUNT(*) FROM farm_users WHERE role='farmer' AND is_active=true`).Scan(&farmers)
	database.DB.QueryRow(`SELECT COUNT(*) FROM farm_users WHERE role='viewer' AND is_active=true`).Scan(&viewers)
	database.DB.QueryRow(`SELECT COUNT(*) FROM farms WHERE is_active=true`).Scan(&farms)
	database.DB.QueryRow(`SELECT COUNT(*) FROM registration_keys WHERE is_used=false AND (expires_at IS NULL OR expires_at > NOW())`).Scan(&regKeysUnused)
	database.DB.QueryRow(`SELECT COUNT(*) FROM registration_keys WHERE is_used=true`).Scan(&regKeysUsed)
	database.DB.QueryRow(`SELECT COUNT(*) FROM admins WHERE is_active=true`).Scan(&admins)

	stats["farmers"] = farmers
	stats["viewers"] = viewers
	stats["farms"] = farms
	stats["reg_keys_unused"] = regKeysUnused
	stats["reg_keys_used"] = regKeysUsed
	stats["admins"] = admins

	return utils.SuccessResponse(c, fiber.StatusOK, stats, "Stats fetched")
}

// ===== GET /v1/admin/farmers =====
// ListFarmersHandler returns all farmers with their profiles and farm info
func ListFarmersHandler(c *fiber.Ctx) error {
	rows, err := database.DB.Query(`
	SELECT
		u.id, u.name, u.email, u.phone, u.language, u.is_active, u.contact_verified,
		u.created_at, u.last_login,
		f.id AS farm_id, f.name AS farm_name, f.province,
		fp.national_id_number, fp.full_name_kh, fp.full_name_en, fp.sex,
		fp.province AS profile_province, fp.district,
		rk.key_code AS reg_key, rk.is_used AS reg_key_used
	FROM users u
	JOIN farm_users fu ON fu.user_id = u.id AND fu.role = 'farmer' AND fu.is_active = true
	JOIN farms f ON f.id = fu.farm_id
	LEFT JOIN farmer_profiles fp ON fp.user_id = u.id
	LEFT JOIN registration_keys rk ON rk.used_by_user_id = u.id
	ORDER BY u.created_at DESC
	`)
	if err != nil {
		log.Printf("ListFarmers query error: %v", err)
		return utils.InternalError(c, "Failed to fetch farmers")
	}
	defer rows.Close()

	type FarmerRow struct {
		ID               uuid.UUID  `json:"id"`
		Name             string     `json:"name"`
		Email            *string    `json:"email,omitempty"`
		Phone            *string    `json:"phone,omitempty"`
		Language         *string    `json:"language,omitempty"`
		IsActive         bool       `json:"is_active"`
		ContactVerified  bool       `json:"contact_verified"`
		CreatedAt        time.Time  `json:"created_at"`
		LastLogin        *time.Time `json:"last_login,omitempty"`
		FarmID           uuid.UUID  `json:"farm_id"`
		FarmName         string     `json:"farm_name"`
		Province         *string    `json:"province,omitempty"`
		NationalIDNumber *string    `json:"national_id_number,omitempty"`
		FullNameKh       *string    `json:"full_name_kh,omitempty"`
		FullNameEn       *string    `json:"full_name_en,omitempty"`
		Sex              *string    `json:"sex,omitempty"`
		ProfileProvince  *string    `json:"profile_province,omitempty"`
		District         *string    `json:"district,omitempty"`
		RegKey           *string    `json:"reg_key,omitempty"`
		RegKeyUsed       *bool      `json:"reg_key_used,omitempty"`
	}

	var results []FarmerRow
	for rows.Next() {
		var r FarmerRow
		if err := rows.Scan(
			&r.ID, &r.Name, &r.Email, &r.Phone, &r.Language, &r.IsActive, &r.ContactVerified,
			&r.CreatedAt, &r.LastLogin,
			&r.FarmID, &r.FarmName, &r.Province,
			&r.NationalIDNumber, &r.FullNameKh, &r.FullNameEn, &r.Sex,
			&r.ProfileProvince, &r.District,
			&r.RegKey, &r.RegKeyUsed,
		); err != nil {
			log.Printf("ListFarmers scan error: %v", err)
			continue
		}
		results = append(results, r)
	}
	if results == nil {
		results = []FarmerRow{}
	}

	return utils.SuccessResponse(c, fiber.StatusOK, results, "Farmers fetched")
}

// ===== POST /v1/admin/farmers =====
// RegisterFarmerHandler creates a new farmer account during on-site setup
// Required: name, phone (or email), password
// Optional: national_id_number, full_name_kh, full_name_en, sex, province, district, notes
func RegisterFarmerHandler(c *fiber.Ctx) error {
	adminID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	var req struct {
		Name             string  `json:"name"`
		Email            *string `json:"email,omitempty"`
		Phone            *string `json:"phone,omitempty"`
		Password         string  `json:"password"`
		Language         *string `json:"language,omitempty"`
		NationalIDNumber *string `json:"national_id_number,omitempty"`
		FullNameKh       *string `json:"full_name_kh,omitempty"`
		FullNameEn       *string `json:"full_name_en,omitempty"`
		Sex              *string `json:"sex,omitempty"`
		Province         *string `json:"province,omitempty"`
		District         *string `json:"district,omitempty"`
		Notes            *string `json:"notes,omitempty"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid_request", "Invalid request body")
	}

	if req.Name == "" {
		return utils.BadRequest(c, "missing_name", "Name is required")
	}
	if req.Password == "" {
		return utils.BadRequest(c, "missing_password", "Password is required")
	}
	if (req.Email == nil || *req.Email == "") && (req.Phone == nil || *req.Phone == "") {
		return utils.BadRequest(c, "missing_contact", "Phone or email is required")
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return utils.InternalError(c, "Failed to hash password")
	}

	lang := "km"
	if req.Language != nil && *req.Language != "" {
		lang = *req.Language
	}

	// Create user — admin-registered farmers are pre-verified
	userID := uuid.New()
	_, err = database.DB.Exec(`
	INSERT INTO users (id, name, email, phone, password_hash, language, is_active, contact_verified, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, true, true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, userID, req.Name, req.Email, req.Phone, string(hash), lang)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			return utils.BadRequest(c, "duplicate_contact", "Phone or email already registered")
		}
		log.Printf("RegisterFarmer user insert error: %v", err)
		return utils.InternalError(c, "Failed to create user")
	}

	// Auto-generate registration key for this farmer
	regKey, err := generateRegKey()
	if err != nil {
		return utils.InternalError(c, "Failed to generate registration key")
	}

	// Ensure key is unique (retry once on collision)
	var existing int
	database.DB.QueryRow(`SELECT COUNT(*) FROM registration_keys WHERE key_code = $1`, regKey).Scan(&existing)
	if existing > 0 {
		regKey, _ = generateRegKey()
	}

	regKeyID := uuid.New()
	expires := time.Now().AddDate(5, 0, 0) // 5 years
	farmName := req.Name + "'s Farm"
	farmLocation := ""
	if req.Province != nil {
		farmLocation = *req.Province
	}
	custName := req.Name
	custPhone := ""
	if req.Phone != nil {
		custPhone = *req.Phone
	}
	_, err = database.DB.Exec(`
	INSERT INTO registration_keys (id, key_code, farm_name, farm_location, customer_name, customer_phone, is_used, used_by_user_id, used_at, expires_at, created_by, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, true, $7, CURRENT_TIMESTAMP, $8, 'admin', CURRENT_TIMESTAMP)
	`, regKeyID, regKey, farmName, farmLocation, custName, custPhone, userID, expires)
	if err != nil {
		log.Printf("RegisterFarmer reg key insert error: %v", err)
		// Non-fatal — continue
	}

	// Create farmer_profile if any profile data provided
	hasProfile := req.NationalIDNumber != nil || req.FullNameKh != nil || req.FullNameEn != nil ||
		req.Sex != nil || req.Province != nil || req.District != nil || req.Notes != nil
	if hasProfile {
		profileID := uuid.New()
		_, perr := database.DB.Exec(`
		INSERT INTO farmer_profiles
		  (id, user_id, national_id_number, full_name_kh, full_name_en, sex, province, district, notes, created_by_admin, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		`, profileID, userID, req.NationalIDNumber, req.FullNameKh, req.FullNameEn, req.Sex, req.Province, req.District, req.Notes, adminID)
		if perr != nil {
			log.Printf("RegisterFarmer profile insert error: %v", perr)
		}
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, fiber.Map{
		"user_id":  userID,
		"name":     req.Name,
		"phone":    req.Phone,
		"email":    req.Email,
		"reg_key":  regKey,
		"language": lang,
		"message":  "Farmer registered. They can login with phone+password.",
	}, "Farmer registered successfully")
}

// ===== GET /v1/admin/reg-keys =====
// ListRegKeysHandler returns all registration keys
func ListRegKeysHandler(c *fiber.Ctx) error {
	rows, err := database.DB.Query(`
	SELECT rk.id, rk.key_code, rk.farm_name, rk.customer_name, rk.customer_phone,
	       rk.is_used, rk.used_at, rk.expires_at, rk.created_by, rk.created_at,
	       u.name AS used_by_name, u.phone AS used_by_phone
	FROM registration_keys rk
	LEFT JOIN users u ON u.id = rk.used_by_user_id
	ORDER BY rk.created_at DESC
	`)
	if err != nil {
		log.Printf("ListRegKeys error: %v", err)
		return utils.InternalError(c, "Failed to fetch registration keys")
	}
	defer rows.Close()

	type RegKeyRow struct {
		ID            uuid.UUID  `json:"id"`
		KeyCode       string     `json:"key_code"`
		FarmName      *string    `json:"farm_name,omitempty"`
		CustomerName  *string    `json:"customer_name,omitempty"`
		CustomerPhone *string    `json:"customer_phone,omitempty"`
		IsUsed        bool       `json:"is_used"`
		UsedAt        *time.Time `json:"used_at,omitempty"`
		ExpiresAt     *time.Time `json:"expires_at,omitempty"`
		CreatedBy     *string    `json:"created_by,omitempty"`
		CreatedAt     time.Time  `json:"created_at"`
		UsedByName    *string    `json:"used_by_name,omitempty"`
		UsedByPhone   *string    `json:"used_by_phone,omitempty"`
	}

	var results []RegKeyRow
	for rows.Next() {
		var r RegKeyRow
		if err := rows.Scan(
			&r.ID, &r.KeyCode, &r.FarmName, &r.CustomerName, &r.CustomerPhone,
			&r.IsUsed, &r.UsedAt, &r.ExpiresAt, &r.CreatedBy, &r.CreatedAt,
			&r.UsedByName, &r.UsedByPhone,
		); err != nil {
			log.Printf("ListRegKeys scan error: %v", err)
			continue
		}
		results = append(results, r)
	}
	if results == nil {
		results = []RegKeyRow{}
	}
	return utils.SuccessResponse(c, fiber.StatusOK, results, "Registration keys fetched")
}

// ===== GET /v1/admin/viewers =====
// ListViewersHandler returns all viewers (workers) with their farm info
func ListViewersHandler(c *fiber.Ctx) error {
	rows, err := database.DB.Query(`
	SELECT u.id, u.name, u.email, u.phone, u.language, u.is_active, u.created_at,
	       f.id AS farm_id, f.name AS farm_name,
	       fu2.name AS farmer_name
	FROM users u
	JOIN farm_users fv ON fv.user_id = u.id AND fv.role = 'viewer' AND fv.is_active = true
	JOIN farms f ON f.id = fv.farm_id
	LEFT JOIN farm_users ffu ON ffu.farm_id = f.id AND ffu.role = 'farmer' AND ffu.is_active = true
	LEFT JOIN users fu2 ON fu2.id = ffu.user_id
	ORDER BY u.created_at DESC
	`)
	if err != nil {
		log.Printf("ListViewers error: %v", err)
		return utils.InternalError(c, "Failed to fetch viewers")
	}
	defer rows.Close()

	type ViewerRow struct {
		ID         uuid.UUID `json:"id"`
		Name       string    `json:"name"`
		Email      *string   `json:"email,omitempty"`
		Phone      *string   `json:"phone,omitempty"`
		Language   *string   `json:"language,omitempty"`
		IsActive   bool      `json:"is_active"`
		CreatedAt  time.Time `json:"created_at"`
		FarmID     uuid.UUID `json:"farm_id"`
		FarmName   string    `json:"farm_name"`
		FarmerName *string   `json:"farmer_name,omitempty"`
	}

	var results []ViewerRow
	for rows.Next() {
		var r ViewerRow
		if serr := rows.Scan(
			&r.ID, &r.Name, &r.Email, &r.Phone, &r.Language, &r.IsActive, &r.CreatedAt,
			&r.FarmID, &r.FarmName, &r.FarmerName,
		); serr != nil {
			log.Printf("ListViewers scan error: %v", serr)
			continue
		}
		results = append(results, r)
	}
	if results == nil {
		results = []ViewerRow{}
	}
	return utils.SuccessResponse(c, fiber.StatusOK, results, "Viewers fetched")
}

// ===== PUT /v1/admin/profile =====
// UpdateAdminProfileHandler lets admin update their own name/language
func UpdateAdminProfileHandler(c *fiber.Ctx) error {
	adminID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}
	var req struct {
		Name     *string `json:"name,omitempty"`
		Language *string `json:"language,omitempty"`
	}
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid_request", "Invalid request body")
	}
	_, err := database.DB.Exec(`
	UPDATE admins SET
		name = COALESCE($1, name),
		language = COALESCE($2, language)
	WHERE id = $3
	`, req.Name, req.Language, adminID)
	if err != nil {
		log.Printf("UpdateAdminProfile error: %v", err)
		return utils.InternalError(c, "Failed to update profile")
	}
	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Profile updated")
}

// ===== DELETE /v1/admin/farmers/:user_id =====
// DeactivateFarmerHandler deactivates a farmer account
func DeactivateFarmerHandler(c *fiber.Ctx) error {
	farmerIDStr := c.Params("user_id")
	farmerID, err := uuid.Parse(farmerIDStr)
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid user ID")
	}
	_, err = database.DB.Exec(`UPDATE users SET is_active = false WHERE id = $1`, farmerID)
	if err != nil {
		return utils.InternalError(c, "Failed to deactivate farmer")
	}
	return utils.SuccessResponse(c, fiber.StatusOK, nil, "Farmer deactivated")
}

// ===== GET /v1/admin/farmer-profiles/:user_id =====
// GetFarmerProfileHandler returns detailed profile for a specific farmer
func GetFarmerProfileHandler(c *fiber.Ctx) error {
	farmerIDStr := c.Params("user_id")
	farmerID, err := uuid.Parse(farmerIDStr)
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid user ID")
	}

	type Profile struct {
		UserID           uuid.UUID `json:"user_id"`
		NationalIDNumber *string   `json:"national_id_number,omitempty"`
		FullNameKh       *string   `json:"full_name_kh,omitempty"`
		FullNameEn       *string   `json:"full_name_en,omitempty"`
		Sex              *string   `json:"sex,omitempty"`
		Province         *string   `json:"province,omitempty"`
		District         *string   `json:"district,omitempty"`
		Notes            *string   `json:"notes,omitempty"`
	}
	var p Profile
	p.UserID = farmerID
	serr := database.DB.QueryRow(`
	SELECT national_id_number, full_name_kh, full_name_en, sex, province, district, notes
	FROM farmer_profiles WHERE user_id = $1
	`, farmerID).Scan(&p.NationalIDNumber, &p.FullNameKh, &p.FullNameEn, &p.Sex, &p.Province, &p.District, &p.Notes)
	if serr == sql.ErrNoRows {
		return utils.SuccessResponse(c, fiber.StatusOK, p, "No profile data yet")
	}
	if serr != nil {
		return utils.InternalError(c, "Failed to fetch profile")
	}
	return utils.SuccessResponse(c, fiber.StatusOK, p, "Profile fetched")
}
