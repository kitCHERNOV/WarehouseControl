package repository

import (
	"context"
	"database/sql"
	"fmt"

	"wbtechschool-L3/WarehouseControl/internal/models"
)

// GetHistoryByItemID retrieves the history of changes for a specific item
func (r *Repository) GetHistoryByItemID(ctx context.Context, itemID int) ([]*models.History, error) {
	const op = "repository.history.GetHistoryByItemID"

	query := `
		SELECT id, item_id, action, old_name, new_name, old_quantity, new_quantity, user_id, user_name, changed_at
		FROM item_history
		WHERE item_id = $1
		ORDER BY changed_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to get history: Loc:%s, Err:%v", op, err)
	}
	defer rows.Close()

	var history []*models.History
	for rows.Next() {
		var h models.History
		if err := rows.Scan(
			&h.ID,
			&h.ItemID,
			&h.Action,
			&h.OldName,
			&h.NewName,
			&h.OldQty,
			&h.NewQty,
			&h.UserID,
			&h.UserName,
			&h.ChangedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan history: Loc:%s, Err:%v", op, err)
		}
		history = append(history, &h)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating history: Loc:%s, Err:%v", op, err)
	}

	return history, nil
}

// GetAllHistory retrieves all history records
func (r *Repository) GetAllHistory(ctx context.Context) ([]*models.History, error) {
	const op = "repository.history.GetAllHistory"

	query := `
		SELECT id, item_id, action, old_name, new_name, old_quantity, new_quantity, user_id, user_name, changed_at
		FROM item_history
		ORDER BY changed_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get history: Loc:%s, Err:%v", op, err)
	}
	defer rows.Close()

	var history []*models.History
	for rows.Next() {
		var h models.History
		if err := rows.Scan(
			&h.ID,
			&h.ItemID,
			&h.Action,
			&h.OldName,
			&h.NewName,
			&h.OldQty,
			&h.NewQty,
			&h.UserID,
			&h.UserName,
			&h.ChangedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan history: Loc:%s, Err:%v", op, err)
		}
		history = append(history, &h)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating history: Loc:%s, Err:%v", op, err)
	}

	return history, nil
}

// CreateHistory creates a new history record
func (r *Repository) CreateHistory(ctx context.Context, history *models.History) error {
	const op = "repository.history.CreateHistory"

	query := `
		INSERT INTO item_history (item_id, action, old_name, new_name, old_quantity, new_quantity, user_id, user_name, changed_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		history.ItemID,
		history.Action,
		history.OldName,
		history.NewName,
		history.OldQty,
		history.NewQty,
		history.UserID,
		history.UserName,
		history.ChangedAt,
	).Scan(&history.ID)

	if err != nil {
		return fmt.Errorf("failed to create history: Loc:%s, Err:%v", op, err)
	}

	return nil
}

// GetHistoryByID retrieves a history record by ID
func (r *Repository) GetHistoryByID(ctx context.Context, id int) (*models.History, error) {
	const op = "repository.history.GetHistoryByID"

	query := `
		SELECT id, item_id, action, old_name, new_name, old_quantity, new_quantity, user_id, user_name, changed_at
		FROM item_history
		WHERE id = $1
	`

	var h models.History
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&h.ID,
		&h.ItemID,
		&h.Action,
		&h.OldName,
		&h.NewName,
		&h.OldQty,
		&h.NewQty,
		&h.UserID,
		&h.UserName,
		&h.ChangedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("history not found: Loc:%s, Err:%v", op, err)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get history: Loc:%s, Err:%v", op, err)
	}

	return &h, nil
}
