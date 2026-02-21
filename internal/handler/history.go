package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"wbtechschool-L3/WarehouseControl/internal/models"
	"wbtechschool-L3/WarehouseControl/internal/service"
)

// HistoryHandler handles HTTP requests for item history
type HistoryHandler struct {
	historyService *service.HistoryService
	authService    *service.AuthService
}

// NewHistoryHandler creates a new history handler
func NewHistoryHandler(historyService *service.HistoryService, authService *service.AuthService) *HistoryHandler {
	return &HistoryHandler{
		historyService: historyService,
		authService:    authService,
	}
}

// GetHistoryByItemID handles GET /items/{id}/history
func (h *HistoryHandler) GetHistoryByItemID(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user, ok := r.Context().Value("user").(*models.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Check permissions
	if !h.authService.CanRead(user.Role) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid item id", http.StatusBadRequest)
		return
	}

	history, err := h.historyService.GetHistoryByItemID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}

// GetAllHistory handles GET /history
func (h *HistoryHandler) GetAllHistory(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user, ok := r.Context().Value("user").(*models.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Check permissions
	if !h.authService.CanRead(user.Role) {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	history, err := h.historyService.GetAllHistory(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}
