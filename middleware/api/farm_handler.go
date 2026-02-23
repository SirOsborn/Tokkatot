package api

import (
	"database/sql"
	"log"
	"strconv"

	"middleware/database"
	"middleware/models"
	"middleware/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ===== FARM MANAGEMENT HANDLERS =====

// ListFarmsHandler returns all farms the current user has access to
// @GET /v1/farms?limit=20&offset=0
func ListFarmsHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	// Pagination
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	if limit > 100 {
		limit = 100
	}
	if limit < 1 {
		limit = 20
	}

	// Get farms user has access to
	query := `
	SELECT f.id, f.owner_id, f.name, f.location, f.timezone, f.latitude, f.longitude, 
	       f.description, f.image_url, f.is_active, f.created_at, f.updated_at,
	       fu.role,
	       COUNT(c.id) AS coop_count
	FROM farms f
	INNER JOIN farm_users fu ON f.id = fu.farm_id
	LEFT JOIN coops c ON f.id = c.farm_id AND c.is_active = true
	WHERE fu.user_id = $1 AND fu.is_active = true AND f.is_active = true
	GROUP BY f.id, f.owner_id, f.name, f.location, f.timezone, f.latitude, f.longitude,
	         f.description, f.image_url, f.is_active, f.created_at, f.updated_at, fu.role
	ORDER BY f.created_at DESC
	LIMIT $2 OFFSET $3
	`

	rows, err := database.DB.Query(query, userID, limit, offset)
	if err != nil {
		log.Printf("List farms error: %v", err)
		return utils.InternalError(c, "Failed to fetch farms")
	}
	defer rows.Close()

	type FarmWithRole struct {
		models.Farm
		Role      string `json:"role"`
		CoopCount int    `json:"coop_count"`
	}

	farms := []FarmWithRole{}
	for rows.Next() {
		var farm FarmWithRole
		err := rows.Scan(
			&farm.ID, &farm.OwnerID, &farm.Name, &farm.Location, &farm.Timezone,
			&farm.Latitude, &farm.Longitude, &farm.Description, &farm.ImageURL,
			&farm.IsActive, &farm.CreatedAt, &farm.UpdatedAt,
			&farm.Role, &farm.CoopCount,
		)
		if err != nil {
			log.Printf("Scan farm error: %v", err)
			continue
		}
		farms = append(farms, farm)
	}

	// Get total count
	var total int64
	countQuery := `
	SELECT COUNT(DISTINCT f.id)
	FROM farms f
	INNER JOIN farm_users fu ON f.id = fu.farm_id
	WHERE fu.user_id = $1 AND fu.is_active = true AND f.is_active = true
	`
	database.DB.QueryRow(countQuery, userID).Scan(&total)

	return utils.SuccessListResponse(c, farms, total, offset/limit+1, limit)
}

// GetFarmHandler returns a single farm by ID
// @GET /v1/farms/:farm_id
func GetFarmHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	farmIDStr := c.Params("farm_id")
	farmID, err := uuid.Parse(farmIDStr)
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	// Check user access to farm
	var userRole string
	err = database.DB.QueryRow(`
		SELECT role FROM farm_users 
		WHERE farm_id = $1 AND user_id = $2 AND is_active = true
	`, farmID, userID).Scan(&userRole)

	if err == sql.ErrNoRows {
		return utils.Forbidden(c, "You do not have access to this farm")
	}
	if err != nil {
		log.Printf("Check farm access error: %v", err)
		return utils.InternalError(c, "Failed to verify access")
	}

	// Get farm details with coop count
	var farm struct {
		models.Farm
		Role      string `json:"role"`
		CoopCount int    `json:"coop_count"`
	}
	farm.Role = userRole

	query := `
	SELECT f.id, f.owner_id, f.name, f.location, f.timezone, f.latitude, f.longitude,
	       f.description, f.image_url, f.is_active, f.created_at, f.updated_at,
	       COUNT(c.id) AS coop_count
	FROM farms f
	LEFT JOIN coops c ON f.id = c.farm_id AND c.is_active = true
	WHERE f.id = $1
	GROUP BY f.id, f.owner_id, f.name, f.location, f.timezone, f.latitude, f.longitude,
	         f.description, f.image_url, f.is_active, f.created_at, f.updated_at
	`

	err = database.DB.QueryRow(query, farmID).Scan(
		&farm.ID, &farm.OwnerID, &farm.Name, &farm.Location, &farm.Timezone,
		&farm.Latitude, &farm.Longitude, &farm.Description, &farm.ImageURL,
		&farm.IsActive, &farm.CreatedAt, &farm.UpdatedAt, &farm.CoopCount,
	)

	if err == sql.ErrNoRows {
		return utils.NotFound(c, "Farm not found")
	}
	if err != nil {
		log.Printf("Get farm error: %v", err)
		return utils.InternalError(c, "Failed to fetch farm")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, farm, "Farm fetched successfully")
}

