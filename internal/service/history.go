package service

import (
	"context"
	"fmt"

	"wbtechschool-L3/WarehouseControl/internal/models"
	"wbtechschool-L3/WarehouseControl/internal/repository"
)

// HistoryService handles business logic for item history
type HistoryService struct {
	repo *repository.Repository
}

// NewHistoryService creates a new history service
func NewHistoryService(repo *repository.Repository) *HistoryService {
	return &HistoryService{repo: repo}
}

// GetHistoryByItemID retrieves the history of changes for a specific item
func (s *HistoryService) GetHistoryByItemID(ctx context.Context, itemID int) ([]*models.History, error) {
	const op = "service.history.GetHistoryByItemID"

	history, err := s.repo.GetHistoryByItemID(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to get history: Loc:%s, Err:%v", op, err)
	}

	return history, nil
}

// GetAllHistory retrieves all history records
func (s *HistoryService) GetAllHistory(ctx context.Context) ([]*models.History, error) {
	const op = "service.history.GetAllHistory"

	history, err := s.repo.GetAllHistory(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get history: Loc:%s, Err:%v", op, err)
	}

	return history, nil
}
