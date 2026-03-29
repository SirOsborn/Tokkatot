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

func (s *CoopService) attachLatestTelemetry(coopID uuid.UUID, c *schemas.CoopWithDevices) {
	if c == nil || coopID == uuid.Nil {
		return
	}

	var (
		tempVal sql.NullFloat64
		tempTs  sql.NullTime
		humVal  sql.NullFloat64
		humTs   sql.NullTime
	)

	_ = database.DB.QueryRow(`
		SELECT dr.value, dr.timestamp
		FROM device_readings dr
		JOIN devices d ON dr.device_id = d.id
		WHERE d.coop_id = $1 AND dr.sensor_type = 'temperature'
		ORDER BY dr.timestamp DESC
		LIMIT 1
	`, coopID).Scan(&tempVal, &tempTs)

	_ = database.DB.QueryRow(`
		SELECT dr.value, dr.timestamp
		FROM device_readings dr
		JOIN devices d ON dr.device_id = d.id
		WHERE d.coop_id = $1 AND dr.sensor_type = 'humidity'
		ORDER BY dr.timestamp DESC
		LIMIT 1
	`, coopID).Scan(&humVal, &humTs)

	if tempVal.Valid {
		v := tempVal.Float64
		c.Temperature = &v
	}
	if humVal.Valid {
		v := humVal.Float64
		c.Humidity = &v
	}

	var latest time.Time
	if tempTs.Valid && tempTs.Time.After(latest) {
		latest = tempTs.Time
	}
	if humTs.Valid && humTs.Time.After(latest) {
		latest = humTs.Time
	}
	if !latest.IsZero() {
		c.LastUpdated = &latest
	}
}

// ListCoops returns all coops for a farm
func (s *CoopService) ListCoops(userID, farmID uuid.UUID) ([]schemas.CoopWithDevices, error) {
	if err := s.farmService.CheckAccess(userID, farmID, "viewer"); err != nil {
		return nil, err
	}

	rows, err := database.DB.Query(`
		SELECT c.id, c.farm_id, c.number, c.name, c.capacity, c.current_count, c.chicken_type, 
		       c.main_device_id, c.temp_min, c.temp_max, c.water_level_half_threshold, c.description, c.is_active, c.created_at,
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
			&c.MainDeviceID, &c.TempMin, &c.TempMax, &c.WaterLevelHalfThreshold, &c.Description, &c.IsActive, &c.CreatedAt, &c.DeviceCount); err != nil {
			continue
		}
		s.attachLatestTelemetry(c.ID, &c)
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
		       c.main_device_id, c.temp_min, c.temp_max, c.water_level_half_threshold, c.description, c.is_active, c.created_at,
		       (SELECT COUNT(*) FROM devices WHERE coop_id = c.id AND is_active = true) as device_count
		FROM coops c
		WHERE c.id = $1 AND c.farm_id = $2 AND c.is_active = true
	`, coopID, farmID).Scan(&c.ID, &c.FarmID, &c.Number, &c.Name, &c.Capacity, &c.CurrentCount, &c.ChickenType, 
		&c.MainDeviceID, &c.TempMin, &c.TempMax, &c.WaterLevelHalfThreshold, &c.Description, &c.IsActive, &c.CreatedAt, &c.DeviceCount)

	if err == sql.ErrNoRows {
		return c, ErrCoopNotFound
	}
	if err == nil {
		s.attachLatestTelemetry(c.ID, &c)
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
		INSERT INTO coops (id, farm_id, number, name, capacity, current_count, chicken_type, temp_min, temp_max, water_level_half_threshold, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`, coopID, farmID, req.Number, req.Name, req.Capacity, req.CurrentCount, req.ChickenType, req.TempMin, req.TempMax, req.WaterLevelHalfThreshold, req.Description, now, now)

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
func (s *CoopService) UpdateCoop(userID, farmID, coopID uuid.UUID, number *int, name *string, capacity *int, currentCount *int, chickenType *string, tempMin *float64, tempMax *float64, waterHalf *float64, description *string) (*models.Coop, error) {
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
			temp_min = COALESCE($6, temp_min),
			temp_max = COALESCE($7, temp_max),
			water_level_half_threshold = COALESCE($8, water_level_half_threshold),
			description = COALESCE($9, description),
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $10 AND farm_id = $11 AND is_active = true
		RETURNING id, farm_id, number, name, capacity, current_count, chicken_type, main_device_id, temp_min, temp_max, water_level_half_threshold, description, is_active, created_at, updated_at
	`, number, name, capacity, currentCount, chickenType, tempMin, tempMax, waterHalf, description, coopID, farmID).Scan(
		&c.ID, &c.FarmID, &c.Number, &c.Name, &c.Capacity, &c.CurrentCount, &c.ChickenType, &c.MainDeviceID, &c.TempMin, &c.TempMax, &c.WaterLevelHalfThreshold, &c.Description, &c.IsActive, &c.CreatedAt, &c.UpdatedAt,
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
