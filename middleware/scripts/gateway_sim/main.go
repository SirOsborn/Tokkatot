package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type DeviceReportItem struct {
	Type            string  `json:"type"`
	Model           string  `json:"model"`
	Name            string  `json:"name"`
	Active          *bool   `json:"active,omitempty"`
	FirmwareVersion *string `json:"firmware_version,omitempty"`
}

type HTTPResponse struct {
	Status int
	Body   map[string]interface{}
}

var defaultDevices = []DeviceReportItem{
	{Type: "relay", Model: "feeder_motor", Name: "Feeder Motor"},
	{Type: "relay", Model: "conveyor_belt", Name: "Conveyor Belt"},
	{Type: "relay", Model: "fan", Name: "Cooling Fan"},
	{Type: "relay", Model: "heater", Name: "Heater"},
	{Type: "sensor", Model: "temp_humidity", Name: "Temp/Humidity"},
	{Type: "sensor", Model: "water_level", Name: "Water Level"},
}

func httpJSON(method, url string, token string, body interface{}) (HTTPResponse, error) {
	var buf io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			return HTTPResponse{}, err
		}
		buf = bytes.NewBuffer(raw)
	}
	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		return HTTPResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return HTTPResponse{}, err
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	var parsed map[string]interface{}
	if len(raw) > 0 {
		_ = json.Unmarshal(raw, &parsed)
	}
	return HTTPResponse{Status: resp.StatusCode, Body: parsed}, nil
}

func getData(payload map[string]interface{}) map[string]interface{} {
	if payload == nil {
		return map[string]interface{}{}
	}
	if data, ok := payload["data"].(map[string]interface{}); ok {
		return data
	}
	return payload
}

func login(baseURL, email, phone, password string) (string, error) {
	payload := map[string]interface{}{
		"password": password,
	}
	if email != "" {
		payload["email"] = email
	}
	if phone != "" {
		payload["phone"] = phone
	}
	resp, err := httpJSON("POST", baseURL+"/v1/auth/login", "", payload)
	if err != nil {
		return "", err
	}
	if resp.Status != 200 {
		return "", fmt.Errorf("login failed (%d): %v", resp.Status, resp.Body)
	}
	data := getData(resp.Body)
	token, _ := data["access_token"].(string)
	if token == "" {
		return "", errors.New("login response missing access_token")
	}
	return token, nil
}

func pickFarm(baseURL, token, farmID, farmName string, createFarm bool) (string, error) {
	if farmID != "" {
		return farmID, nil
	}
	resp, err := httpJSON("GET", baseURL+"/v1/farms", token, nil)
	if err != nil {
		return "", err
	}
	if resp.Status != 200 {
		return "", fmt.Errorf("list farms failed (%d): %v", resp.Status, resp.Body)
	}
	data := getData(resp.Body)
	items, _ := data["data"].([]interface{})
	if len(items) > 0 {
		if first, ok := items[0].(map[string]interface{}); ok {
			if id, ok := first["id"].(string); ok {
				return id, nil
			}
		}
	}
	if !createFarm {
		return "", errors.New("no farms found, use --create-farm")
	}
	payload := map[string]interface{}{"name": farmName}
	resp, err = httpJSON("POST", baseURL+"/v1/farms", token, payload)
	if err != nil {
		return "", err
	}
	if resp.Status != 200 && resp.Status != 201 {
		return "", fmt.Errorf("create farm failed (%d): %v", resp.Status, resp.Body)
	}
	data = getData(resp.Body)
	if id, ok := data["id"].(string); ok {
		return id, nil
	}
	return "", errors.New("create farm response missing id")
}

func pickCoop(baseURL, token, farmID, coopID, coopName string, coopNumber int, createCoop bool) (string, error) {
	if coopID != "" {
		return coopID, nil
	}
	resp, err := httpJSON("GET", fmt.Sprintf("%s/v1/farms/%s/coops", baseURL, farmID), token, nil)
	if err != nil {
		return "", err
	}
	if resp.Status != 200 {
		return "", fmt.Errorf("list coops failed (%d): %v", resp.Status, resp.Body)
	}
	data := getData(resp.Body)
	items, _ := data["coops"].([]interface{})
	if len(items) > 0 {
		if first, ok := items[0].(map[string]interface{}); ok {
			if id, ok := first["id"].(string); ok {
				return id, nil
			}
		}
	}
	if !createCoop {
		return "", errors.New("no coops found, use --create-coop")
	}
	payload := map[string]interface{}{
		"number": coopNumber,
		"name":   coopName,
	}
	resp, err = httpJSON("POST", fmt.Sprintf("%s/v1/farms/%s/coops", baseURL, farmID), token, payload)
	if err != nil {
		return "", err
	}
	if resp.Status != 200 && resp.Status != 201 {
		return "", fmt.Errorf("create coop failed (%d): %v", resp.Status, resp.Body)
	}
	data = getData(resp.Body)
	if id, ok := data["id"].(string); ok {
		return id, nil
	}
	return "", errors.New("create coop response missing id")
}

