package models

import "time"

// Item represents a warehouse inventory item
type Item struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Quantity  int       `json:"quantity" db:"quantity"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CreateItemRequest represents the request to create a new item
type CreateItemRequest struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

// UpdateItemRequest represents the request to update an item
type UpdateItemRequest struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}
