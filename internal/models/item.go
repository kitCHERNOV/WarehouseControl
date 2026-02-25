package models

import "time"

// Item represents a warehouse inventory item
type Item struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Quantity    int       `json:"quantity" db:"quantity"`
	Price       float64   `json:"price" db:"price"`
	SKU         string    `json:"sku" db:"sku"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// CreateItemRequest represents the request to create a new item
type CreateItemRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	SKU         string  `json:"sku"`
}

// UpdateItemRequest represents the request to update an item
type UpdateItemRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	SKU         string  `json:"sku"`
}
