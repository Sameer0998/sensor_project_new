package domain

import (
	"time"
)

// SensorData represents a single sensor reading
type SensorData struct {
	SensorValue float64 `json:"sensor_value"`
	SensorType  string  `json:"sensor_type"`
	ID1         string  `json:"id1"`
	ID2         int     `json:"id2"`
	Timestamp   int64   `json:"timestamp"`
}

// SensorConfig holds the configuration for sensor data generation
type SensorConfig struct {
	SensorType     string        `json:"sensor_type"`
	MinValue       float64       `json:"min_value"`
	MaxValue       float64       `json:"max_value"`
	GenerationRate time.Duration `json:"generation_rate_ms"`
}

// SensorGenerator defines the interface for generating sensor data
type SensorGenerator interface {
	GenerateSensorData() *SensorData
	SetGenerationRate(rate time.Duration)
	GetGenerationRate() time.Duration
	GetSensorType() string
}

// SensorSender defines the interface for sending sensor data
type SensorSender interface {
	SendSensorData(data *SensorData) error
	Close() error
}