// CreateFarmHandler creates a new farm (user becomes owner)
// @POST /v1/farms
func CreateFarmHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	var req struct {
		Name        string   `json:"name"`
		Location    *string  `json:"location,omitempty"`
		Timezone    *string  `json:"timezone,omitempty"`
		Latitude    *float64 `json:"latitude,omitempty"`
		Longitude   *float64 `json:"longitude,omitempty"`
		Description *string  `json:"description,omitempty"`
		ImageURL    *string  `json:"image_url,omitempty"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid_request", "Invalid request body")
	}

	if req.Name == "" {
		return utils.BadRequest(c, "missing_name", "Farm name is required")
	}

	timezone := "Asia/Phnom_Penh"
	if req.Timezone != nil && *req.Timezone != "" {
		timezone = *req.Timezone
	}

	// Start transaction
	tx, err := database.DB.Begin()
	if err != nil {
		log.Printf("Begin transaction error: %v", err)
		return utils.InternalError(c, "Failed to create farm")
	}
	defer tx.Rollback()

	// Create farm
	farmID := uuid.New()
	query := `
	INSERT INTO farms (id, owner_id, name, location, timezone, latitude, longitude, description, image_url, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	RETURNING id, owner_id, name, location, timezone, latitude, longitude, description, image_url, is_active, created_at, updated_at
	`

	var farm models.Farm
	err = tx.QueryRow(
		query,
		farmID, userID, req.Name, req.Location, timezone,
		req.Latitude, req.Longitude, req.Description, req.ImageURL,
	).Scan(
		&farm.ID, &farm.OwnerID, &farm.Name, &farm.Location, &farm.Timezone,
		&farm.Latitude, &farm.Longitude, &farm.Description, &farm.ImageURL,
		&farm.IsActive, &farm.CreatedAt, &farm.UpdatedAt,
	)

	if err != nil {
		log.Printf("Create farm error: %v", err)
		return utils.InternalError(c, "Failed to create farm")
	}

	// Add creator as owner in farm_users
	farmUserID := uuid.New()
	_, err = tx.Exec(`
		INSERT INTO farm_users (id, farm_id, user_id, role, invited_by, created_at, updated_at)
		VALUES ($1, $2, $3, 'owner', $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, farmUserID, farmID, userID)

	if err != nil {
		log.Printf("Create farm_user error: %v", err)
		return utils.InternalError(c, "Failed to assign ownership")
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Printf("Commit transaction error: %v", err)
		return utils.InternalError(c, "Failed to create farm")
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, farm, "Farm created successfully")
}

// UpdateFarmHandler updates farm details (owner/manager only)
// @PUT /v1/farms/:farm_id
func UpdateFarmHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	farmIDStr := c.Params("farm_id")
	farmID, err := uuid.Parse(farmIDStr)
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	// Check user has owner or manager role
	var userRole string
	err = database.DB.QueryRow(`
		SELECT role FROM farm_users 
		WHERE farm_id = $1 AND user_id = $2 AND is_active = true
	`, farmID, userID).Scan(&userRole)

	if err == sql.ErrNoRows {
		return utils.Forbidden(c, "You do not have access to this farm")
	}
	if err != nil {
		log.Printf("Check farm access error: %v", err)
		return utils.InternalError(c, "Failed to verify access")
	}

	if userRole != "owner" && userRole != "manager" {
		return utils.Forbidden(c, "Only farm owners and managers can update farm details")
	}

	var req struct {
		Name        *string  `json:"name,omitempty"`
		Location    *string  `json:"location,omitempty"`
		Timezone    *string  `json:"timezone,omitempty"`
		Latitude    *float64 `json:"latitude,omitempty"`
		Longitude   *float64 `json:"longitude,omitempty"`
		Description *string  `json:"description,omitempty"`
		ImageURL    *string  `json:"image_url,omitempty"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid_request", "Invalid request body")
	}

	// Build dynamic UPDATE query (similar to UpdateProfileHandler)
	// For brevity, implementing a simpler approach here
	query := `
	UPDATE farms SET
		name = COALESCE($1, name),
		location = COALESCE($2, location),
		timezone = COALESCE($3, timezone),
		latitude = COALESCE($4, latitude),
		longitude = COALESCE($5, longitude),
		description = COALESCE($6, description),
		image_url = COALESCE($7, image_url),
		updated_at = CURRENT_TIMESTAMP
	WHERE id = $8
	RETURNING id, owner_id, name, location, timezone, latitude, longitude, description, image_url, is_active, created_at, updated_at
	`

	var farm models.Farm
	err = database.DB.QueryRow(
		query,
		req.Name, req.Location, req.Timezone, req.Latitude, req.Longitude,
		req.Description, req.ImageURL, farmID,
	).Scan(
		&farm.ID, &farm.OwnerID, &farm.Name, &farm.Location, &farm.Timezone,
		&farm.Latitude, &farm.Longitude, &farm.Description, &farm.ImageURL,
		&farm.IsActive, &farm.CreatedAt, &farm.UpdatedAt,
	)

	if err != nil {
		log.Printf("Update farm error: %v", err)
		return utils.InternalError(c, "Failed to update farm")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, farm, "Farm updated successfully")
}

// DeleteFarmHandler soft-deletes a farm (owner only)
// @DELETE /v1/farms/:farm_id
func DeleteFarmHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	farmIDStr := c.Params("farm_id")
	farmID, err := uuid.Parse(farmIDStr)
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	// Check user is owner
	var userRole string
	err = database.DB.QueryRow(`
		SELECT role FROM farm_users 
		WHERE farm_id = $1 AND user_id = $2 AND is_active = true
	`, farmID, userID).Scan(&userRole)

	if err == sql.ErrNoRows {
		return utils.Forbidden(c, "You do not have access to this farm")
	}
	if err != nil {
		log.Printf("Check farm access error: %v", err)
		return utils.InternalError(c, "Failed to verify access")
	}

	if userRole != "owner" {
		return utils.Forbidden(c, "Only farm owners can delete farms")
	}

	// Soft delete (set is_active = false)
	_, err = database.DB.Exec("UPDATE farms SET is_active = false, updated_at = CURRENT_TIMESTAMP WHERE id = $1", farmID)
	if err != nil {
		log.Printf("Delete farm error: %v", err)
		return utils.InternalError(c, "Failed to delete farm")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"message": "Farm deleted successfully",
	}, "Farm deleted")
}
