package repository

import (
	"database/sql"
	"fmt"
	"wbtechschool-L3/WarehouseControl/internal/config"

	_ "github.com/lib/pq"
)

// Repository represents the data access layer
type Repository struct {
	db *sql.DB
}

// New creates a new repository instance with database connection
func New(cfg config.Database) (*Repository, error) {
	const op = "repository.repository.New"

	// Build PostgreSQL connection string
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.SSLMode,
	)

	// Open database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: Loc:%s, Err:%v", op, err)
	}

	// Verify connection is working
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: Loc:%s, Err:%v", op, err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	return &Repository{db: db}, nil
}

// Close closes the database connection
func (r *Repository) Close() error {
	const op = "repository.repository.Close"

	if r.db != nil {
		return r.db.Close()
	}
	return nil
}

// DB returns the underlying sql.DB instance
func (r *Repository) DB() *sql.DB {
	return r.db
}
