package config

import (
	"sensor_project/microservice-a/internal/domain"
	"time"
)

// Config holds the application configuration
type Config struct {
	ServerPort     string
	GRPCServerAddr string
	SensorConfig   domain.SensorConfig
}

// LoadConfig loads the configuration from environment variables
func LoadConfig() *Config {
	// Default values
	config := &Config{
		ServerPort:     "8090",
		GRPCServerAddr: "127.0.0.1:50051",
		SensorConfig: domain.SensorConfig{
			SensorType:     "temperature",
			MinValue:       0.0,
			MaxValue:       100.0,
			GenerationRate: 1000 * time.Millisecond,
		},
	}

	return config
}
