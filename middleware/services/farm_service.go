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
	ErrFarmNotFound     = errors.New("farm not found")
	ErrFarmAccessDenied = errors.New("access denied to farm")
)

// FarmService handles all business logic related to farm management
type FarmService struct{}

func NewFarmService() *FarmService {
	return &FarmService{}
}

// CheckAccess verifies if a user has a specific role (or better) in a farm
func (s *FarmService) CheckAccess(userID, farmID uuid.UUID, minRole string) error {
	var role string
	err := database.DB.QueryRow(`
		SELECT role FROM farm_users 
		WHERE farm_id = $1 AND user_id = $2
	`, farmID, userID).Scan(&role)

	if err == sql.ErrNoRows {
		return ErrFarmAccessDenied
	}

	// Simple role hierarchy: farmer > worker > viewer
	roleWeight := map[string]int{"farmer": 3, "worker": 2, "viewer": 1}
	if roleWeight[role] < roleWeight[minRole] {
		return ErrFarmAccessDenied
	}

	return nil
}

// ListFarms returns all farms the user is a member of with pagination
func (s *FarmService) ListFarms(userID uuid.UUID, limit, offset int) ([]schemas.FarmWithRole, int64, error) {
	query := `
		SELECT f.id, f.name, f.location, f.description, fu.role, f.created_at,
		       (SELECT COUNT(*) FROM coops WHERE farm_id = f.id AND is_active = true) as coop_count
		FROM farms f
		JOIN farm_users fu ON f.id = fu.farm_id
		WHERE fu.user_id = $1 AND f.is_active = true
		ORDER BY f.created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := database.DB.Query(query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var farms []schemas.FarmWithRole
	for rows.Next() {
		var f schemas.FarmWithRole
		if err := rows.Scan(&f.ID, &f.Name, &f.Location, &f.Description, &f.Role, &f.CreatedAt, &f.CoopCount); err != nil {
			continue
		}
		farms = append(farms, f)
	}

	var total int64
	database.DB.QueryRow("SELECT COUNT(*) FROM farm_users WHERE user_id = $1 AND is_active = true", userID).Scan(&total)

	return farms, total, nil
}

// GetFarm retrieves details for a single farm
func (s *FarmService) GetFarm(userID, farmID uuid.UUID) (schemas.FarmWithRole, error) {
	var f schemas.FarmWithRole
	err := database.DB.QueryRow(`
		SELECT f.id, f.name, f.location, f.description, fu.role, f.created_at,
		       (SELECT COUNT(*) FROM coops WHERE farm_id = f.id AND is_active = true) as coop_count
		FROM farms f
		JOIN farm_users fu ON f.id = fu.farm_id
		WHERE f.id = $1 AND fu.user_id = $2 AND f.is_active = true
	`, farmID, userID).Scan(&f.ID, &f.Name, &f.Location, &f.Description, &f.Role, &f.CreatedAt, &f.CoopCount)

	if err == sql.ErrNoRows {
		return f, ErrFarmNotFound
	}
	return f, err
}

// CreateFarm creates a new farm and makes the creator the 'farmer'
func (s *FarmService) CreateFarm(userID uuid.UUID, req schemas.CreateFarmRequest) (*models.Farm, error) {
	tx, err := database.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	farmID := uuid.New()
	now := time.Now()
	_, err = tx.Exec(`
		INSERT INTO farms (id, name, location, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, farmID, req.Name, req.Location, req.Description, now, now)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(`
		INSERT INTO farm_users (farm_id, user_id, role, created_at)
		VALUES ($1, $2, 'farmer', $3)
	`, farmID, userID, now)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.Farm{
		ID:          farmID,
		Name:        req.Name,
		Location:    req.Location,
		Description: req.Description,
		CreatedAt:   now,
	}, nil
}

// UpdateFarm updates farm details
func (s *FarmService) UpdateFarm(userID, farmID uuid.UUID, req schemas.UpdateFarmRequest) (*models.Farm, error) {
	if err := s.CheckAccess(userID, farmID, "farmer"); err != nil {
		return nil, err
	}

	query := `
		UPDATE farms SET
			name = COALESCE($1, name),
			location = COALESCE($2, location),
			description = COALESCE($3, description),
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
		RETURNING id, name, location, description, created_at
	`
	var f models.Farm
	err := database.DB.QueryRow(query, req.Name, req.Location, req.Description, farmID).
		Scan(&f.ID, &f.Name, &f.Location, &f.Description, &f.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, ErrFarmNotFound
	}
	return &f, err
}

// DeleteFarm soft-deletes a farm
func (s *FarmService) DeleteFarm(userID, farmID uuid.UUID) error {
	if err := s.CheckAccess(userID, farmID, "farmer"); err != nil {
		return err
	}

	_, err := database.DB.Exec("UPDATE farms SET is_active = false, updated_at = CURRENT_TIMESTAMP WHERE id = $1", farmID)
	return err
}

// GetFarmMembers returns membership list for a farm
func (s *FarmService) GetFarmMembers(userID, farmID uuid.UUID) ([]schemas.MemberInfo, int64, error) {
	if err := s.CheckAccess(userID, farmID, "viewer"); err != nil {
		return nil, 0, err
	}

	rows, err := database.DB.Query(`
		SELECT fu.id, fu.user_id, u.name, u.email, u.phone, fu.role, fu.created_at
		FROM farm_users fu
		JOIN users u ON fu.user_id = u.id
		WHERE fu.farm_id = $1 AND fu.is_active = true
		ORDER BY fu.created_at ASC
	`, farmID)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var members []schemas.MemberInfo
	for rows.Next() {
		var m schemas.MemberInfo
		if err := rows.Scan(&m.ID, &m.UserID, &m.Name, &m.Email, &m.Phone, &m.Role, &m.JoinedAt); err != nil {
			continue
		}
		members = append(members, m)
	}

	var total int64
	database.DB.QueryRow("SELECT COUNT(*) FROM farm_users WHERE farm_id = $1 AND is_active = true", farmID).Scan(&total)

	return members, total, nil
}
