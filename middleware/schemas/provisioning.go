package schemas

import "github.com/google/uuid"

type ProvisionRequest struct {
	HardwareID string `json:"hardware_id"`
}

type ProvisionResponse struct {
	SetupCode string `json:"setup_code"`
	ExpiresAt string `json:"expires_at"`
}

type ProvisionStatusResponse struct {
	IsClaimed bool       `json:"is_claimed"`
	FarmID    *uuid.UUID `json:"farm_id,omitempty"`
	CoopID    *uuid.UUID `json:"coop_id,omitempty"`
	Token     *string    `json:"token,omitempty"` // Hashed token/API Key
}

type ClaimGatewayRequest struct {
	SetupCode string    `json:"setup_code"`
	CoopID    uuid.UUID `json:"coop_id"`
}
