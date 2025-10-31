package usecase

import (
	"math/rand"
	"sync"
	"time"

	"sensor_project/microservice-a/internal/domain"
)

// DefaultSensorGenerator implements the SensorGenerator interface
type DefaultSensorGenerator struct {
	sensorType     string
	minValue       float64
	maxValue       float64
	generationRate time.Duration
	mu             sync.RWMutex
}

// NewSensorGenerator creates a new sensor generator with the given configuration
func NewSensorGenerator(config domain.SensorConfig) domain.SensorGenerator {
	return &DefaultSensorGenerator{
		sensorType:     config.SensorType,
		minValue:       config.MinValue,
		maxValue:       config.MaxValue,
		generationRate: config.GenerationRate,
	}
}

// GenerateSensorData generates a new sensor reading
func (g *DefaultSensorGenerator) GenerateSensorData() *domain.SensorData {
	g.mu.RLock()
	defer g.mu.RUnlock()

	// Generate random sensor value between min and max
	value := g.minValue + rand.Float64()*(g.maxValue-g.minValue)

	// Generate random ID1 (capital letters)
	id1 := string(rune('A' + rand.Intn(26)))

	// Generate random ID2 (integer)
	id2 := rand.Intn(100)

	return &domain.SensorData{
		SensorValue: value,
		SensorType:  g.sensorType,
		ID1:         id1,
		ID2:         id2,
		Timestamp:   time.Now().UnixMilli(),
	}
}

// SetGenerationRate updates the data generation rate
func (g *DefaultSensorGenerator) SetGenerationRate(rate time.Duration) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.generationRate = rate
}

// GetGenerationRate returns the current data generation rate
func (g *DefaultSensorGenerator) GetGenerationRate() time.Duration {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.generationRate
}

// GetSensorType returns the sensor type
func (g *DefaultSensorGenerator) GetSensorType() string {
	return g.sensorType
}
