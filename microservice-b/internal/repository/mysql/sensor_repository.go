package mysql

import (
	"database/sql"
	"fmt"
	"strings"

	"sensor_project/microservice-b/internal/domain"

	_ "github.com/go-sql-driver/mysql"
)

// MySQLSensorRepository implements the SensorDataRepository interface
type MySQLSensorRepository struct {
	db *sql.DB
}

// NewMySQLSensorRepository creates a new MySQL repository
func NewMySQLSensorRepository(db *sql.DB) domain.SensorDataRepository {
	return &MySQLSensorRepository{
		db: db,
	}
}

// Store saves a sensor data record to the database
func (r *MySQLSensorRepository) Store(data *domain.SensorData) error {
	query := `
		INSERT INTO sensor_data (sensor_value, sensor_type_id, id1, id2, created_at)
		SELECT ?, id, ?, ?, ?
		FROM sensor_types
		WHERE name = ?
	`

	_, err := r.db.Exec(
		query,
		data.SensorValue,
		data.ID1,
		data.ID2,
		data.CreatedAt,
		data.SensorType,
	)

	return err
}

// GetByID retrieves a sensor data record by ID
func (r *MySQLSensorRepository) GetByID(id int64) (*domain.SensorData, error) {
	fmt.Println("reached here ??")
	query := `
		SELECT sd.id, sd.sensor_value, st.name, sd.id1, sd.id2, sd.created_at
		FROM sensor_data sd
		JOIN sensor_types st ON sd.sensor_type_id = st.id
		WHERE sd.id = ?
	`

	var data domain.SensorData
	err := r.db.QueryRow(query, id).Scan(
		&data.ID,
		&data.SensorValue,
		&data.SensorType,
		&data.ID1,
		&data.ID2,
		&data.CreatedAt,
	)

	fmt.Println("errror : ", err)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &data, nil
}

// GetByFilter retrieves sensor data records based on filter criteria
func (r *MySQLSensorRepository) GetByFilter(filter *domain.SensorDataFilter) ([]*domain.SensorData, int, error) {
	// Build the WHERE clause based on filter criteria
	whereClause, args := r.buildWhereClause(filter)

	// Count total records matching the filter
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM sensor_data sd
		JOIN sensor_types st ON sd.sensor_type_id = st.id
		%s
	`, whereClause)

	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Build the main query with pagination
	query := fmt.Sprintf(`
		SELECT sd.id, sd.sensor_value, st.name, sd.id1, sd.id2, sd.created_at
		FROM sensor_data sd
		JOIN sensor_types st ON sd.sensor_type_id = st.id
		%s
		ORDER BY sd.created_at DESC
		LIMIT ? OFFSET ?
	`, whereClause)

	// Add pagination parameters
	offset := (filter.Page - 1) * filter.PageSize
	args = append(args, filter.PageSize, offset)

	// Execute the query
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// Parse the results
	var results []*domain.SensorData
	for rows.Next() {
		var data domain.SensorData
		err := rows.Scan(
			&data.ID,
			&data.SensorValue,
			&data.SensorType,
			&data.ID1,
			&data.ID2,
			&data.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		results = append(results, &data)
	}

	return results, total, nil
}

// Update updates a sensor data record
func (r *MySQLSensorRepository) Update(id int64, update *domain.SensorDataUpdate) error {
	if update.SensorValue == nil {
		return nil // Nothing to update
	}

	query := `
		UPDATE sensor_data
		SET sensor_value = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(query, *update.SensorValue, id)
	return err
}

// Delete removes a sensor data record by ID
func (r *MySQLSensorRepository) Delete(id int64) error {
	query := `DELETE FROM sensor_data WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

// DeleteByFilter removes sensor data records based on filter criteria
func (r *MySQLSensorRepository) DeleteByFilter(filter *domain.SensorDataFilter) (int, error) {
	// Build the WHERE clause based on filter criteria
	whereClause, args := r.buildWhereClause(filter)

	// Build the delete query
	query := fmt.Sprintf(`
		DELETE sd FROM sensor_data sd
		JOIN sensor_types st ON sd.sensor_type_id = st.id
		%s
	`, whereClause)

	// Execute the query
	result, err := r.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}

	// Get the number of affected rows
	affected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(affected), nil
}

// buildWhereClause constructs a WHERE clause based on filter criteria
func (r *MySQLSensorRepository) buildWhereClause(filter *domain.SensorDataFilter) (string, []interface{}) {
	var conditions []string
	var args []interface{}

	if filter.ID1 != nil {
		conditions = append(conditions, "sd.id1 = ?")
		args = append(args, *filter.ID1)
	}

	if filter.ID2 != nil {
		conditions = append(conditions, "sd.id2 = ?")
		args = append(args, *filter.ID2)
	}

	if filter.SensorType != nil {
		conditions = append(conditions, "st.name = ?")
		args = append(args, *filter.SensorType)
	}

	if filter.StartTime != nil {
		conditions = append(conditions, "sd.created_at >= ?")
		args = append(args, *filter.StartTime)
	}

	if filter.EndTime != nil {
		conditions = append(conditions, "sd.created_at <= ?")
		args = append(args, *filter.EndTime)
	}

	if len(conditions) > 0 {
		return "WHERE " + strings.Join(conditions, " AND "), args
	}

	return "", args
}
