package services

import (
	"database/sql"
	"errors"
	"middleware/database"
	"middleware/models"
	"middleware/schemas"
	"time"

	"github.com/google/uuid"
)

var (
	ErrCoopNotFound = errors.New("coop not found")
)

// CoopService handles all business logic related to coop management
type CoopService struct {
	farmService *FarmService
}

func NewCoopService() *CoopService {
	return &CoopService{
		farmService: NewFarmService(),
	}
}

// ListCoops returns all coops for a farm
func (s *CoopService) ListCoops(userID, farmID uuid.UUID) ([]schemas.CoopWithDevices, error) {
	if err := s.farmService.CheckAccess(userID, farmID, "viewer"); err != nil {
		return nil, err
	}

	rows, err := database.DB.Query(`
		SELECT c.id, c.farm_id, c.number, c.name, c.capacity, c.current_count, c.chicken_type, 
		       c.main_device_id, c.description, c.is_active, c.created_at,
		       (SELECT COUNT(*) FROM devices WHERE coop_id = c.id AND is_active = true) as device_count
		FROM coops c
		WHERE c.farm_id = $1 AND c.is_active = true
		ORDER BY c.number ASC
	`, farmID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var coops []schemas.CoopWithDevices
	for rows.Next() {
		var c schemas.CoopWithDevices
		if err := rows.Scan(&c.ID, &c.FarmID, &c.Number, &c.Name, &c.Capacity, &c.CurrentCount, &c.ChickenType, 
			&c.MainDeviceID, &c.Description, &c.IsActive, &c.CreatedAt, &c.DeviceCount); err != nil {
			continue
		}
		coops = append(coops, c)
	}
	return coops, nil
}

// GetCoop retrieves details for a single coop
func (s *CoopService) GetCoop(userID, farmID, coopID uuid.UUID) (schemas.CoopWithDevices, error) {
	var c schemas.CoopWithDevices
	if err := s.farmService.CheckAccess(userID, farmID, "viewer"); err != nil {
		return c, err
	}

	err := database.DB.QueryRow(`
		SELECT c.id, c.farm_id, c.number, c.name, c.capacity, c.current_count, c.chicken_type, 
		       c.main_device_id, c.description, c.is_active, c.created_at,
		       (SELECT COUNT(*) FROM devices WHERE coop_id = c.id AND is_active = true) as device_count
		FROM coops c
		WHERE c.id = $1 AND c.farm_id = $2 AND c.is_active = true
	`, coopID, farmID).Scan(&c.ID, &c.FarmID, &c.Number, &c.Name, &c.Capacity, &c.CurrentCount, &c.ChickenType, 
		&c.MainDeviceID, &c.Description, &c.IsActive, &c.CreatedAt, &c.DeviceCount)

	if err == sql.ErrNoRows {
		return c, ErrCoopNotFound
	}
	return c, err
}

// CreateCoop creates a new coop
func (s *CoopService) CreateCoop(userID, farmID uuid.UUID, req models.Coop) (*models.Coop, error) {
	if err := s.farmService.CheckAccess(userID, farmID, "farmer"); err != nil {
		return nil, err
	}

	coopID := uuid.New()
	now := time.Now()
	_, err := database.DB.Exec(`
		INSERT INTO coops (id, farm_id, number, name, capacity, current_count, chicken_type, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`, coopID, farmID, req.Number, req.Name, req.Capacity, req.CurrentCount, req.ChickenType, req.Description, now, now)

	if err != nil {
		return nil, err
	}

	req.ID = coopID
	req.FarmID = farmID
	req.CreatedAt = now
	req.UpdatedAt = now
	return &req, nil
}

// UpdateCoop updates an existing coop
func (s *CoopService) UpdateCoop(userID, farmID, coopID uuid.UUID, number *int, name *string, capacity *int, currentCount *int, chickenType *string, description *string) (*models.Coop, error) {
	if err := s.farmService.CheckAccess(userID, farmID, "farmer"); err != nil {
		return nil, err
	}

	var c models.Coop
	err := database.DB.QueryRow(`
		UPDATE coops SET
			number = COALESCE($1, number),
			name = COALESCE($2, name),
			capacity = COALESCE($3, capacity),
			current_count = COALESCE($4, current_count),
			chicken_type = COALESCE($5, chicken_type),
			description = COALESCE($6, description),
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $7 AND farm_id = $8 AND is_active = true
		RETURNING id, farm_id, number, name, capacity, current_count, chicken_type, main_device_id, description, is_active, created_at, updated_at
	`, number, name, capacity, currentCount, chickenType, description, coopID, farmID).Scan(
		&c.ID, &c.FarmID, &c.Number, &c.Name, &c.Capacity, &c.CurrentCount, &c.ChickenType, &c.MainDeviceID, &c.Description, &c.IsActive, &c.CreatedAt, &c.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrCoopNotFound
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// DeleteCoop deactivates a coop
func (s *CoopService) DeleteCoop(userID, farmID, coopID uuid.UUID) error {
	if err := s.farmService.CheckAccess(userID, farmID, "farmer"); err != nil {
		return err
	}
	res, err := database.DB.Exec(`
		UPDATE coops SET is_active = false, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND farm_id = $2
	`, coopID, farmID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return ErrCoopNotFound
	}
	return nil
}
