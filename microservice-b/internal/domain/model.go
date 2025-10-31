package domain

import (
	"time"
)

// SensorData represents a single sensor reading
type SensorData struct {
	ID          int64   `json:"id"`
	SensorValue float64 `json:"sensor_value"`
	SensorType  string  `json:"sensor_type"`
	ID1         string  `json:"id1"`
	ID2         int     `json:"id2"`
	CreatedAt   int64   `json:"created_at"`
}

// SensorDataFilter represents filter criteria for querying sensor data
type SensorDataFilter struct {
	ID1        *string    `json:"id1"`
	ID2        *int       `json:"id2"`
	SensorType *string    `json:"sensor_type"`
	StartTime  *time.Time `json:"start_time"`
	EndTime    *time.Time `json:"end_time"`
	Page       int        `json:"page"`
	PageSize   int        `json:"page_size"`
}

// SensorDataUpdate represents fields that can be updated
type SensorDataUpdate struct {
	SensorValue *float64 `json:"sensor_value"`
}

// User represents a user in the system
type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	Email        string    `json:"email"`
	Roles        []string  `json:"roles"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// APIKey represents an API key for authentication
type APIKey struct {
	ID        int       `json:"id"`
	Key       string    `json:"key"`
	UserID    int       `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// SensorDataRepository defines the interface for sensor data storage
type SensorDataRepository interface {
	Store(data *SensorData) error
	GetByID(id int64) (*SensorData, error)
	GetByFilter(filter *SensorDataFilter) ([]*SensorData, int, error)
	Update(id int64, update *SensorDataUpdate) error
	Delete(id int64) error
	DeleteByFilter(filter *SensorDataFilter) (int, error)
}

// UserRepository defines the interface for user storage
type UserRepository interface {
	GetByID(id int) (*User, error)
	GetByUsername(username string) (*User, error)
	GetByAPIKey(key string) (*User, error)
	CreateAPIKey(userID int, expiresAt time.Time) (*APIKey, error)
}

// SensorDataUseCase defines the interface for sensor data business logic
type SensorDataUseCase interface {
	Store(data *SensorData) error
	GetByID(id int64) (*SensorData, error)
	GetByFilter(filter *SensorDataFilter) ([]*SensorData, int, error)
	Update(id int64, update *SensorDataUpdate) error
	Delete(id int64) error
	DeleteByFilter(filter *SensorDataFilter) (int, error)
}

// AuthUseCase defines the interface for authentication business logic
type AuthUseCase interface {
	Authenticate(username, password string) (*User, error)
	ValidateAPIKey(key string) (*User, error)
	GenerateToken(user *User) (string, error)
	ValidateToken(token string) (*User, error)
}
