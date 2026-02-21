package service

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"wbtechschool-L3/WarehouseControl/internal/config"
	"wbtechschool-L3/WarehouseControl/internal/models"
	"wbtechschool-L3/WarehouseControl/internal/repository"
)

// AuthService handles authentication and authorization
type AuthService struct {
	repo     *repository.Repository
	jwtCfg   config.JWT
	password string // For demo purposes, we'll use a simple password check
}

// NewAuthService creates a new auth service
func NewAuthService(repo *repository.Repository, jwtCfg config.JWT) *AuthService {
	return &AuthService{
		repo:   repo,
		jwtCfg: jwtCfg,
	}
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error) {
	// Get user from repository
	user, err := s.repo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials: %w", err)
	}

	// Check password (in production, use bcrypt)
	if user.Password != req.Password {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate JWT token
	token, err := s.generateToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Remove password from response
	user.Password = ""

	return &models.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

// generateToken generates a JWT token for a user
func (s *AuthService) generateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Duration(s.jwtCfg.Expiration) * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtCfg.Secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *AuthService) ValidateToken(tokenString string) (*models.Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtCfg.Secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	return &models.Claims{
		UserID:   claims["user_id"].(string),
		Username: claims["username"].(string),
		Role:     models.Role(claims["role"].(string)),
	}, nil
}

// HasPermission checks if a user has the required permission
func (s *AuthService) HasPermission(userRole models.Role, requiredRole models.Role) bool {
	roleHierarchy := map[models.Role]int{
		models.RoleViewer:  1,
		models.RoleManager: 2,
		models.RoleAdmin:   3,
	}

	return roleHierarchy[userRole] >= roleHierarchy[requiredRole]
}

// CanRead checks if a role can read items
func (s *AuthService) CanRead(role models.Role) bool {
	return true // All roles can read
}

// CanCreate checks if a role can create items
func (s *AuthService) CanCreate(role models.Role) bool {
	return role == models.RoleAdmin || role == models.RoleManager
}

// CanUpdate checks if a role can update items
func (s *AuthService) CanUpdate(role models.Role) bool {
	return role == models.RoleAdmin || role == models.RoleManager
}

// CanDelete checks if a role can delete items
func (s *AuthService) CanDelete(role models.Role) bool {
	return role == models.RoleAdmin
}