func reportDevices(baseURL, token, farmID, coopID, hardwareID string, devices []DeviceReportItem, inactiveSet map[string]bool) error {
	// Modify the payload to include all required fields based on server-side expectations
	report := make([]map[string]interface{}, 0, len(devices))
	for _, d := range devices {
		item := map[string]interface{}{
			"type":   d.Type,
			"model":  d.Model,
			"name":   d.Name,
			"active": true, // Default to active unless specified in inactiveSet
		}
		if inactiveSet[d.Model] {
			item["active"] = false
		}
		report = append(report, item)
	}

	payload := map[string]interface{}{
		"hardware_id": hardwareID,
		"devices":     report,
	}
	url := fmt.Sprintf("%s/v1/farms/%s/coops/%s/devices/report", baseURL, farmID, coopID)

	resp, err := httpJSON("POST", url, token, payload)
	if err != nil {
		return err
	}
	if resp.Status != 200 && resp.Status != 201 {
		return fmt.Errorf("device report failed (%d): %v", resp.Status, resp.Body)
	}
	return nil
}

func sendTelemetry(baseURL, token, farmID, coopID, hardwareID string, temp, humidity, water float64) error {
	payload := map[string]interface{}{
		"hardware_id": hardwareID,
		"timestamp":   time.Now().UTC().Format(time.RFC3339),
		"sensors": map[string]interface{}{
			"temperature_c":   temp,
			"humidity_pct":    humidity,
			"water_level_raw": water,
		},
	}
	url := fmt.Sprintf("%s/v1/farms/%s/coops/%s/telemetry", baseURL, farmID, coopID)
	resp, err := httpJSON("POST", url, token, payload)
	if err != nil {
		return err
	}
	if resp.Status != 200 && resp.Status != 201 {
		return fmt.Errorf("telemetry failed (%d): %v", resp.Status, resp.Body)
	}
	return nil
}

func listDevices(baseURL, token, farmID, coopID string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/v1/farms/%s/devices?coop_id=%s", baseURL, farmID, coopID)
	resp, err := httpJSON("GET", url, token, nil)
	if err != nil {
		return nil, err
	}
	if resp.Status != 200 {
		return nil, fmt.Errorf("list devices failed (%d): %v", resp.Status, resp.Body)
	}
	return resp.Body, nil
}

func pickActuatorDevice(devicesPayload map[string]interface{}) (map[string]interface{}, error) {
	data := getData(devicesPayload)
	items, _ := data["devices"].([]interface{})
	for _, it := range items {
		dev, ok := it.(map[string]interface{})
		if !ok {
			continue
		}
		model, _ := dev["model"].(string)
		model = strings.ToLower(model)
		name, _ := dev["name"].(string)
		name = strings.ToLower(name)
		devType, _ := dev["type"].(string)
		devType = strings.ToLower(devType)
		if model == "feeder_motor" || model == "conveyor_belt" || model == "fan" || model == "heater" ||
			strings.Contains(name, "feeder") || strings.Contains(name, "conveyor") || strings.Contains(name, "fan") || strings.Contains(name, "heater") ||
			devType == "relay" {
			return dev, nil
		}
	}
	return nil, errors.New("no actuator device found")
}

