package db

import (
	"context"
	"fmt"
	"log"
	"smart-home-service/models"
)

// CreateSensorLink сохраняет связь между домом и внешним датчиком
func (db *DB) CreateSensorLink(ctx context.Context, s models.Sensor) error {
	query := `
		INSERT INTO sensors (home_id, service_id)
		VALUES ($1, $2)
	`
	_, err := db.Pool.Exec(ctx, query, s.HomeID, s.ServiceID)
	if err != nil {
		log.Printf("ERROR: error linking sensor: %w", err)
		return fmt.Errorf("error linking sensor: %w", err)
	}
	return nil
}

// GetSensorIDsByHomeID возвращает список service_id для указанного дома
func (db *DB) GetSensorIDsByHomeID(ctx context.Context, homeID int) ([]int, error) {
	query := `SELECT service_id FROM sensors WHERE home_id = $1`

	rows, err := db.Pool.Query(ctx, query, homeID)
	if err != nil {
		log.Printf("ERROR: querying sensor IDs: %v", err)
		return nil, fmt.Errorf("error querying sensor IDs: %w", err)
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			log.Printf("ERROR: scanning sensor ID: %v", err)
			continue
		}
		ids = append(ids, id)
	}

	return ids, nil
}