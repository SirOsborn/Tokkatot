package services

import (
	"crypto/rand"
	"database/sql"
	"math/big"
	"middleware/database"
	"middleware/models"
	"time"

	"github.com/google/uuid"
)

type AdminService struct{}

func NewAdminService() *AdminService {
	return &AdminService{}
}

// GenerateRegistrationKey creates a new one-time use key for farmer signup
func (s *AdminService) GenerateRegistrationKey(farmName, customerPhone, nationalID, fullName, sex, province string, expiryDays int) (*models.RegistrationKey, error) {
	keyCode, err := s.generateSecureKey(25)
	if err != nil {
		return nil, err
	}

	// Format key with hyphens: XXXXX-XXXXX-XXXXX-XXXXX-XXXXX
	formattedKey := ""
	for i := 0; i < 25; i++ {
		if i > 0 && i%5 == 0 {
			formattedKey += "-"
		}
		formattedKey += string(keyCode[i])
	}

	id := uuid.New()
	now := time.Now()
	expiresAt := now.AddDate(0, 0, expiryDays)

	regKey := &models.RegistrationKey{
		ID:               id,
		KeyCode:          formattedKey,
		FarmName:         &farmName,
		CustomerPhone:    &customerPhone,
		IsUsed:           false,
		ExpiresAt:        &expiresAt,
		CreatedBy:        "admin_api",
		NationalIDNumber: &nationalID,
		FullName:         &fullName,
		Sex:              &sex,
		Province:         &province,
		CreatedAt:        now,
	}

	_, err = database.DB.Exec(`
		INSERT INTO registration_keys 
		(id, key_code, farm_name, customer_phone, expires_at, created_by, national_id_number, full_name, sex, province, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`, regKey.ID, regKey.KeyCode, regKey.FarmName, regKey.CustomerPhone, regKey.ExpiresAt, regKey.CreatedBy, regKey.NationalIDNumber, regKey.FullName, regKey.Sex, regKey.Province, regKey.CreatedAt)

	if err != nil {
		return nil, err
	}

	return regKey, nil
}

// ListRegistrationKeys returns all keys in the system
func (s *AdminService) ListRegistrationKeys() ([]models.RegistrationKey, error) {
	rows, err := database.DB.Query(`
		SELECT rk.id, rk.key_code, rk.farm_name, 
		       COALESCE(u.phone, rk.customer_phone) as phone, 
		       rk.is_used, rk.expires_at, rk.created_at,
		       COALESCE(u.national_id_number, rk.national_id_number) as national_id,
		       COALESCE(u.name, rk.full_name) as full_name,
		       COALESCE(u.sex, rk.sex) as sex,
		       COALESCE(u.province, rk.province) as province
		FROM registration_keys rk
		LEFT JOIN users u ON rk.used_by_user_id = u.id
		ORDER BY rk.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	keys := make([]models.RegistrationKey, 0)
	for rows.Next() {
		var k models.RegistrationKey
		var farmName, customerPhone, nationalID, fullName, sex, province sql.NullString
		err := rows.Scan(&k.ID, &k.KeyCode, &farmName, &customerPhone, &k.IsUsed, &k.ExpiresAt, &k.CreatedAt, &nationalID, &fullName, &sex, &province)
		if err != nil {
			return nil, err
		}
		if farmName.Valid { k.FarmName = &farmName.String }
		if customerPhone.Valid { k.CustomerPhone = &customerPhone.String }
		if nationalID.Valid { k.NationalIDNumber = &nationalID.String }
		if fullName.Valid { k.FullName = &fullName.String }
		if sex.Valid { k.Sex = &sex.String }
		if province.Valid { k.Province = &province.String }
		keys = append(keys, k)
	}
	return keys, nil
}

// ListAllFarmers returns all users with farmer role
func (s *AdminService) ListAllFarmers() ([]models.User, error) {
	rows, err := database.DB.Query(`
		SELECT
			u.id,
			u.email,
			u.phone,
			u.name,
			u.is_active,
			u.created_at,
			u.updated_at,
			u.last_login,
			u.national_id_number,
			u.sex,
			u.province,
			u.full_name,
			picked.farm_id,
			picked.farm_name,
			picked.role
		FROM users u
		JOIN LATERAL (
			SELECT
				fu.farm_id,
				f.name AS farm_name,
				fu.role
			FROM farm_users fu
			JOIN farms f ON f.id = fu.farm_id
			WHERE fu.user_id = u.id
			  AND fu.role = 'farmer'
			  AND fu.is_active = true
			  AND f.is_active = true
			ORDER BY
				CASE WHEN f.owner_id = u.id THEN 0 ELSE 1 END,
				f.created_at DESC
			LIMIT 1
		) AS picked ON true
		ORDER BY u.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]models.User, 0)
	for rows.Next() {
		var u models.User
		var lastLogin sql.NullTime
		var nationalID, sex, province, fullName sql.NullString
		var farmName sql.NullString
		var farmID sql.NullString
		err := rows.Scan(
			&u.ID,
			&u.Email,
			&u.Phone,
			&u.Name,
			&u.IsActive,
			&u.CreatedAt,
			&u.UpdatedAt,
			&lastLogin,
			&nationalID,
			&sex,
			&province,
			&fullName,
			&farmID,
			&farmName,
			&u.Role,
		)
		if err != nil {
			return nil, err
		}
		if lastLogin.Valid {
			u.LastLogin = &lastLogin.Time
		}
		if nationalID.Valid { u.NationalIDNumber = &nationalID.String }
		if sex.Valid { u.Sex = &sex.String }
		if province.Valid { u.Province = &province.String }
		if fullName.Valid { u.FullName = &fullName.String }
		if farmID.Valid {
			parsed, perr := uuid.Parse(farmID.String)
			if perr == nil {
				u.FarmID = &parsed
			}
		}
		if farmName.Valid { u.FarmName = &farmName.String }
		
		users = append(users, u)
	}
	return users, nil
}

