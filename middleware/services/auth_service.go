package services

import (
	"database/sql"
	"errors"
	"middleware/database"
	"middleware/models"
	"middleware/schemas"
	"middleware/utils"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound      = errors.New("user_not_found")
	ErrInvalidPassword   = errors.New("invalid_password")
	ErrDuplicateEmail    = errors.New("email_already_exists")
	ErrDuplicatePhone    = errors.New("phone_already_exists")
	ErrInvalidRegKey     = errors.New("invalid_reg_key")
	ErrKeyUsed           = errors.New("reg_key_used")
	ErrKeyExpired        = errors.New("reg_key_expired")
	ErrInvalidFarmerID   = errors.New("invalid_farmer_id")
	ErrMissingCredentials = errors.New("reg_key_req")
	ErrAccountInactive   = errors.New("account_inactive")
)

// AuthService handles all business logic related to authentication
type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

// Signup creates a new farmer user and returns user + token
func (s *AuthService) Signup(req schemas.SignupRequest) (*models.User, string, error) {
	// Check if email or phone already exists
	var count int
	database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", req.Email).Scan(&count)
	if count > 0 {
		return nil, "", ErrDuplicateEmail
	}
	database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE phone = $1", req.Phone).Scan(&count)
	if count > 0 {
		return nil, "", ErrDuplicatePhone
	}

	// Validate registration key and get associated farm_id (if any)
	var regKeyID uuid.UUID
	var farmID uuid.NullUUID
	var nationalID, fullName, sex, province *string
	var farmName sql.NullString
	if req.RegistrationKey != nil && *req.RegistrationKey != "" {
		var isUsed bool
		var expiresAt time.Time
		err := database.DB.QueryRow(`
			SELECT id, farm_id, farm_name, is_used, expires_at, national_id_number, full_name, sex, province
			FROM registration_keys 
			WHERE key_code = $1
		`, *req.RegistrationKey).Scan(&regKeyID, &farmID, &farmName, &isUsed, &expiresAt, &nationalID, &fullName, &sex, &province)
		
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, "", ErrInvalidRegKey
			}
			return nil, "", err
		}
		
		if isUsed {
			return nil, "", ErrKeyUsed
		}
		if expiresAt.Before(time.Now()) {
			return nil, "", ErrKeyExpired
		}
	} else if req.FarmerID != nil && *req.FarmerID != "" {
		farmerUserID, parseErr := uuid.Parse(*req.FarmerID)
		if parseErr != nil {
			return nil, "", ErrInvalidFarmerID
		}

		err := database.DB.QueryRow(`
			SELECT farm_id FROM farm_users
			WHERE user_id = $1 AND role = 'farmer' AND is_active = true
			LIMIT 1
		`, farmerUserID).Scan(&farmID)

		if err == sql.ErrNoRows {
			return nil, "", ErrInvalidFarmerID
		}
	} else {
		return nil, "", ErrMissingCredentials
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	// Create user transaction
	tx, err := database.DB.Begin()
	if err != nil {
		return nil, "", err
	}
	defer tx.Rollback()

	userID := uuid.New()
	now := time.Now()

	// Use request data if provided, otherwise fallback to key data
	finalNationalID := nationalID
	if req.NationalIDNumber != nil && *req.NationalIDNumber != "" {
		finalNationalID = req.NationalIDNumber
	}

	finalSex := sex
	if req.Sex != nil && *req.Sex != "" {
		finalSex = req.Sex
	}

	finalProvince := province
	if req.Province != nil && *req.Province != "" {
		finalProvince = req.Province
	}

	// fullName from key is used as fallback for users.full_name
	finalFullName := fullName
	if finalFullName == nil || *finalFullName == "" {
		finalFullName = &req.Name
	}

	_, err = tx.Exec(`
		INSERT INTO users (id, email, phone, name, password_hash, national_id_number, sex, province, full_name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $10)
	`, userID, req.Email, req.Phone, req.Name, string(hashedPassword), finalNationalID, finalSex, finalProvince, finalFullName, now)
	if err != nil {
		return nil, "", err
	}

	if regKeyID != uuid.Nil {
		_, err = tx.Exec("UPDATE registration_keys SET is_used = true, used_by_user_id = $1, used_at = $2 WHERE id = $3", userID, now, regKeyID)
		if err != nil {
			return nil, "", err
		}
	}

	// If key/farmer is attached to a farm, link user to it
	// If it's a new farmer (reg key used) and no farm exists yet, create one
	if regKeyID != uuid.Nil && !farmID.Valid {
		fID := uuid.New()
		fName := "My Farm"
		if farmName.Valid && farmName.String != "" {
			fName = farmName.String
		}
		
		_, err = tx.Exec(`
			INSERT INTO farms (id, owner_id, name, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $4)
		`, fID, userID, fName, now)
		if err != nil {
			return nil, "", err
		}
		farmID.UUID = fID
		farmID.Valid = true
	}

	if farmID.Valid {
		role := "farmer"
		if req.FarmerID != nil && *req.FarmerID != "" {
			role = "worker"
		}
		_, err = tx.Exec(`
			INSERT INTO farm_users (id, farm_id, user_id, role, invited_by, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $6)
		`, uuid.New(), farmID.UUID, userID, role, userID, now) // self-invited for first farmer
		if err != nil {
			return nil, "", err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, "", err
	}

	// Fetch created user with farm info
	user, err := s.GetProfile(userID)
	if err != nil {
		return nil, "", err
	}

	// In real app, generate JWT here
	token, err := utils.GenerateAccessToken(userID, req.Email, req.Phone, farmID.UUID, "farmer")
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

// Login verifies credentials and returns user info + token
func (s *AuthService) Login(identifier, password string) (*models.User, string, string, error) {
	var user models.User
	var passwordHash string
	var role = "farmer"
	var farmID uuid.UUID

	// Try admin table first for role prioritization
	err := database.DB.QueryRow(`
		SELECT id, email, phone, name, password_hash, is_active 
		FROM admins WHERE email = $1 OR phone = $1
	`, identifier).Scan(&user.ID, &user.Email, &user.Phone, &user.Name, &passwordHash, &user.IsActive)

	if err == nil {
		role = "admin"
	} else if err == sql.ErrNoRows {
		// Fallback to regular users table
		err = database.DB.QueryRow(`
			SELECT id, email, phone, name, password_hash, is_active 
			FROM users WHERE email = $1 OR phone = $1
		`, identifier).Scan(&user.ID, &user.Email, &user.Phone, &user.Name, &passwordHash, &user.IsActive)

		if err == sql.ErrNoRows {
			return nil, "", "", ErrUserNotFound
		}
	} else {
		return nil, "", "", err
	}

	if !user.IsActive {
		return nil, "", "", ErrAccountInactive
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return nil, "", "", ErrInvalidPassword
	}

	// Update last login
	if role == "admin" {
		database.DB.Exec("UPDATE admins SET last_login = CURRENT_TIMESTAMP WHERE id = $1", user.ID)
	} else {
		database.DB.Exec("UPDATE users SET last_login = CURRENT_TIMESTAMP WHERE id = $1", user.ID)

		// Get farm ID for farmer/worker
		database.DB.QueryRow("SELECT farm_id FROM farm_users WHERE user_id = $1 LIMIT 1", user.ID).Scan(&farmID)
	}

	// Generate JWT
	token, err := utils.GenerateAccessToken(user.ID, user.Email, user.Phone, farmID, role)
	if err != nil {
		return nil, "", "", err
	}

	// Fetch full profile with farm info for the response
	fullUser, _ := s.GetProfile(user.ID)
	if fullUser != nil {
		return fullUser, token, role, nil
	}

	return &user, token, role, nil
}
// GetProfile returns the full profile including farm info
func (s *AuthService) GetProfile(userID uuid.UUID) (*models.User, error) {
	var user models.User
	var lastLogin sql.NullTime
	var nationalID, sex, province, fullName, farmName sql.NullString
	var farmID uuid.NullUUID
	var role string

	// Try to get role from farm_users first, default to 'farmer' or 'viewer'
	database.DB.QueryRow("SELECT role FROM farm_users WHERE user_id = $1 LIMIT 1", userID).Scan(&role)
	if role == "" {
		// Might be an admin
		database.DB.QueryRow("SELECT 'admin' FROM admins WHERE id = $1", userID).Scan(&role)
	}
	user.Role = role

	err := database.DB.QueryRow(`
		SELECT u.id, u.email, u.phone, u.name, u.is_active, u.created_at, u.last_login,
		       u.national_id_number, u.sex, u.province, u.full_name,
		       f.id as farm_id, f.name as farm_name
		FROM users u
		LEFT JOIN farm_users fu ON u.id = fu.user_id
		LEFT JOIN farms f ON fu.farm_id = f.id
		WHERE u.id = $1
		LIMIT 1
	`, userID).Scan(
		&user.ID, &user.Email, &user.Phone, &user.Name, 
		&user.IsActive, &user.CreatedAt, &lastLogin,
		&nationalID, &sex, &province, &fullName,
		&farmID, &farmName,
	)

	if err == sql.ErrNoRows {
		// Check admins table if not in users
		err = database.DB.QueryRow(`
			SELECT id, email, phone, name, is_active, created_at, last_login
			FROM admins WHERE id = $1
		`, userID).Scan(
			&user.ID, &user.Email, &user.Phone, &user.Name, 
			&user.IsActive, &user.CreatedAt, &lastLogin,
		)
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		user.Role = "admin"
	} else if err != nil {
		return nil, err
	}

	if lastLogin.Valid { user.LastLogin = &lastLogin.Time }
	if nationalID.Valid { user.NationalIDNumber = &nationalID.String }
	if sex.Valid { user.Sex = &sex.String }
	if province.Valid { user.Province = &province.String }
	if fullName.Valid { user.FullName = &fullName.String }
	
	if farmID.Valid { user.FarmID = &farmID.UUID }
	if farmName.Valid { user.FarmName = &farmName.String }

	return &user, nil
}

// UpdateProfile updates user info
func (s *AuthService) UpdateProfile(userID uuid.UUID, req schemas.UpdateProfileRequest) (*models.User, error) {
	// Update users table directly with all fields
	_, err := database.DB.Exec(`
		UPDATE users 
		SET name = COALESCE($1, name),
		    email = COALESCE($2, email),
		    phone = COALESCE($3, phone),
		    national_id_number = COALESCE($4, national_id_number),
		    sex = COALESCE($5, sex),
		    province = COALESCE($6, province),
		    full_name = COALESCE($1, full_name), -- sync name to full_name if name changed
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $7
	`, req.Name, req.Email, req.Phone, req.NationalIDNumber, req.Sex, req.Province, userID)
	
	if err != nil {
		return nil, err
	}

	return s.GetProfile(userID)
}
