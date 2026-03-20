package services

import (
	"database/sql"
	"fmt"
	"middleware/database"
	"middleware/models"
	"middleware/schemas"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

type TelemetryService struct {
	farmService *FarmService
	coopService *CoopService
}

func NewTelemetryService() *TelemetryService {
	return &TelemetryService{
		farmService: NewFarmService(),
		coopService: NewCoopService(),
	}
}

type tempReading struct {
	Temp float64
	Time time.Time
}

func (s *TelemetryService) CleanupOldReadings(retentionDays int) (int64, error) {
	if retentionDays <= 0 {
		return 0, nil
	}
	cutoff := time.Now().AddDate(0, 0, -retentionDays)
	res, err := database.DB.Exec(`DELETE FROM device_readings WHERE timestamp < $1`, cutoff)
	if err != nil {
		return 0, err
	}
	affected, _ := res.RowsAffected()
	return affected, nil
}

var allowedDeviceTypes = map[string]bool{
	"gpio":  true,
	"relay": true,
	"pwm":   true,
	"adc":   true,
	"servo": true,
	"sensor": true,
}

func (s *TelemetryService) ReportCoopDevices(userID, farmID, coopID uuid.UUID, req schemas.DeviceReportRequest) ([]models.Device, error) {
	if err := s.farmService.CheckAccess(userID, farmID, "farmer"); err != nil {
		return nil, err
	}
	if _, err := s.coopService.GetCoop(userID, farmID, coopID); err != nil {
		return nil, err
	}
	if strings.TrimSpace(req.HardwareID) == "" {
		return nil, fmt.Errorf("hardware_id required")
	}

	reported := make(map[string]struct{})
	now := time.Now()
	devices := make([]models.Device, 0)

	for _, d := range req.Devices {
		model := strings.ToLower(strings.TrimSpace(d.Model))
		if model == "" {
			return nil, fmt.Errorf("device model required")
		}
		deviceType := strings.ToLower(strings.TrimSpace(d.Type))
		if deviceType == "" {
			if model == "temp_humidity" || model == "water_level" {
				deviceType = "sensor"
			} else {
				deviceType = "relay"
			}
		}
		if !allowedDeviceTypes[deviceType] {
			return nil, fmt.Errorf("invalid device type: %s", deviceType)
		}

		active := true
		if d.Active != nil {
			active = *d.Active
		}

		name := strings.TrimSpace(d.Name)
		if name == "" {
			name = strings.Title(strings.ReplaceAll(model, "_", " "))
		}

		firmware := "gateway"
		if d.FirmwareVersion != nil && strings.TrimSpace(*d.FirmwareVersion) != "" {
			firmware = strings.TrimSpace(*d.FirmwareVersion)
		}

		deviceID := fmt.Sprintf("%s:%s", req.HardwareID, model)
		reported[deviceID] = struct{}{}

		var dev models.Device
		err := database.DB.QueryRow(`
			INSERT INTO devices (
				id, farm_id, coop_id, device_id, name, type, model,
				is_main_controller, firmware_version, hardware_id,
				is_active, is_online, created_at, updated_at
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7,
				false, $8, $9,
				$10, $11, $12, $12
			)
			ON CONFLICT (device_id) DO UPDATE SET
				name = EXCLUDED.name,
				type = EXCLUDED.type,
				model = EXCLUDED.model,
				farm_id = EXCLUDED.farm_id,
				coop_id = EXCLUDED.coop_id,
				firmware_version = EXCLUDED.firmware_version,
				hardware_id = EXCLUDED.hardware_id,
				is_active = EXCLUDED.is_active,
				is_online = EXCLUDED.is_online,
				updated_at = EXCLUDED.updated_at
			RETURNING id, farm_id, coop_id, device_id, name, type, model,
				is_main_controller, firmware_version, hardware_id, location,
				is_active, is_online, last_heartbeat, created_at, updated_at
		`, uuid.New(), farmID, coopID, deviceID, name, deviceType, model, firmware, req.HardwareID, active, active, now).
			Scan(&dev.ID, &dev.FarmID, &dev.CoopID, &dev.DeviceID, &dev.Name, &dev.Type, &dev.Model,
				&dev.IsMainController, &dev.FirmwareVersion, &dev.HardwareID, &dev.Location,
				&dev.IsActive, &dev.IsOnline, &dev.LastHeartbeat, &dev.CreatedAt, &dev.UpdatedAt)
		if err != nil {
			return nil, err
		}
		devices = append(devices, dev)
	}

	// Mark missing devices inactive
	if len(reported) > 0 {
		placeholders := make([]string, 0, len(reported))
		args := []interface{}{coopID}
		i := 2
		for id := range reported {
			placeholders = append(placeholders, fmt.Sprintf("$%d", i))
			args = append(args, id)
			i++
		}
		query := fmt.Sprintf(`UPDATE devices SET is_active = false, is_online = false, updated_at = $%d
			WHERE coop_id = $1 AND device_id NOT IN (%s)`, i, strings.Join(placeholders, ","))
		args = append(args, now)
		_, _ = database.DB.Exec(query, args...)
	}

	return devices, nil
}

func (s *TelemetryService) IngestTelemetry(userID, farmID, coopID uuid.UUID, req schemas.TelemetryRequest) error {
	if err := s.farmService.CheckAccess(userID, farmID, "worker"); err != nil {
		return err
	}
	if _, err := s.coopService.GetCoop(userID, farmID, coopID); err != nil {
		return err
	}
	ts := time.Now()
	if req.Timestamp != nil {
		ts = *req.Timestamp
	}

	if req.Sensors.TemperatureC == nil && req.Sensors.HumidityPct == nil && req.Sensors.WaterLevel == nil {
		return nil
	}

	tempDeviceID, _ := s.ensureSensorDevice(farmID, coopID, req.HardwareID, "temp_humidity", "Temp/Humidity")
	waterDeviceID, _ := s.ensureSensorDevice(farmID, coopID, req.HardwareID, "water_level", "Water Level")

	if req.Sensors.TemperatureC != nil && tempDeviceID != nil {
		_ = insertReading(*tempDeviceID, "temperature", *req.Sensors.TemperatureC, "C", ts)
	}
	if req.Sensors.HumidityPct != nil && tempDeviceID != nil {
		_ = insertReading(*tempDeviceID, "humidity", *req.Sensors.HumidityPct, "%", ts)
	}
	if req.Sensors.WaterLevel != nil && waterDeviceID != nil {
		_ = insertReading(*waterDeviceID, "water_level", *req.Sensors.WaterLevel, "raw", ts)
	}

	// Update heartbeats
	if tempDeviceID != nil {
		_, _ = database.DB.Exec("UPDATE devices SET is_online = true, last_heartbeat = $1, updated_at = $1 WHERE id = $2", ts, *tempDeviceID)
	}
	if waterDeviceID != nil {
		_, _ = database.DB.Exec("UPDATE devices SET is_online = true, last_heartbeat = $1, updated_at = $1 WHERE id = $2", ts, *waterDeviceID)
	}

	if req.Sensors.WaterLevel != nil {
		_ = s.checkWaterAlert(farmID, coopID, *req.Sensors.WaterLevel, ts)
	}

	return nil
}

func (s *TelemetryService) checkWaterAlert(farmID, coopID uuid.UUID, latest float64, ts time.Time) error {
	var threshold sql.NullFloat64
	err := database.DB.QueryRow(`SELECT water_level_half_threshold FROM coops WHERE id = $1`, coopID).Scan(&threshold)
	if err != nil || !threshold.Valid {
		return nil
	}
	if latest >= threshold.Float64 {
		return nil
	}

	oneMinAgo := ts.Add(-1 * time.Minute)
	var maxVal sql.NullFloat64
	err = database.DB.QueryRow(`
		SELECT MAX(dr.value) FROM device_readings dr
		JOIN devices d ON dr.device_id = d.id
		WHERE d.coop_id = $1 AND dr.sensor_type = 'water_level' AND dr.timestamp >= $2
	`, coopID, oneMinAgo).Scan(&maxVal)
	if err != nil || !maxVal.Valid {
		return nil
	}
	if maxVal.Float64 >= threshold.Float64 {
		return nil
	}

	var exists bool
	_ = database.DB.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM alerts 
			WHERE coop_id = $1 AND alert_type = 'water_level_low' AND is_active = true
		)
	`, coopID).Scan(&exists)
	if exists {
		return nil
	}

	msg := "Water level below half for 1 minute"
	_, err = database.DB.Exec(`
		INSERT INTO alerts (id, farm_id, coop_id, alert_type, severity, message, threshold_value, actual_value, is_active, is_acknowledged, triggered_at, created_at)
		VALUES ($1,$2,$3,'water_level_low','warning',$4,$5,$6,true,false,$7,$7)
	`, uuid.New(), farmID, coopID, msg, threshold.Float64, latest, ts)
	return err
}

func (s *TelemetryService) ensureSensorDevice(farmID, coopID uuid.UUID, hardwareID, model, name string) (*uuid.UUID, error) {
	if model == "" {
		return nil, nil
	}
	var id uuid.UUID
	err := database.DB.QueryRow(`SELECT id FROM devices WHERE coop_id = $1 AND model = $2 LIMIT 1`, coopID, model).Scan(&id)
	if err == nil {
		return &id, nil
	}
	if err != sql.ErrNoRows {
		return nil, err
	}
	if strings.TrimSpace(hardwareID) == "" {
		return nil, nil
	}

	deviceID := fmt.Sprintf("%s:%s", hardwareID, model)
	now := time.Now()
	err = database.DB.QueryRow(`
		INSERT INTO devices (
			id, farm_id, coop_id, device_id, name, type, model,
			is_main_controller, firmware_version, hardware_id,
			is_active, is_online, created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,false,$8,$9,true,true,$10,$10)
		ON CONFLICT (device_id) DO UPDATE SET
			name = EXCLUDED.name, type = EXCLUDED.type, model = EXCLUDED.model,
			farm_id = EXCLUDED.farm_id, coop_id = EXCLUDED.coop_id,
			hardware_id = EXCLUDED.hardware_id, is_active = true, is_online = true, updated_at = EXCLUDED.updated_at
		RETURNING id
	`, uuid.New(), farmID, coopID, deviceID, name, "sensor", model, "gateway", hardwareID, now).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func insertReading(deviceID uuid.UUID, sensorType string, value float64, unit string, ts time.Time) error {
	_, err := database.DB.Exec(`
		INSERT INTO device_readings (id, device_id, sensor_type, value, unit, quality, timestamp)
		VALUES ($1, $2, $3, $4, $5, 'good', $6)
	`, uuid.New(), deviceID, sensorType, value, unit, ts)
	return err
}

func (s *TelemetryService) GetTemperatureTimeline(userID, farmID, coopID uuid.UUID, days int) (schemas.TemperatureTimelineResponse, error) {
	if days < 1 {
		days = 7
	}
	if err := s.farmService.CheckAccess(userID, farmID, "viewer"); err != nil {
		return schemas.TemperatureTimelineResponse{SensorFound: false}, err
	}
	coop, err := s.coopService.GetCoop(userID, farmID, coopID)
	if err != nil {
		return schemas.TemperatureTimelineResponse{SensorFound: false}, err
	}

	loc, _ := time.LoadLocation("Asia/Phnom_Penh")
	now := time.Now().In(loc)
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	start := todayStart.AddDate(0, 0, -(days-1))

	// Latest readings
	var currentTemp *float64
	var currentHumidity *float64
	var tempVal sql.NullFloat64
	var humidVal sql.NullFloat64
	_ = database.DB.QueryRow(`
		SELECT dr.value FROM device_readings dr
		JOIN devices d ON dr.device_id = d.id
		WHERE d.coop_id = $1 AND dr.sensor_type = 'temperature'
		ORDER BY dr.timestamp DESC LIMIT 1
	`, coopID).Scan(&tempVal)
	if tempVal.Valid {
		v := tempVal.Float64
		currentTemp = &v
	}
	_ = database.DB.QueryRow(`
		SELECT dr.value FROM device_readings dr
		JOIN devices d ON dr.device_id = d.id
		WHERE d.coop_id = $1 AND dr.sensor_type = 'humidity'
		ORDER BY dr.timestamp DESC LIMIT 1
	`, coopID).Scan(&humidVal)
	if humidVal.Valid {
		v := humidVal.Float64
		currentHumidity = &v
	}

	rows, err := database.DB.Query(`
		SELECT dr.value, dr.timestamp
		FROM device_readings dr
		JOIN devices d ON dr.device_id = d.id
		WHERE d.coop_id = $1 AND dr.sensor_type = 'temperature' AND dr.timestamp >= $2
		ORDER BY dr.timestamp ASC
	`, coopID, start)
	if err != nil {
		return schemas.TemperatureTimelineResponse{SensorFound: false}, err
	}
	defer rows.Close()

	readings := make([]tempReading, 0)
	for rows.Next() {
		var r tempReading
		if err := rows.Scan(&r.Temp, &r.Time); err == nil {
			readings = append(readings, r)
		}
	}
	if len(readings) == 0 {
		return schemas.TemperatureTimelineResponse{
			SensorFound: false,
			CoopName:    coop.Name,
			History:     []schemas.DaySummary{},
		}, nil
	}

	// Bucket by day
	dayBuckets := map[string][]tempReading{}
	for _, r := range readings {
		t := r.Time.In(loc)
		dayKey := t.Format("2006-01-02")
		dayBuckets[dayKey] = append(dayBuckets[dayKey], tempReading{Temp: r.Temp, Time: t})
	}

	// Today hourly
	hourMap := map[int][]float64{}
	todayKey := todayStart.Format("2006-01-02")
	if todays, ok := dayBuckets[todayKey]; ok {
		for _, r := range todays {
			hour := r.Time.Hour()
			hourMap[hour] = append(hourMap[hour], r.Temp)
		}
	}
	hours := make([]int, 0, len(hourMap))
	for h := range hourMap {
		hours = append(hours, h)
	}
	sort.Ints(hours)

	hourly := make([]schemas.HourlyPoint, 0, len(hours))
	for _, h := range hours {
		vals := hourMap[h]
		sum := 0.0
		for _, v := range vals { sum += v }
		avg := sum / float64(len(vals))
		hourly = append(hourly, schemas.HourlyPoint{
			Hour: fmt.Sprintf("%02d:00", h),
			Temp: round1(avg),
		})
	}

	todaySummary := summarizeDay(dayBuckets[todayKey], hourly, "")

	// History (previous days)
	history := make([]schemas.DaySummary, 0)
	for i := 1; i < days; i++ {
		dt := todayStart.AddDate(0, 0, -i)
		key := dt.Format("2006-01-02")
		if items, ok := dayBuckets[key]; ok {
			label := dt.Format("Jan 02")
			if i == 1 {
				label = "Yesterday"
			}
			history = append(history, summarizeDay(items, nil, label))
		}
	}

	bgHint := "neutral"
	if currentTemp != nil {
		bgHint = tempToBg(*currentTemp)
	}

	return schemas.TemperatureTimelineResponse{
		SensorFound:     true,
		CoopName:        coop.Name,
		CurrentTemp:     currentTemp,
		CurrentHumidity: currentHumidity,
		BgHint:          bgHint,
		Today:           &todaySummary,
		History:         history,
	}, nil
}

func summarizeDay(items []tempReading, hourly []schemas.HourlyPoint, label string) schemas.DaySummary {
	if len(items) == 0 {
		return schemas.DaySummary{Label: label, Hourly: hourly}
	}
	min := items[0]
	max := items[0]
	for _, r := range items[1:] {
		if r.Temp < min.Temp { min = r }
		if r.Temp > max.Temp { max = r }
	}
	high := &schemas.TempPoint{Temp: round1(max.Temp), Time: max.Time.Format("15:04")}
	low := &schemas.TempPoint{Temp: round1(min.Temp), Time: min.Time.Format("15:04")}
	return schemas.DaySummary{
		Label:  label,
		High:   high,
		Low:    low,
		Hourly: hourly,
	}
}

func round1(v float64) float64 {
	return float64(int(v*10+0.5)) / 10
}

func tempToBg(temp float64) string {
	switch {
	case temp >= 35:
		return "scorching"
	case temp >= 32:
		return "hot"
	case temp >= 28:
		return "warm"
	case temp >= 22:
		return "neutral"
	case temp >= 18:
		return "cool"
	default:
		return "cold"
	}
}
