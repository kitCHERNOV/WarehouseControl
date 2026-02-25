package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"wbtechschool-L3/WarehouseControl/internal/models"
)

// setUserContext sets the current user context for triggers
// This uses SET LOCAL which is transaction-specific
func setUserContext(ctx context.Context, tx *sql.Tx, userID, userName string) error {
	const op = "repository.item.setUserContext"

	_, err := tx.ExecContext(ctx, "SET LOCAL app.current_user = $1", fmt.Sprintf("%s:%s", userID, userName))
	if err != nil {
		return fmt.Errorf("failed to set user context: Loc:%s, Err:%w", op, err)
	}
	return nil
}

// CreateItem creates a new item in the database
func (r *Repository) CreateItem(ctx context.Context, item *models.Item, userID, userName string) error {
	const op = "repository.item.CreateItem"

	// Start transaction for trigger support
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: Loc:%s, Err:%v", op, err)
	}
	defer tx.Rollback() // Will be ignored if transaction is committed

	// Set user context for triggers
	if err := setUserContext(ctx, tx, userID, userName); err != nil {
		return fmt.Errorf("failed to set user context: Loc:%s, Err:%v", op, err)
	}

	query := `
		INSERT INTO items (name, description, quantity, price, sku, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	now := time.Now()
	item.CreatedAt = now
	item.UpdatedAt = now

	err = tx.QueryRowContext(ctx, query,
		item.Name,
		item.Description,
		item.Quantity,
		item.Price,
		item.SKU,
		now,
		now,
	).Scan(&item.ID)
	if err != nil {
		return fmt.Errorf("failed to create item: Loc:%s, Err:%v", op, err)
	}

	// Commit transaction to apply changes
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: Loc:%s, Err:%v", op, err)
	}

	return nil
}

// GetItem retrieves an item by ID
func (r *Repository) GetItem(ctx context.Context, id int) (*models.Item, error) {
	const op = "repository.item.GetItem"

	query := `
		SELECT id, name, description, quantity, price, sku, created_at, updated_at
		FROM items
		WHERE id = $1
	`

	var item models.Item
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&item.ID,
		&item.Name,
		&item.Description,
		&item.Quantity,
		&item.Price,
		&item.SKU,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("item not found: Loc:%s, Err:%v", op, err)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get item: Loc:%s, Err:%v", op, err)
	}

	return &item, nil
}

// GetAllItems retrieves all items from the database
func (r *Repository) GetAllItems(ctx context.Context) ([]*models.Item, error) {
	const op = "repository.item.GetAllItems"

	query := `
		SELECT id, name, description, quantity, price, sku, created_at, updated_at
		FROM items
		ORDER BY id
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get items: Loc:%s, Err:%v", op, err)
	}
	defer rows.Close()

	var items []*models.Item
	for rows.Next() {
		var item models.Item
		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Description,
			&item.Quantity,
			&item.Price,
			&item.SKU,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan item: Loc:%s, Err:%v", op, err)
		}
		items = append(items, &item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating items: Loc:%s, Err:%v", op, err)
	}

	return items, nil
}

// UpdateItem updates an existing item
func (r *Repository) UpdateItem(ctx context.Context, item *models.Item, userID, userName string) error {
	const op = "repository.item.UpdateItem"

	// Start transaction for trigger support
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: Loc:%s, Err:%v", op, err)
	}
	defer tx.Rollback() // Will be ignored if transaction is committed

	// Set user context for triggers
	if err := setUserContext(ctx, tx, userID, userName); err != nil {
		return fmt.Errorf("failed to set user context: Loc:%s, Err:%v", op, err)
	}

	query := `
		UPDATE items
		SET name = $1, description = $2, quantity = $3, price = $4, sku = $5, updated_at = $6
		WHERE id = $7
	`

	item.UpdatedAt = time.Now()

	result, err := tx.ExecContext(ctx, query,
		item.Name,
		item.Description,
		item.Quantity,
		item.Price,
		item.SKU,
		item.UpdatedAt,
		item.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update item: Loc:%s, Err:%v", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: Loc:%s, Err:%v", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("item not found")
	}

	// Commit transaction to apply changes
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: Loc:%s, Err:%v", op, err)
	}

	return nil
}

// DeleteItem deletes an item by ID
func (r *Repository) DeleteItem(ctx context.Context, id int, userID, userName string) error {
	const op = "repository.item.DeleteItem"

	// Start transaction for trigger support
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: Loc:%s, Err:%v", op, err)
	}
	defer tx.Rollback() // Will be ignored if transaction is committed

	// Set user context for triggers
	if err := setUserContext(ctx, tx, userID, userName); err != nil {
		return fmt.Errorf("failed to set user context: Loc:%s, Err:%v", op, err)
	}

	query := `DELETE FROM items WHERE id = $1`

	result, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete item: Loc:%s, Err:%v", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: Loc:%s, Err:%v", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("item not found")
	}

	// Commit transaction to apply changes
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: Loc:%s, Err:%v", op, err)
	}

	return nil
}
