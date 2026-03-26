package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// Configuration
var (
	CloudBaseURL = os.Getenv("CLOUD_BASE_URL")
	HardwareID   = os.Getenv("HARDWARE_ID")
	FarmID       = os.Getenv("FARM_ID")
	CoopID       = os.Getenv("COOP_ID")
	CloudToken   = os.Getenv("CLOUD_TOKEN") // Should be fetched via Login or Provisioning
)

type TelemetryPayload struct {
	HardwareID string                 `json:"hardware_id"`
	Timestamp  string                 `json:"timestamp"`
	Sensors    map[string]interface{} `json:"sensors"`
}

func main() {
	log.Printf("🚀 Starting Tokkatot Gateway Service [%s]", HardwareID)

	if CloudBaseURL == "" {
		log.Fatal("❌ CLOUD_BASE_URL not set")
	}

	// Disable TLS verification for local ESP32 self-signed certs (production should use real certs/keys)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr, Timeout: 5 * time.Second}

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 1. Poll ESP32 (Assuming static IP for now, but should be discovered)
			espURL := "https://10.0.0.2/get-current-data"
			resp, err := client.Get(espURL)
			if err != nil {
				log.Printf("⚠️  Failed to poll ESP32: %v", err)
				continue
			}
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)
			var espData map[string]interface{}
			if err := json.Unmarshal(body, &espData); err != nil {
				log.Printf("⚠️  Failed to parse ESP32 data: %v", err)
				continue
			}

			// 2. Prepare payload for Cloud
			payload := TelemetryPayload{
				HardwareID: HardwareID,
				Timestamp:  time.Now().UTC().Format(time.RFC3339),
				Sensors: map[string]interface{}{
					"temperature_c":   espData["temperature"],
					"humidity_pct":    espData["humidity"],
					"water_level_raw": espData["water_level"],
				},
			}

			// 3. Post to Cloud
			postToCloud(payload)
		}
	}
}

func postToCloud(payload TelemetryPayload) {
	url := fmt.Sprintf("%s/v1/farms/%s/coops/%s/telemetry", CloudBaseURL, FarmID, CoopID)
	
	raw, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(raw))
	if err != nil {
		log.Printf("❌ Failed to create cloud request: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	if CloudToken != "" {
		req.Header.Set("Authorization", "Bearer "+CloudToken)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("⚠️  Cloud communication error: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		log.Printf("⚠️  Cloud rejected telemetry: %d", resp.StatusCode)
		return
	}

	log.Printf("📡 Telemetry uploaded to Cloud")
}
