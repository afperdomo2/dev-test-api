package users

import (
	"time"

	"github.com/felipe/dev-test-api/internal/common"
	"github.com/felipe/dev-test-api/internal/models"
	"github.com/google/uuid"
)

var sortConfig = common.SortConfig{
	Allowed: []string{"email", "created_at", "updated_at"},
	Default: "created_at desc",
}

type CreateUserRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=72"`
	IsAdmin  bool   `json:"is_admin"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToUserResponse(u models.User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		IsAdmin:   u.IsAdmin,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
