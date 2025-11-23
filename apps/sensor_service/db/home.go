package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"smart-home-service/models"
	"github.com/jackc/pgx/v5"
)

// GetHomes получает все дома из базы данных
func (db *DB) GetHomes(ctx context.Context) ([]models.Home, error) {
	query := `
		SELECT home_id, user_id, name, city, street, num, created_at
		FROM homes
		ORDER BY name
	`

	rows, err := db.Pool.Query(ctx, query)
	if err != nil {
	    log.Printf("ERROR: querying homes: %w", err)
		return nil, fmt.Errorf("error querying homes: %w", err)
	}
	defer rows.Close()

	var homes []models.Home
	for rows.Next() {
		var h models.Home
		err := rows.Scan(
		    &h.HomeID, &h.UserID, &h.Name, &h.City, &h.Street, &h.Num, &h.CreatedAt,)
		if err != nil {
		    log.Printf("ERROR: scanning home row: %w", err)
			return nil, fmt.Errorf("error scanning home row: %w", err)
		}
		homes = append(homes, h)
	}

	if err := rows.Err(); err != nil {
	    log.Printf("ERROR: error collecting home rows: %w", err)
		return nil, fmt.Errorf("error iterating home rows: %w", err)
	}

	return homes, nil
}

// GetHomeByID получает дом по его ID
func (db *DB) GetHomeByID(ctx context.Context, id int) (models.Home, error) {
	query := `
		SELECT home_id, user_id, name, city, street, num, created_at
		FROM homes
		WHERE home_id = $1
	`
	var h models.Home
	err := db.Pool.QueryRow(ctx, query, id).Scan(
		&h.HomeID, &h.UserID, &h.Name, &h.City, &h.Street, &h.Num, &h.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
		    log.Printf("ERROR: home with id %s not found", id)
			return models.Home{}, fmt.Errorf("home with id %s not found", id)
		}
        log.Printf("ERROR: getting home by ID: %w", err)
		return models.Home{}, fmt.Errorf("error getting home by ID: %w", err)
	}
	return h, nil
}

// CreateHome создает новый дом в базе данных
func (db *DB) CreateHome(ctx context.Context, h models.HomeCreate) (models.Home, error) {
	query := `
		INSERT INTO homes (user_id, name, city, street, num)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING home_id, user_id, name, city, street, num, created_at
	`
	var newHome models.Home
	err := db.Pool.QueryRow(ctx, query, h.UserID, h.Name, h.City, h.Street, h.Num).Scan(
		&newHome.HomeID, &newHome.UserID, &newHome.Name, &newHome.City, &newHome.Street, &newHome.Num, &newHome.CreatedAt,
	)
	if err != nil {
	    log.Printf("ERROR: error creating home: %w", err)
		return models.Home{}, fmt.Errorf("error creating home: %w", err)
	}
	return newHome, nil
}

// UpdateHome обновляет существующий дом
func (db *DB) UpdateHome(ctx context.Context, id int, h models.HomeUpdate) (models.Home, error) {
	var setClauses []string
	args := []interface{}{}
	argCount := 1

	if h.Name != "" {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", argCount))
		args = append(args, h.Name)
		argCount++
	}
	if h.City != "" {
		setClauses = append(setClauses, fmt.Sprintf("city = $%d", argCount))
		args = append(args, h.City)
		argCount++
	}
	if h.Street != "" {
		setClauses = append(setClauses, fmt.Sprintf("street = $%d", argCount))
		args = append(args, h.Street)
		argCount++
	}
	if h.Num > 0 {
		setClauses = append(setClauses, fmt.Sprintf("num = $%d", argCount))
		args = append(args, h.Num)
		argCount++
	}

	if len(setClauses) == 0 {
		return db.GetHomeByID(ctx, id)
	}

	query := `UPDATE homes SET ` + strings.Join(setClauses, ", ") + `
		WHERE home_id = $` + fmt.Sprintf("%d", argCount) + `
		RETURNING home_id, user_id, name, city, street, num, created_at`
	args = append(args, id)

	var updatedHome models.Home
	err := db.Pool.QueryRow(ctx, query, args...).Scan(
		&updatedHome.HomeID, &updatedHome.UserID, &updatedHome.Name, &updatedHome.City, &updatedHome.Street, &updatedHome.Num, &updatedHome.CreatedAt,
	)
	if err != nil {
		// Добавим проверку на 'not found', которая может прийти из GetHomeByID, если запрос ничего не обновил
		if errors.Is(err, pgx.ErrNoRows) {
		    log.Printf("ERROR: home with id %s not found or no fields to update", id)
			return models.Home{}, fmt.Errorf("home with id %s not found or no fields to update", id)
		}
        log.Printf("ERROR: updating home: %w", err)
		return models.Home{}, fmt.Errorf("error updating home: %w", err)
	}
	return updatedHome, nil
}

// DeleteHome удаляет дом по его ID
func (db *DB) DeleteHome(ctx context.Context, id int) error {
	query := "DELETE FROM homes WHERE home_id = $1"
	result, err := db.Pool.Exec(ctx, query, id)
	if err != nil {
	    log.Printf("ERROR: deleting home: %w", err)
		return fmt.Errorf("error deleting home: %w", err)
	}
	if result.RowsAffected() == 0 {
	    log.Printf("ERROR: home with id %s not found", id)
		return fmt.Errorf("home with id %s not found", id)
	}
	return nil
}
