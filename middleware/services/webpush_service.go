package services

import (
	"encoding/json"
	"errors"
	"log"
	"middleware/config"
	"middleware/database"
	"middleware/schemas"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/google/uuid"
)

var (
	ErrWebPushNotConfigured = errors.New("web_push_not_configured")
	ErrSubscriptionFailed   = errors.New("subscription_failed")
)

type WebPushService struct{}

func NewWebPushService() *WebPushService {
	// Auto-generate keys if they are missing in dev, but warn the user.
	if config.AppConfig.VapidPublicKey == "" || config.AppConfig.VapidPrivateKey == "" {
		privateKey, publicKey, err := webpush.GenerateVAPIDKeys()
		if err == nil {
			log.Println("⚠️  WARNING: VAPID keys were not found in .env!")
			log.Println("⚠️  Auto-generated temporary keys for this session.")
			log.Printf("⚠️  VAPID_PUBLIC_KEY=%s\n", publicKey)
			log.Printf("⚠️  VAPID_PRIVATE_KEY=%s\n", privateKey)
			config.AppConfig.VapidPublicKey = publicKey
			config.AppConfig.VapidPrivateKey = privateKey
		}
	}
	return &WebPushService{}
}

// GetPublicKey returns the VAPID public key
func (s *WebPushService) GetPublicKey() string {
	return config.AppConfig.VapidPublicKey
}

// Subscribe saves a browser push subscription to the database
func (s *WebPushService) Subscribe(userID uuid.UUID, req schemas.WebPushSubscriptionReq, userAgent string) error {
	id := uuid.New()
	
	// Upsert based on endpoint (a browser device)
	_, err := database.DB.Exec(`
		INSERT INTO web_push_subscriptions (id, user_id, endpoint, p256dh, auth, user_agent, created_at, last_used)
		VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		ON CONFLICT (endpoint) 
		DO UPDATE SET 
			user_id = EXCLUDED.user_id,
			p256dh = EXCLUDED.p256dh,
			auth = EXCLUDED.auth,
			user_agent = EXCLUDED.user_agent,
			last_used = CURRENT_TIMESTAMP
	`, id, userID, req.Endpoint, req.Keys.P256dh, req.Keys.Auth, userAgent)
	
	if err != nil {
		log.Printf("Failed to save push subscription: %v", err)
		return ErrSubscriptionFailed
	}
	return nil
}

// Unsubscribe removes a browser push subscription
func (s *WebPushService) Unsubscribe(userID uuid.UUID, endpoint string) error {
	_, err := database.DB.Exec(`
		DELETE FROM web_push_subscriptions WHERE user_id = $1 AND endpoint = $2
	`, userID, endpoint)
	return err
}

// SendPushToUser sends a Web Push notification to all registered devices for a user
func (s *WebPushService) SendPushToUser(userID uuid.UUID, title, body, url string) error {
	if config.AppConfig.VapidPrivateKey == "" {
		return ErrWebPushNotConfigured
	}

	rows, err := database.DB.Query(`
		SELECT endpoint, p256dh, auth FROM web_push_subscriptions WHERE user_id = $1
	`, userID)
	
	if err != nil {
		return err
	}
	defer rows.Close()

	payload, _ := json.Marshal(map[string]interface{}{
		"title": title,
		"body":  body,
		"url":   url,
		"icon":  "/assets/images/tokkatot logo-02.png",
		"badge": "/assets/images/tokkatot logo-02.png",
	})

	for rows.Next() {
		var sub webpush.Subscription
		if err := rows.Scan(&sub.Endpoint, &sub.Keys.P256dh, &sub.Keys.Auth); err != nil {
			continue
		}

		// Send notification
		resp, err := webpush.SendNotification(payload, &sub, &webpush.Options{
			Subscriber:      config.AppConfig.VapidSubject, // Admin email
			VAPIDPublicKey:  config.AppConfig.VapidPublicKey,
			VAPIDPrivateKey: config.AppConfig.VapidPrivateKey,
			TTL:             43200, // 12 hours
		})
		
		if err != nil || resp.StatusCode >= 400 {
			// If subscription is expired or invalid (410 or 404), remove it
			if resp != nil && (resp.StatusCode == 410 || resp.StatusCode == 404) {
				database.DB.Exec("DELETE FROM web_push_subscriptions WHERE endpoint = $1", sub.Endpoint)
			}
		} else {
			// Update last_used
			database.DB.Exec("UPDATE web_push_subscriptions SET last_used = CURRENT_TIMESTAMP WHERE endpoint = $1", sub.Endpoint)
		}
		
		if resp != nil {
			resp.Body.Close()
		}
	}
	
	return nil
}
