package http

import (
	"time"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse represents a generic success response
type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// DeleteResponse represents a response for delete operations
type DeleteResponse struct {
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	DeletedRows int    `json:"deleted_rows"`
}

// UserResponse represents user information in responses
type UserResponse struct {
	ID       int      `json:"id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles"`
}

// SensorDataResponse represents a sensor data record in responses
type SensorDataResponse struct {
	ID          int64   `json:"id"`
	SensorValue float64 `json:"sensor_value"`
	SensorType  string  `json:"sensor_type"`
	ID1         string  `json:"id1"`
	ID2         int     `json:"id2"`
	CreatedAt   string  `json:"created_at"`
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Data       []SensorDataResponse `json:"data"`
	Total      int                  `json:"total"`
	Page       int                  `json:"page"`
	PageSize   int                  `json:"page_size"`
	TotalPages int                  `json:"total_pages"`
}

// UpdateSensorDataRequest represents a request to update sensor data
type UpdateSensorDataRequest struct {
	SensorValue float64 `json:"sensor_value"`
}

// FilterRequest represents filter criteria for querying or deleting sensor data
type FilterRequest struct {
	ID1        *string    `json:"id1"`
	ID2        *int       `json:"id2"`
	SensorType *string    `json:"sensor_type"`
	StartTime  *time.Time `json:"start_time"`
	EndTime    *time.Time `json:"end_time"`
}
