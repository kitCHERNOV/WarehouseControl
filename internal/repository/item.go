package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"wbtechschool-L3/WarehouseControl/internal/models"
)

// CreateItem creates a new item in the database
func (r *Repository) CreateItem(ctx context.Context, item *models.Item) error {
	query := `
		INSERT INTO items (name, quantity, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	now := time.Now()
	item.CreatedAt = now
	item.UpdatedAt = now

	err := r.db.QueryRowContext(ctx, query, item.Name, item.Quantity, now, now).Scan(&item.ID)
	if err != nil {
		return fmt.Errorf("failed to create item: %w", err)
	}

	return nil
}

// GetItem retrieves an item by ID
func (r *Repository) GetItem(ctx context.Context, id int) (*models.Item, error) {
	query := `
		SELECT id, name, quantity, created_at, updated_at
		FROM items
		WHERE id = $1
	`

	var item models.Item
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&item.ID,
		&item.Name,
		&item.Quantity,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("item not found: %w", err)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get item: %w", err)
	}

	return &item, nil
}

// GetAllItems retrieves all items from the database
func (r *Repository) GetAllItems(ctx context.Context) ([]*models.Item, error) {
	query := `
		SELECT id, name, quantity, created_at, updated_at
		FROM items
		ORDER BY id
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get items: %w", err)
	}
	defer rows.Close()

	var items []*models.Item
	for rows.Next() {
		var item models.Item
		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Quantity,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan item: %w", err)
		}
		items = append(items, &item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating items: %w", err)
	}

	return items, nil
}

// UpdateItem updates an existing item
func (r *Repository) UpdateItem(ctx context.Context, item *models.Item) error {
	query := `
		UPDATE items
		SET name = $1, quantity = $2, updated_at = $3
		WHERE id = $4
	`

	item.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query, item.Name, item.Quantity, item.UpdatedAt, item.ID)
	if err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("item not found")
	}

	return nil
}

// DeleteItem deletes an item by ID
func (r *Repository) DeleteItem(ctx context.Context, id int) error {
	query := `DELETE FROM items WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("item not found")
	}

	return nil
}