// DeactivateUser toggles the active status of a user
func (s *AdminService) DeactivateUser(userID uuid.UUID, active bool) error {
	_, err := database.DB.Exec("UPDATE users SET is_active = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2", active, userID)
	return err
}

// IsUserActive checks if a user (farmer or admin) is currently active
func (s *AdminService) IsUserActive(userID uuid.UUID) (bool, error) {
	var isActive bool
	err := database.DB.QueryRow("SELECT is_active FROM users WHERE id = $1", userID).Scan(&isActive)
	if err == sql.ErrNoRows {
		// Check admin table
		err = database.DB.QueryRow("SELECT is_active FROM admins WHERE id = $1", userID).Scan(&isActive)
	}
	
	if err != nil {
		return false, err
	}
	return isActive, nil
}

// GetAdminStats returns counts of farmers, workers, farms, and unused keys
func (s *AdminService) GetAdminStats() (map[string]int, error) {
	var totalFarmers, totalWorkers, totalFarms, activeKeys int

	// Count Farmers
	err := database.DB.QueryRow(`
		SELECT COUNT(DISTINCT u.id) 
		FROM users u 
		JOIN farm_users fu ON u.id = fu.user_id 
		WHERE fu.role = 'farmer'
	`).Scan(&totalFarmers)
	if err != nil {
		return nil, err
	}

	// Count Workers (viewer + worker)
	err = database.DB.QueryRow(`
		SELECT COUNT(DISTINCT u.id) 
		FROM users u 
		JOIN farm_users fu ON u.id = fu.user_id 
		WHERE fu.role IN ('viewer', 'worker')
	`).Scan(&totalWorkers)
	if err != nil {
		return nil, err
	}

	// Count Farms
	err = database.DB.QueryRow("SELECT COUNT(*) FROM farms").Scan(&totalFarms)
	if err != nil {
		return nil, err
	}

	// Count Unused Keys
	err = database.DB.QueryRow("SELECT COUNT(*) FROM registration_keys WHERE is_used = false").Scan(&activeKeys)
	if err != nil {
		return nil, err
	}

	// Count Active Gateways
	var activeGateways int
	_ = database.DB.QueryRow("SELECT COUNT(*) FROM gateway_tokens WHERE is_active = true").Scan(&activeGateways)

	return map[string]int{
		"total_farmers":   totalFarmers,
		"total_workers":   totalWorkers,
		"total_farms":     totalFarms,
		"active_keys":     activeKeys,
		"active_gateways": activeGateways,
	}, nil
}

// ListGateways returns all persistent gateway tokens
func (s *AdminService) ListGateways() ([]map[string]interface{}, error) {
	rows, err := database.DB.Query(`
		SELECT id, farm_id, user_id, token_hash, name, is_active, last_used_at, created_at
		FROM gateway_tokens
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	gateways := make([]map[string]interface{}, 0)
	for rows.Next() {
		var id, farmID, userID uuid.UUID
		var tokenHash, name string
		var isActive bool
		var createdAt time.Time
		var lastUsedAtNull sql.NullTime

		err := rows.Scan(&id, &farmID, &userID, &tokenHash, &name, &isActive, &lastUsedAtNull, &createdAt)
		if err != nil {
			return nil, err
		}

		gw := map[string]interface{}{
			"id":           id,
			"farm_id":      farmID,
			"user_id":      userID,
			"token_hash":   tokenHash,
			"name":         name,
			"is_active":    isActive,
			"created_at":   createdAt,
		}
		if lastUsedAtNull.Valid {
			gw["last_used_at"] = lastUsedAtNull.Time
		}
		gateways = append(gateways, gw)
	}
	return gateways, nil
}

// RevokeGateway deactivates a gateway token
func (s *AdminService) RevokeGateway(id uuid.UUID) error {
	_, err := database.DB.Exec("UPDATE gateway_tokens SET is_active = false WHERE id = $1", id)
	return err
}

// generateSecureKey creates a random string of specified length
func (s *AdminService) generateSecureKey(length int) (string, error) {
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}
	return string(result), nil
}
