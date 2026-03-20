package schemas

import "time"

type DeviceReportItem struct {
	Type            string  `json:"type"`
	Model           string  `json:"model"`
	Name            string  `json:"name"`
	Active          *bool   `json:"active,omitempty"`
	FirmwareVersion *string `json:"firmware_version,omitempty"`
}

type DeviceReportRequest struct {
	HardwareID string             `json:"hardware_id"`
	Devices    []DeviceReportItem `json:"devices"`
}

type TelemetrySensors struct {
	TemperatureC *float64 `json:"temperature_c,omitempty"`
	HumidityPct  *float64 `json:"humidity_pct,omitempty"`
	WaterLevel   *float64 `json:"water_level_raw,omitempty"`
}

type TelemetryRequest struct {
	HardwareID string           `json:"hardware_id,omitempty"`
	Timestamp  *time.Time       `json:"timestamp,omitempty"`
	Sensors    TelemetrySensors `json:"sensors"`
}

type TempPoint struct {
	Temp float64 `json:"temp"`
	Time string  `json:"time"`
}

type HourlyPoint struct {
	Hour string  `json:"hour"`
	Temp float64 `json:"temp"`
}

type DaySummary struct {
	Label  string        `json:"label,omitempty"`
	High   *TempPoint    `json:"high,omitempty"`
	Low    *TempPoint    `json:"low,omitempty"`
	Hourly []HourlyPoint `json:"hourly,omitempty"`
}

type TemperatureTimelineResponse struct {
	SensorFound     bool        `json:"sensor_found"`
	CoopName        string      `json:"coop_name"`
	CurrentTemp     *float64    `json:"current_temp,omitempty"`
	CurrentHumidity *float64    `json:"current_humidity,omitempty"`
	BgHint          string      `json:"bg_hint"`
	Today           *DaySummary `json:"today,omitempty"`
	History         []DaySummary `json:"history"`
}
