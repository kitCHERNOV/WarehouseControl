package service

import (
	"context"
	"fmt"

	"wbtechschool-L3/WarehouseControl/internal/models"
	"wbtechschool-L3/WarehouseControl/internal/repository"
)

// ItemService handles business logic for items
type ItemService struct {
	repo *repository.Repository
}

// NewItemService creates a new item service
func NewItemService(repo *repository.Repository) *ItemService {
	return &ItemService{repo: repo}
}

// CreateItem creates a new item
func (s *ItemService) CreateItem(ctx context.Context, req *models.CreateItemRequest, userID, userName string) (*models.Item, error) {
	// Validate input
	if req.Name == "" {
		return nil, fmt.Errorf("item name is required")
	}
	if req.Quantity < 0 {
		return nil, fmt.Errorf("quantity cannot be negative")
	}

	item := &models.Item{
		Name:     req.Name,
		Quantity: req.Quantity,
	}

	if err := s.repo.CreateItem(ctx, item); err != nil {
		return nil, fmt.Errorf("failed to create item: %w", err)
	}

	return item, nil
}

// GetItem retrieves an item by ID
func (s *ItemService) GetItem(ctx context.Context, id int) (*models.Item, error) {
	item, err := s.repo.GetItem(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get item: %w", err)
	}

	return item, nil
}

// GetAllItems retrieves all items
func (s *ItemService) GetAllItems(ctx context.Context) ([]*models.Item, error) {
	items, err := s.repo.GetAllItems(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get items: %w", err)
	}

	return items, nil
}

// UpdateItem updates an existing item
func (s *ItemService) UpdateItem(ctx context.Context, id int, req *models.UpdateItemRequest, userID, userName string) (*models.Item, error) {
	// Validate input
	if req.Name == "" {
		return nil, fmt.Errorf("item name is required")
	}
	if req.Quantity < 0 {
		return nil, fmt.Errorf("quantity cannot be negative")
	}

	// Check if item exists
	_, err := s.repo.GetItem(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get item: %w", err)
	}

	// Update item
	item := &models.Item{
		ID:       id,
		Name:     req.Name,
		Quantity: req.Quantity,
	}

	if err := s.repo.UpdateItem(ctx, item); err != nil {
		return nil, fmt.Errorf("failed to update item: %w", err)
	}

	return item, nil
}

// DeleteItem deletes an item by ID
func (s *ItemService) DeleteItem(ctx context.Context, id int, userID, userName string) error {
	// Check if item exists
	_, err := s.repo.GetItem(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get item: %w", err)
	}

	if err := s.repo.DeleteItem(ctx, id); err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}

	return nil
}
