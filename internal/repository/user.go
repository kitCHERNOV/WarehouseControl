package repository

import (
	"context"
	"database/sql"
	"fmt"

	"wbtechschool-L3/WarehouseControl/internal/models"
)

// GetUserByUsername retrieves a user by username
func (r *Repository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	const op = "repository.user.GetUserByUsername"

	query := `
		SELECT id, username, password, role
		FROM users
		WHERE username = $1
	`

	var user models.User
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Role,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found: Loc:%s, Err:%v", op, err)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: Loc:%s, Err:%v", op, err)
	}

	return &user, nil
}

// GetUserByID retrieves a user by ID
func (r *Repository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	const op = "repository.user.GetUserByID"

	query := `
		SELECT id, username, password, role
		FROM users
		WHERE id = $1
	`

	var user models.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Role,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found: Loc:%s, Err:%v", op, err)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: Loc:%s, Err:%v", op, err)
	}

	return &user, nil
}

// GetAllUsers retrieves all users
func (r *Repository) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	const op = "repository.user.GetAllUsers"

	query := `
		SELECT id, username, password, role
		FROM users
		ORDER BY id
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: Loc:%s, Err:%v", op, err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Password,
			&user.Role,
		); err != nil {
			return nil, fmt.Errorf("failed to scan user: Loc:%s, Err:%v", op, err)
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: Loc:%s, Err:%v", op, err)
	}

	return users, nil
}

// CreateUser creates a new user
func (r *Repository) CreateUser(ctx context.Context, user *models.User) error {
	const op = "repository.user.CreateUser"

	query := `
		INSERT INTO users (id, username, password, role)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.ExecContext(ctx, query, user.ID, user.Username, user.Password, user.Role)
	if err != nil {
		return fmt.Errorf("failed to create user: Loc:%s, Err:%v", op, err)
	}

	return nil
}
