package http

import (
	"net/http"
	"strconv"
	"time"

	"sensor_project/microservice-b/internal/domain"

	"github.com/labstack/echo/v4"
)

// Handler handles HTTP requests for the sensor data API
type Handler struct {
	sensorUseCase domain.SensorDataUseCase
}

// NewHandler creates a new HTTP handler
func NewHandler(sensorUseCase domain.SensorDataUseCase) *Handler {
	return &Handler{
		sensorUseCase: sensorUseCase,
	}
}

// SetupRoutes configures the HTTP routes
func (h *Handler) SetupRoutes(e *echo.Echo) {
	// Health check
	e.GET("/health", h.HealthCheck)

	// API routes
	api := e.Group("/api")

	// Sensor data routes
	api.GET("/sensor-data/:id", h.GetSensorDataByID)
	api.GET("/sensor-data", h.GetSensorDataByFilter)
	api.PUT("/sensor-data/:id", h.UpdateSensorData)
	api.DELETE("/sensor-data/:id", h.DeleteSensorData)
	api.DELETE("/sensor-data", h.DeleteSensorDataByFilter)
}

// HealthCheck handles health check requests
// @Summary Health check endpoint
// @Description Check if the service is running
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (h *Handler) HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
	})
}

// GetSensorDataByID retrieves a sensor data record by ID
// @Summary Get sensor data by ID
// @Description Retrieve a single sensor data record by its ID
// @Tags sensor-data
// @Produce json
// @Param id path int true "Sensor Data ID"
// @Success 200 {object} SensorDataResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/sensor-data/{id} [get]
func (h *Handler) GetSensorDataByID(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid ID format"})
	}

	data, err := h.sensorUseCase.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to retrieve sensor data"})
	}

	if data == nil {
		return c.JSON(http.StatusNotFound, ErrorResponse{Error: "Sensor data not found"})
	}

	return c.JSON(http.StatusOK, SensorDataResponse{
		ID:          data.ID,
		SensorValue: data.SensorValue,
		SensorType:  data.SensorType,
		ID1:         data.ID1,
		ID2:         data.ID2,
	})
}

// GetSensorDataByFilter retrieves sensor data records based on filter criteria
// @Summary Get sensor data by filter
// @Description Retrieve sensor data records based on filter criteria with pagination
// @Tags sensor-data
// @Produce json
// @Param id1 query string false "ID1 filter"
// @Param id2 query int false "ID2 filter"
// @Param sensor_type query string false "Sensor type filter"
// @Param start_time query string false "Start time filter (RFC3339 format)"
// @Param end_time query string false "End time filter (RFC3339 format)"
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10, max: 100)"
// @Success 200 {object} PaginatedResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/sensor-data [get]
func (h *Handler) GetSensorDataByFilter(c echo.Context) error {
	filter, err := parseFilterFromQuery(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}

	data, total, err := h.sensorUseCase.GetByFilter(filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to retrieve sensor data"})
	}

	// Convert domain models to response models
	var results []SensorDataResponse
	for _, item := range data {
		results = append(results, SensorDataResponse{
			ID:          item.ID,
			SensorValue: item.SensorValue,
			SensorType:  item.SensorType,
			ID1:         item.ID1,
			ID2:         item.ID2,
		})
	}

	return c.JSON(http.StatusOK, PaginatedResponse{
		Data:       results,
		Total:      total,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		TotalPages: (total + filter.PageSize - 1) / filter.PageSize,
	})
}

// UpdateSensorData updates a sensor data record
// @Summary Update sensor data
// @Description Update a sensor data record by ID
// @Tags sensor-data
// @Accept json
// @Produce json
// @Param id path int true "Sensor Data ID"
// @Param data body UpdateSensorDataRequest true "Sensor data update"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/sensor-data/{id} [put]
func (h *Handler) UpdateSensorData(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid ID format"})
	}

	// Check if the record exists
	data, err := h.sensorUseCase.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to retrieve sensor data"})
	}

	if data == nil {
		return c.JSON(http.StatusNotFound, ErrorResponse{Error: "Sensor data not found"})
	}

	// Parse the update request
	req := new(UpdateSensorDataRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request format"})
	}

	// Create the update model
	update := &domain.SensorDataUpdate{
		SensorValue: &req.SensorValue,
	}

	// Update the record
	err = h.sensorUseCase.Update(id, update)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to update sensor data"})
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Message: "Sensor data updated successfully",
	})
}

// DeleteSensorData deletes a sensor data record by ID
// @Summary Delete sensor data
// @Description Delete a sensor data record by ID
// @Tags sensor-data
// @Produce json
// @Param id path int true "Sensor Data ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/sensor-data/{id} [delete]
func (h *Handler) DeleteSensorData(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid ID format"})
	}

	err = h.sensorUseCase.Delete(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to delete sensor data"})
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Message: "Sensor data deleted successfully",
	})
}

// DeleteSensorDataByFilter deletes sensor data records based on filter criteria
// @Summary Delete sensor data by filter
// @Description Delete sensor data records based on filter criteria
// @Tags sensor-data
// @Accept json
// @Produce json
// @Param filter body FilterRequest true "Filter criteria"
// @Success 200 {object} DeleteResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security ApiKeyAuth
// @Router /api/sensor-data [delete]
func (h *Handler) DeleteSensorDataByFilter(c echo.Context) error {
	req := new(FilterRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request format"})
	}

	filter := &domain.SensorDataFilter{
		ID1:        req.ID1,
		ID2:        req.ID2,
		SensorType: req.SensorType,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
	}

	count, err := h.sensorUseCase.DeleteByFilter(filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to delete sensor data"})
	}

	return c.JSON(http.StatusOK, DeleteResponse{
		Success:     true,
		Message:     "Sensor data deleted successfully",
		DeletedRows: count,
	})
}

// parseFilterFromQuery parses filter parameters from the query string
func parseFilterFromQuery(c echo.Context) (*domain.SensorDataFilter, error) {
	filter := &domain.SensorDataFilter{}

	// Parse ID1
	if id1 := c.QueryParam("id1"); id1 != "" {
		filter.ID1 = &id1
	}

	// Parse ID2
	if id2Str := c.QueryParam("id2"); id2Str != "" {
		id2, err := strconv.Atoi(id2Str)
		if err != nil {
			return nil, err
		}
		filter.ID2 = &id2
	}

	// Parse sensor type
	if sensorType := c.QueryParam("sensor_type"); sensorType != "" {
		filter.SensorType = &sensorType
	}

	// Parse start time
	if startTimeStr := c.QueryParam("start_time"); startTimeStr != "" {
		startTime, err := time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			return nil, err
		}
		filter.StartTime = &startTime
	}

	// Parse end time
	if endTimeStr := c.QueryParam("end_time"); endTimeStr != "" {
		endTime, err := time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			return nil, err
		}
		filter.EndTime = &endTime
	}

	// Parse pagination
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 {
		page = 1
	}
	filter.Page = page

	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	filter.PageSize = pageSize

	return filter, nil
}
