package services

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"math/big"
	"middleware/database"
	"middleware/utils"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ProvisioningService struct {
	farmService *FarmService
}

func NewProvisioningService() *ProvisioningService {
	return &ProvisioningService{
		farmService: NewFarmService(),
	}
}

// GenerateSetupCode creates a human-readable 6-digit code
func (s *ProvisioningService) GenerateSetupCode() (string, error) {
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789" // No confusing chars like 0/O or 1/I
	code := make([]byte, 6)
	for i := range code {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		code[i] = charset[num.Int64()]
	}
	// Format as XXX-XXX
	return fmt.Sprintf("%s-%s", string(code[:3]), string(code[3:])), nil
}

// RequestProvisioning generates a setup code for a new hardware_id
func (s *ProvisioningService) RequestProvisioning(hardwareID string) (string, time.Time, error) {
	code, err := s.GenerateSetupCode()
	if err != nil {
		return "", time.Time{}, err
	}

	expiresAt := time.Now().Add(30 * time.Minute)

	_, err = database.DB.Exec(`
		INSERT INTO gateway_provisions (id, setup_code, hardware_id, expires_at, is_claimed)
		VALUES ($1, $2, $3, $4, false)
		ON CONFLICT (setup_code) DO UPDATE SET hardware_id = $3, expires_at = $4, is_claimed = false
	`, uuid.New(), code, hardwareID, expiresAt)

	return code, expiresAt, err
}

// CheckProvisioningStatus returns if a code has been claimed and returns the config
func (s *ProvisioningService) CheckProvisioningStatus(code string) (bool, *uuid.UUID, *uuid.UUID, *string, error) {
	var farmID, coopID sql.NullString
	var isClaimed bool
	var expiresAt time.Time
	var hardwareID string

	err := database.DB.QueryRow(`
		SELECT is_claimed, farm_id, coop_id, expires_at, hardware_id
		FROM gateway_provisions 
		WHERE setup_code = $1
	`, code).Scan(&isClaimed, &farmID, &coopID, &expiresAt, &hardwareID)

	if err == sql.ErrNoRows {
		return false, nil, nil, nil, fmt.Errorf("invalid code")
	}
	if err != nil {
		return false, nil, nil, nil, err
	}

	if time.Now().After(expiresAt) {
		return false, nil, nil, nil, fmt.Errorf("code expired")
	}

	if !isClaimed {
		return false, nil, nil, nil, nil
	}

	// Code was claimed! Let's generate a permanent token for this Pi.
	fID, _ := uuid.Parse(farmID.String)
	cID, _ := uuid.Parse(coopID.String)

	// Get owner_id for the farm to associate with the token
	var ownerID uuid.UUID
	err = database.DB.QueryRow("SELECT owner_id FROM farms WHERE id = $1", fID).Scan(&ownerID)
	if err != nil {
		return true, &fID, &cID, nil, err
	}

	// Generate a unique token
	token := strings.ReplaceAll(uuid.New().String(), "-", "")
	tokenHash := utils.HashToken(token)

	// Create a Device entry if it doesn't exist (linked to the main controller)
	var deviceID uuid.UUID
	_ = database.DB.QueryRow(`
		INSERT INTO devices (id, farm_id, coop_id, device_id, name, type, firmware_version, hardware_id, is_active, is_online)
		VALUES ($1, $2, $3, $4, 'Main Gateway', 'sensor', '1.0', $5, true, true)
		ON CONFLICT (device_id) DO UPDATE SET
			coop_id = $3,
			hardware_id = $5,
			is_active = true,
			is_online = true,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id
	`, uuid.New(), fID, cID, hardwareID, hardwareID).Scan(&deviceID)

	// Save token to database
	if deviceID != uuid.Nil {
		_, _ = database.DB.Exec(`UPDATE gateway_tokens SET is_active = false WHERE device_id = $1`, deviceID)
	}
	_, err = database.DB.Exec(`
		INSERT INTO gateway_tokens (farm_id, device_id, user_id, token_hash, name, is_active)
		VALUES ($1, $2, $3, $4, $5, true)
	`, fID, deviceID, ownerID, tokenHash, "Gateway ("+hardwareID+")")

	if err != nil {
		return true, &fID, &cID, nil, err
	}

	return true, &fID, &cID, &token, nil
}

// ClaimGateway associates a setup code with a farm/coop
func (s *ProvisioningService) ClaimGateway(userID, farmID, coopID uuid.UUID, setupCode string) error {
	// 1. Verify user owns the farm
	if err := s.farmService.CheckAccess(userID, farmID, "farmer"); err != nil {
		return err
	}

	// 2. Check if code exists and is not expired
	var exists bool
	var expiresAt time.Time
	err := database.DB.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM gateway_provisions WHERE setup_code = $1), expires_at 
		FROM gateway_provisions WHERE setup_code = $1
		GROUP BY expires_at
	`, setupCode).Scan(&exists, &expiresAt)

	if err != nil || !exists {
		return fmt.Errorf("invalid setup code")
	}
	if time.Now().After(expiresAt) {
		return fmt.Errorf("setup code has expired")
	}

	// 3. Update provision record
	_, err = database.DB.Exec(`
		UPDATE gateway_provisions 
		SET is_claimed = true, farm_id = $1, coop_id = $2
		WHERE setup_code = $3
	`, farmID, coopID, setupCode)

	return err
}