func createSchedule(baseURL, token, farmID, coopID, deviceID string) (string, error) {
	payload := map[string]interface{}{
		"device_id":       deviceID,
		"coop_id":         coopID,
		"name":            "Gateway E2E Test",
		"schedule_type":   "time_based",
		"cron_expression": "0 6 * * *",
		"action":          "on",
		"action_duration": 60,
	}
	url := fmt.Sprintf("%s/v1/farms/%s/schedules", baseURL, farmID)
	resp, err := httpJSON("POST", url, token, payload)
	if err != nil {
		return "", err
	}
	if resp.Status != 200 && resp.Status != 201 {
		return "", fmt.Errorf("create schedule failed (%d): %v", resp.Status, resp.Body)
	}
	data := getData(resp.Body)
	id, _ := data["id"].(string)
	if id == "" {
		return "", errors.New("schedule response missing id")
	}
	return id, nil
}

func executeScheduleNow(baseURL, token, farmID, scheduleID string) error {
	url := fmt.Sprintf("%s/v1/farms/%s/schedules/%s/execute-now", baseURL, farmID, scheduleID)
	resp, err := httpJSON("POST", url, token, map[string]interface{}{})
	if err != nil {
		return err
	}
	if resp.Status != 200 {
		return fmt.Errorf("execute schedule failed (%d): %v", resp.Status, resp.Body)
	}
	return nil
}

func getDeviceCommands(baseURL, token, farmID, deviceID string) error {
	url := fmt.Sprintf("%s/v1/farms/%s/devices/%s/commands", baseURL, farmID, deviceID)
	resp, err := httpJSON("GET", url, token, nil)
	if err != nil {
		return err
	}
	if resp.Status != 200 {
		return fmt.Errorf("list device commands failed (%d): %v", resp.Status, resp.Body)
	}
	return nil
}

func getTemperatureTimeline(baseURL, token, farmID, coopID string) error {
	url := fmt.Sprintf("%s/v1/farms/%s/coops/%s/temperature-timeline", baseURL, farmID, coopID)
	resp, err := httpJSON("GET", url, token, nil)
	if err != nil {
		return err
	}
	if resp.Status != 200 {
		return fmt.Errorf("temperature timeline failed (%d): %v", resp.Status, resp.Body)
	}
	return nil
}

func getAlerts(baseURL, token, farmID string) error {
	url := fmt.Sprintf("%s/v1/farms/%s/alerts", baseURL, farmID)
	resp, err := httpJSON("GET", url, token, nil)
	if err != nil {
		return err
	}
	if resp.Status != 200 {
		return fmt.Errorf("get alerts failed (%d): %v", resp.Status, resp.Body)
	}
	return nil
}

func parseList(value string) map[string]bool {
	out := map[string]bool{}
	if strings.TrimSpace(value) == "" {
		return out
	}
	for _, item := range strings.Split(value, ",") {
		item = strings.TrimSpace(item)
		if item != "" {
			out[item] = true
		}
	}
	return out
}

