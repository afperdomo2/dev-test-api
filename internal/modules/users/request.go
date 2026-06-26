package users

import "github.com/felipe/dev-test-api/internal/common"

var sortConfig = common.SortConfig{
	Allowed: []string{"email", "created_at", "updated_at"},
	Default: "created_at desc",
}

type CreateUserRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=72"`
	IsAdmin  bool   `json:"isAdmin"`
}

type UpdateUserRequest struct {
	Password string `json:"password,omitempty" binding:"omitempty,min=8,max=72"`
	IsAdmin  *bool  `json:"isAdmin,omitempty"`
}
