package auth

import (
	"github.com/felipe/dev-test-api/internal/modules/users"
)

type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=1"`
}

type SetupRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

type AuthResponse struct {
	Token string             `json:"token"`
	User  users.UserResponse `json:"user"`
}
