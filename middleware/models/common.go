package models

import (
	"encoding/json"
	"fmt"
	"database/sql/driver"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// NullRawMessage is a json.RawMessage that can be scanned from a nullable SQL JSONB column.
type NullRawMessage []byte

func (n *NullRawMessage) Scan(src interface{}) error {
	if src == nil {
		*n = nil
		return nil
	}
	switch v := src.(type) {
	case []byte:
		data := make([]byte, len(v))
		copy(data, v)
		*n = data
	case string:
		*n = []byte(v)
	default:
		return fmt.Errorf("NullRawMessage: unsupported type %T", src)
	}
	return nil
}

func (n NullRawMessage) MarshalJSON() ([]byte, error) {
	if len(n) == 0 {
		return []byte("null"), nil
	}
	return json.RawMessage(n).MarshalJSON()
}

func (n NullRawMessage) Value() (driver.Value, error) {
	if len(n) == 0 {
		return nil, nil
	}
	return []byte(n), nil
}

// JWTClaims for JWT token claims
type JWTClaims struct {
	UserID    uuid.UUID `json:"sub"`
	Email     *string   `json:"email,omitempty"`
	Phone     *string   `json:"phone,omitempty"`
	FarmID    uuid.UUID `json:"farm_id"`
	Role      string    `json:"role"`
	IssuedAt  int64     `json:"iat"`
	ExpiresAt int64     `json:"exp"`
}

func (j JWTClaims) Valid() error {
	if j.ExpiresAt < time.Now().Unix() {
		return jwt.NewValidationError("token expired", jwt.ValidationErrorExpired)
	}
	return nil
}
