package usecase

import (
	"sensor_project/microservice-b/internal/domain"
)

// SensorDataUseCase implements the domain.SensorDataUseCase interface
type SensorDataUseCase struct {
	repo domain.SensorDataRepository
}

// NewSensorDataUseCase creates a new sensor data use case
func NewSensorDataUseCase(repo domain.SensorDataRepository) domain.SensorDataUseCase {
	return &SensorDataUseCase{
		repo: repo,
	}
}

// Store saves a sensor data record
func (uc *SensorDataUseCase) Store(data *domain.SensorData) error {
	return uc.repo.Store(data)
}

// GetByID retrieves a sensor data record by ID
func (uc *SensorDataUseCase) GetByID(id int64) (*domain.SensorData, error) {
	return uc.repo.GetByID(id)
}

// GetByFilter retrieves sensor data records based on filter criteria
func (uc *SensorDataUseCase) GetByFilter(filter *domain.SensorDataFilter) ([]*domain.SensorData, int, error) {
	// Set default pagination values if not provided
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 10
	}

	if filter.PageSize > 100 {
		filter.PageSize = 100
	}

	return uc.repo.GetByFilter(filter)
}

// Update updates a sensor data record
func (uc *SensorDataUseCase) Update(id int64, update *domain.SensorDataUpdate) error {
	return uc.repo.Update(id, update)
}

// Delete removes a sensor data record by ID
func (uc *SensorDataUseCase) Delete(id int64) error {
	return uc.repo.Delete(id)
}

// DeleteByFilter removes sensor data records based on filter criteria
func (uc *SensorDataUseCase) DeleteByFilter(filter *domain.SensorDataFilter) (int, error) {
	return uc.repo.DeleteByFilter(filter)
}
