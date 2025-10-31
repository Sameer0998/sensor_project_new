package http

import (
	"net/http"
	"time"

	"sensor_project/microservice-a/internal/domain"

	"github.com/labstack/echo/v4"
)

type FrequencyRequest struct {
	IntervalMs int `json:"interval_ms"`
}

type Handler struct {
	generator domain.SensorGenerator
}

func NewHandler(generator domain.SensorGenerator) *Handler {
	return &Handler{
		generator: generator,
	}
}

func (h *Handler) SetupRoutes(e *echo.Echo) {
	e.GET("/health", h.HealthCheck)
	e.GET("/config", h.GetConfig)
	e.POST("/config/frequency", h.SetFrequency)
}

func (h *Handler) HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
		"type":   h.generator.GetSensorType(),
	})
}

func (h *Handler) GetConfig(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"sensor_type":        h.generator.GetSensorType(),
		"generation_rate_ms": h.generator.GetGenerationRate().Milliseconds(),
	})
}

func (h *Handler) SetFrequency(c echo.Context) error {
	req := new(FrequencyRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	if req.IntervalMs < 100 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Interval must be at least 100ms",
		})
	}

	h.generator.SetGenerationRate(time.Duration(req.IntervalMs) * time.Millisecond)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success":            true,
		"message":            "Frequency updated successfully",
		"generation_rate_ms": req.IntervalMs,
	})
}
