package config

import (
	"time"
)

// Config holds the application configuration
type Config struct {
	ServerPort     string
	GRPCPort       string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	JWTSecret      string
	TokenExpiry    time.Duration
	MaxConnections int
}

// LoadConfig loads the configuration from environment variables
func LoadConfig() *Config {
	// Default values
	config := &Config{
		ServerPort:     "8080",
		GRPCPort:       "50051",
		DBHost:         "mysql",
		DBPort:         "3306",
		DBUser:         "root",
		DBPassword:     "Sameer@0998$B@",
		DBName:         "sensor_data",
		JWTSecret:      "default_jwt_secret",
		TokenExpiry:    24 * time.Hour,
		MaxConnections: 10,
	}

	return config
}

// GetDSN returns the database connection string
func (c *Config) GetDSN() string {
	return c.DBUser + ":" + c.DBPassword + "@tcp(" + c.DBHost + ":" + c.DBPort + ")/" + c.DBName + "?parseTime=true"
}
