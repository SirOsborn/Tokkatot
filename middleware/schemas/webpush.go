package schemas

import "time"

// WebPushSubscriptionKeys represents the keys required for Web Push
type WebPushSubscriptionKeys struct {
	P256dh string `json:"p256dh" validate:"required"`
	Auth   string `json:"auth" validate:"required"`
}

// WebPushSubscriptionReq represents a subscription sent from the browser
type WebPushSubscriptionReq struct {
	Endpoint string                  `json:"endpoint" validate:"required"`
	Keys     WebPushSubscriptionKeys `json:"keys" validate:"required"`
}

// WebPushSubscriptionRes represents a stored subscription
type WebPushSubscriptionRes struct {
	ID        string    `json:"id"`
	Endpoint  string    `json:"endpoint"`
	UserAgent string    `json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
	LastUsed  time.Time `json:"last_used"`
}

// VapidPublicKeyRes returns the VAPID public key
type VapidPublicKeyRes struct {
	PublicKey string `json:"public_key"`
}
