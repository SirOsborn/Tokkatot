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

// ===== COOP MANAGEMENT HANDLERS =====

// ListCoopsHandler returns all coops in a farm
// @GET /v1/farms/:farm_id/coops?limit=20&offset=0
func ListCoopsHandler(c *fiber.Ctx) error {
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
	err = checkFarmAccess(userID, farmID, "viewer") // Minimum role: viewer
	if err != nil {
		return err
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

	// Get coops with device count
	query := `
	SELECT c.id, c.farm_id, c.number, c.name, c.capacity, c.current_count, 
	       c.chicken_type, c.main_device_id, c.description, c.is_active, 
	       c.created_at, c.updated_at,
	       COUNT(d.id) AS device_count
	FROM coops c
	LEFT JOIN devices d ON c.id = d.coop_id AND d.is_active = true
	WHERE c.farm_id = $1 AND c.is_active = true
	GROUP BY c.id, c.farm_id, c.number, c.name, c.capacity, c.current_count,
	         c.chicken_type, c.main_device_id, c.description, c.is_active,
	         c.created_at, c.updated_at
	ORDER BY c.number ASC
	LIMIT $2 OFFSET $3
	`

	rows, err := database.DB.Query(query, farmID, limit, offset)
	if err != nil {
		log.Printf("List coops error: %v", err)
		return utils.InternalError(c, "Failed to fetch coops")
	}
	defer rows.Close()

	type CoopWithDevices struct {
		models.Coop
		DeviceCount int `json:"device_count"`
	}

	coops := []CoopWithDevices{}
	for rows.Next() {
		var coop CoopWithDevices
		err := rows.Scan(
			&coop.ID, &coop.FarmID, &coop.Number, &coop.Name, &coop.Capacity,
			&coop.CurrentCount, &coop.ChickenType, &coop.MainDeviceID, &coop.Description,
			&coop.IsActive, &coop.CreatedAt, &coop.UpdatedAt, &coop.DeviceCount,
		)
		if err != nil {
			log.Printf("Scan coop error: %v", err)
			continue
		}
		coops = append(coops, coop)
	}

	// Get total count
	var total int64
	database.DB.QueryRow("SELECT COUNT(*) FROM coops WHERE farm_id = $1 AND is_active = true", farmID).Scan(&total)

	return utils.SuccessListResponse(c, coops, total, offset/limit+1, limit)
}

// GetCoopHandler returns a single coop by ID
// @GET /v1/farms/:farm_id/coops/:coop_id
func GetCoopHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	farmIDStr := c.Params("farm_id")
	farmID, err := uuid.Parse(farmIDStr)
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	coopIDStr := c.Params("coop_id")
	coopID, err := uuid.Parse(coopIDStr)
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid coop ID")
	}

	// Check user access
	err = checkFarmAccess(userID, farmID, "viewer")
	if err != nil {
		return err
	}

	// Get coop with device count
	var coop struct {
		models.Coop
		DeviceCount int `json:"device_count"`
	}

	query := `
	SELECT c.id, c.farm_id, c.number, c.name, c.capacity, c.current_count,
	       c.chicken_type, c.main_device_id, c.description, c.is_active,
	       c.created_at, c.updated_at,
	       COUNT(d.id) AS device_count
	FROM coops c
	LEFT JOIN devices d ON c.id = d.coop_id AND d.is_active = true
	WHERE c.id = $1 AND c.farm_id = $2
	GROUP BY c.id, c.farm_id, c.number, c.name, c.capacity, c.current_count,
	         c.chicken_type, c.main_device_id, c.description, c.is_active,
	         c.created_at, c.updated_at
	`

	err = database.DB.QueryRow(query, coopID, farmID).Scan(
		&coop.ID, &coop.FarmID, &coop.Number, &coop.Name, &coop.Capacity,
		&coop.CurrentCount, &coop.ChickenType, &coop.MainDeviceID, &coop.Description,
		&coop.IsActive, &coop.CreatedAt, &coop.UpdatedAt, &coop.DeviceCount,
	)

	if err == sql.ErrNoRows {
		return utils.NotFound(c, "Coop not found")
	}
	if err != nil {
		log.Printf("Get coop error: %v", err)
		return utils.InternalError(c, "Failed to fetch coop")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, coop, "Coop fetched successfully")
}