func logf(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

func main() {
	baseURL := flag.String("base-url", os.Getenv("CLOUD_API_URL"), "Cloud API base URL (defaults to CLOUD_API_URL environment variable)")
	if *baseURL == "" {
		*baseURL = "http://127.0.0.1:3000"
	}
	email := flag.String("email", "test1@tokkatot.com", "")
	phone := flag.String("phone", "", "")
	password := flag.String("password", "TokkatotTest2026!", "")
	farmID := flag.String("farm-id", "", "")
	farmName := flag.String("farm-name", "Demo Farm", "")
	coopID := flag.String("coop-id", "", "")
	coopName := flag.String("coop-name", "Coop 1", "")
	coopNumber := flag.Int("coop-number", 1, "")
	createFarm := flag.Bool("create-farm", false, "")
	createCoop := flag.Bool("create-coop", false, "")
	hardwareID := flag.String("hardware-id", "", "")
	interval := flag.Int("interval", 10, "")
	once := flag.Bool("once", false, "")
	devices := flag.String("devices", "", "")
	inactive := flag.String("inactive", "", "")
	tempBase := flag.Float64("temp-base", 29.5, "")
	tempVar := flag.Float64("temp-var", 0.6, "")
	humidBase := flag.Float64("humidity-base", 62.0, "")
	humidVar := flag.Float64("humidity-var", 2.0, "")
	waterBase := flag.Float64("water-base", 1500.0, "")
	waterVar := flag.Float64("water-var", 30.0, "")
	pollSchedules := flag.Bool("poll-schedules", false, "")
	live := flag.Bool("live", false, "")
	e2e := flag.Bool("e2e", false, "")
	flag.Parse()

	if *password == "" || (*email == "" && *phone == "") {
		logf("❌ Provide --password and either --email or --phone")
		os.Exit(1)
	}

	if *hardwareID == "" {
		*hardwareID = fmt.Sprintf("SIM-ESP32-%d", time.Now().UnixNano())
	}

	deviceList := defaultDevices
	if strings.TrimSpace(*devices) != "" {
		allow := parseList(*devices)
		filtered := make([]DeviceReportItem, 0)
		for _, d := range deviceList {
			if allow[d.Model] {
				filtered = append(filtered, d)
			}
		}
		deviceList = filtered
	}
	inactiveSet := parseList(*inactive)

	token, err := login(*baseURL, *email, *phone, *password)
	if err != nil {
		logf("❌ %v", err)
		os.Exit(1)
	}

	fID, err := pickFarm(*baseURL, token, *farmID, *farmName, *createFarm)
	if err != nil {
		logf("❌ %v", err)
		os.Exit(1)
	}
	cID, err := pickCoop(*baseURL, token, fID, *coopID, *coopName, *coopNumber, *createCoop)
	if err != nil {
		logf("❌ %v", err)
		os.Exit(1)
	}

	logf("✅ Auth OK | farm=%s coop=%s hardware_id=%s", fID, cID, *hardwareID)

	if err := reportDevices(*baseURL, token, fID, cID, *hardwareID, deviceList, inactiveSet); err != nil {
		logf("❌ %v", err)
		os.Exit(1)
	}
	logf("✅ Device report sent")

	if *e2e {
		devs, err := listDevices(*baseURL, token, fID, cID)
		if err != nil {
			logf("❌ %v", err)
			os.Exit(1)
		}
		actuator, err := pickActuatorDevice(devs)
		if err != nil {
			logf("❌ %v", err)
			os.Exit(1)
		}
		deviceID, _ := actuator["id"].(string)
		if deviceID == "" {
			logf("❌ actuator missing id")
			os.Exit(1)
		}
		scheduleID, err := createSchedule(*baseURL, token, fID, cID, deviceID)
		if err != nil {
			logf("❌ %v", err)
			os.Exit(1)
		}
		logf("✅ Schedule created | id=%s", scheduleID)

		if err := executeScheduleNow(*baseURL, token, fID, scheduleID); err != nil {
			logf("❌ %v", err)
			os.Exit(1)
		}
		logf("✅ Schedule execute-now issued")

		if err := getDeviceCommands(*baseURL, token, fID, deviceID); err != nil {
			logf("❌ %v", err)
			os.Exit(1)
		}
		logf("✅ Device commands fetched")

		if err := sendTelemetry(*baseURL, token, fID, cID, *hardwareID, *tempBase, *humidBase, *waterBase); err != nil {
			logf("❌ %v", err)
			os.Exit(1)
		}
		logf("✅ Telemetry sent (E2E)")

		if err := getTemperatureTimeline(*baseURL, token, fID, cID); err != nil {
			logf("❌ %v", err)
			os.Exit(1)
		}
		logf("✅ Temperature timeline fetched")

		if err := getAlerts(*baseURL, token, fID); err != nil {
			logf("❌ %v", err)
			os.Exit(1)
		}
		logf("✅ Alerts fetched")
		return
	}

	rand.Seed(time.Now().UnixNano())
	for {
		temp := *tempBase + (rand.Float64()*2-1)*(*tempVar)
		humid := *humidBase + (rand.Float64()*2-1)*(*humidVar)
		water := *waterBase + (rand.Float64()*2-1)*(*waterVar)

		if err := sendTelemetry(*baseURL, token, fID, cID, *hardwareID, temp, humid, water); err != nil {
			logf("❌ %v", err)
			os.Exit(1)
		}
		logf("📡 Telemetry sent | temp=%.1fC humidity=%.1f%% water=%.0f", temp, humid, water)

		if *pollSchedules {
			if _, err := listDevices(*baseURL, token, fID, cID); err != nil {
				logf("❌ %v", err)
				os.Exit(1)
			}
			logf("🗓️  Schedules fetched")
		}

		if *live {
			if err := getTemperatureTimeline(*baseURL, token, fID, cID); err != nil {
				logf("❌ %v", err)
				os.Exit(1)
			}
			logf("📊 Monitoring updated")
		}

		if *once {
			break
		}
		time.Sleep(time.Duration(*interval) * time.Second)
	}
}