// CreateCoopHandler creates a new coop (manager/owner only)
// @POST /v1/farms/:farm_id/coops
func CreateCoopHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	farmIDStr := c.Params("farm_id")
	farmID, err := uuid.Parse(farmIDStr)
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	// Check user has manager or owner role
	err = checkFarmAccess(userID, farmID, "manager")
	if err != nil {
		return err
	}

	var req struct {
		Number       int     `json:"number"`
		Name         string  `json:"name"`
		Capacity     *int    `json:"capacity,omitempty"`
		CurrentCount *int    `json:"current_count,omitempty"`
		ChickenType  *string `json:"chicken_type,omitempty"`
		Description  *string `json:"description,omitempty"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid_request", "Invalid request body")
	}

	if req.Number < 1 {
		return utils.BadRequest(c, "invalid_number", "Coop number must be at least 1")
	}

	if req.Name == "" {
		return utils.BadRequest(c, "missing_name", "Coop name is required")
	}

	// Check if coop number already exists for this farm
	var exists bool
	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM coops WHERE farm_id = $1 AND number = $2)", farmID, req.Number).Scan(&exists)
	if err != nil {
		log.Printf("Coop number check error: %v", err)
		return utils.InternalError(c, "Failed to validate coop number")
	}
	if exists {
		return utils.Conflict(c, "coop_exists", "A coop with this number already exists in this farm")
	}

	// Create coop
	coopID := uuid.New()
	query := `
	INSERT INTO coops (id, farm_id, number, name, capacity, current_count, chicken_type, description, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	RETURNING id, farm_id, number, name, capacity, current_count, chicken_type, main_device_id, description, is_active, created_at, updated_at
	`

	var coop models.Coop
	err = database.DB.QueryRow(
		query,
		coopID, farmID, req.Number, req.Name, req.Capacity, req.CurrentCount, req.ChickenType, req.Description,
	).Scan(
		&coop.ID, &coop.FarmID, &coop.Number, &coop.Name, &coop.Capacity,
		&coop.CurrentCount, &coop.ChickenType, &coop.MainDeviceID, &coop.Description,
		&coop.IsActive, &coop.CreatedAt, &coop.UpdatedAt,
	)

	if err != nil {
		log.Printf("Create coop error: %v", err)
		return utils.InternalError(c, "Failed to create coop")
	}

	return utils.SuccessResponse(c, fiber.StatusCreated, coop, "Coop created successfully")
}

// UpdateCoopHandler updates coop details (manager/owner only)
// @PUT /v1/farms/:farm_id/coops/:coop_id
func UpdateCoopHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	farmIDStr := c.Params("farm_id")
	farmID, err := uuid.Parse(farmIDStr)
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	coopIDStr := c.Params("coop_id")
	coopID, err := uuid.Parse(coopIDStr)
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid coop ID")
	}

	// Check user has manager or owner role
	err = checkFarmAccess(userID, farmID, "manager")
	if err != nil {
		return err
	}

	var req struct {
		Name         *string `json:"name,omitempty"`
		Capacity     *int    `json:"capacity,omitempty"`
		CurrentCount *int    `json:"current_count,omitempty"`
		ChickenType  *string `json:"chicken_type,omitempty"`
		Description  *string `json:"description,omitempty"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "invalid_request", "Invalid request body")
	}

	query := `
	UPDATE coops SET
		name = COALESCE($1, name),
		capacity = COALESCE($2, capacity),
		current_count = COALESCE($3, current_count),
		chicken_type = COALESCE($4, chicken_type),
		description = COALESCE($5, description),
		updated_at = CURRENT_TIMESTAMP
	WHERE id = $6 AND farm_id = $7
	RETURNING id, farm_id, number, name, capacity, current_count, chicken_type, main_device_id, description, is_active, created_at, updated_at
	`

	var coop models.Coop
	err = database.DB.QueryRow(
		query,
		req.Name, req.Capacity, req.CurrentCount, req.ChickenType, req.Description, coopID, farmID,
	).Scan(
		&coop.ID, &coop.FarmID, &coop.Number, &coop.Name, &coop.Capacity,
		&coop.CurrentCount, &coop.ChickenType, &coop.MainDeviceID, &coop.Description,
		&coop.IsActive, &coop.CreatedAt, &coop.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return utils.NotFound(c, "Coop not found")
	}
	if err != nil {
		log.Printf("Update coop error: %v", err)
		return utils.InternalError(c, "Failed to update coop")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, coop, "Coop updated successfully")
}

// DeleteCoopHandler soft-deletes a coop (owner only)
// @DELETE /v1/farms/:farm_id/coops/:coop_id
func DeleteCoopHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.Unauthorized(c, "Invalid token")
	}

	farmIDStr := c.Params("farm_id")
	farmID, err := uuid.Parse(farmIDStr)
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid farm ID")
	}

	coopIDStr := c.Params("coop_id")
	coopID, err := uuid.Parse(coopIDStr)
	if err != nil {
		return utils.BadRequest(c, "invalid_id", "Invalid coop ID")
	}

	// Check user is owner
	err = checkFarmAccess(userID, farmID, "owner")
	if err != nil {
		return err
	}

	// Soft delete coop
	_, err = database.DB.Exec("UPDATE coops SET is_active = false, updated_at = CURRENT_TIMESTAMP WHERE id = $1 AND farm_id = $2", coopID, farmID)
	if err != nil {
		log.Printf("Delete coop error: %v", err)
		return utils.InternalError(c, "Failed to delete coop")
	}

	return utils.SuccessResponse(c, fiber.StatusOK, fiber.Map{
		"message": "Coop deleted successfully",
	}, "Coop deleted")
}

// ===== HELPER FUNCTIONS =====

// checkFarmAccess verifies user has at least minimum role for farm
func checkFarmAccess(userID, farmID uuid.UUID, minRole string) error {
	var userRole string
	err := database.DB.QueryRow(`
		SELECT role FROM farm_users 
		WHERE farm_id = $1 AND user_id = $2 AND is_active = true
	`, farmID, userID).Scan(&userRole)

	if err == sql.ErrNoRows {
		return fiber.NewError(fiber.StatusForbidden, "You do not have access to this farm")
	}
	if err != nil {
		log.Printf("Check farm access error: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to verify access")
	}

	// Check role hierarchy: owner > manager > viewer
	roleHierarchy := map[string]int{
		"owner":   3,
		"manager": 2,
		"viewer":  1,
	}

	if roleHierarchy[userRole] < roleHierarchy[minRole] {
		return fiber.NewError(fiber.StatusForbidden, "Insufficient permissions")
	}

	return nil
}
